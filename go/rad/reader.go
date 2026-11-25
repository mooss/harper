package rad

import (
	"fmt"
	"io"

	"github.com/mooss/harper/go/core"
)

type Reader[K comparable, V any] struct {
	tokens []core.Token[K]
	pos    int
	spec   *Spec[K, V]
}

func (s *Spec[K, V]) NewReader(tokens []core.Token[K]) *Reader[K, V] {
	return &Reader[K, V]{tokens: tokens, pos: 0, spec: s}
}

//////////////////////
// Public interface //

func (r *Reader[K, V]) Done() bool {
	return r.pos >= len(r.tokens)
}

func (r *Reader[K, V]) Next() core.Token[K] {
	t := r.Peek()
	if r.pos < len(r.tokens) {
		r.pos++
	}
	return t
}

func (r *Reader[K, V]) Peek() core.Token[K] {
	if r.Done() {
		var zero core.Token[K]
		return zero
	}
	return r.tokens[r.pos]
}

func (r *Reader[K, V]) ReadForm() (V, error) {
	var zero V
	if r.Done() {
		return zero, r.erf("%w", io.EOF)
	}
	t := r.Next()
	for _, fn := range []func(core.Token[K]) (V, bool, error){r.readPrefix, r.readCompound} {
		if v, handled, err := fn(t); err != nil {
			return zero, err
		} else if handled {
			return v, nil
		}
	}
	return r.readAtom(t)
}

func (r *Reader[K, V]) Rest() []core.Token[K] {
	if r.pos >= len(r.tokens) {
		return nil
	}
	return r.tokens[r.pos:]
}

/////////////////////
// Private methods //

func (r *Reader[K, V]) buildAtom(t core.Token[K]) (V, error) {
	var zero V
	if atom, ok := r.spec.atoms[t.Type]; ok {
		v, err := atom(t)
		if err != nil {
			return zero, r.erf("%w", err)
		}
		return v, nil
	}
	return zero, r.erf("unexpected token %q", t.Value)
}

func (*Reader[K, V]) erf(format string, a ...any) error {
	return fmt.Errorf("reader: "+format, a...)
}

func (r *Reader[K, V]) readAtom(t core.Token[K]) (V, error) {
	var zero V

	if r.spec.delimKinds[t.Type] {
		if r.spec.closes[t.Value] {
			return zero, r.erf("unexpected closing %q", t.Value)
		}
		return zero, r.erf("unexpected delimiter %q", t.Value)
	}

	if r.spec.prefixKinds[t.Type] {
		return zero, r.erf("unexpected sigil %q", t.Value)
	}

	return r.buildAtom(t)
}

func (r *Reader[K, V]) readCompound(t core.Token[K]) (V, bool, error) {
	var zero V
	open, ok := r.spec.opens[token[K]{t.Type, t.Value}]
	if !ok {
		return zero, false, nil
	}
	forms, err := r.readSequence(open.close, open.kind)
	if err != nil {
		return zero, true, err
	}
	if open.isTable && len(forms)%open.arity != 0 {
		msg := "incomplete table literal %q: expected multiple of %d forms"
		return zero, true, r.erf(msg, open.close, open.arity)
	}
	return r.wrap(open.build(forms))
}

func (r *Reader[K, V]) readPrefix(t core.Token[K]) (V, bool, error) {
	var zero V
	expand, ok := r.spec.prefixes[token[K]{t.Type, t.Value}]
	if !ok {
		return zero, false, nil
	}
	if r.Done() {
		return zero, true, r.erf("dangling prefix %q (missing form)", t.Value)
	}
	n, err := r.ReadForm()
	if err != nil {
		return zero, true, err
	}
	return r.wrap(expand(n))
}

func (r *Reader[K, V]) readSequence(until string, kind K) ([]V, error) {
	var res []V
	for {
		if r.Done() {
			return nil, r.erf("unbalanced expression (missing %q)", until)
		}
		p := r.Peek()
		if p.Type == kind && p.Value == until {
			r.Next()
			break
		}
		f, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (r *Reader[K, V]) wrap(v V, err error) (V, bool, error) {
	if err != nil {
		var zero V
		return zero, true, r.erf("%w", err)
	}
	return v, true, nil
}
