package rad

import (
	"fmt"
	"io"

	"github.com/mooss/harper/go/core"
)

// Reader constructs Lisp values from tokens.
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

// Done returns true when all the input has been read.
func (r *Reader[K, V]) Done() bool {
	return r.pos >= len(r.tokens)
}

// Next returns the next token (zero when exhausted) and advances the input.
func (r *Reader[K, V]) Next() core.Token[K] {
	t := r.Peek()
	if r.pos < len(r.tokens) {
		r.pos++
	}
	return t
}

// Peek returns the next token (zero when exhausted) without advancing the input.
func (r *Reader[K, V]) Peek() core.Token[K] {
	if r.Done() {
		var zero core.Token[K]
		return zero
	}
	return r.tokens[r.pos]
}

func slc[T any](args ...T) []T { return args }

// ReadForm constructs the next type from the input tokens, after ensuring that there are tokens
// left to read.
// Returns an error wrapping io.EOF when there are no more tokens left to read.
func (r *Reader[K, V]) ReadForm() (V, error) {
	var zero V
	if r.Done() {
		return zero, r.erf("%w", io.EOF)
	}
	t := r.Next()
	for _, fn := range slc(r.readPrefix, r.readCompound) {
		if v, handled, err := fn(t); err != nil {
			return zero, err
		} else if handled {
			return v, nil
		}
	}
	return r.readAtom(t)
}

// Rest returns all the tokens that are yet to be read.
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

// erf returns a reader error formatted with msg and args.
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

// readCompound dispatches a read delimiter to the appropriate read* method.
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

// readSequence constructs a sequence type from the remaining input tokens (assumes that the opening
// delimiter has already been consumed).
func (r *Reader[K, V]) readSequence(untilValue string, untilType K) ([]V, error) {
	var res []V
	for {
		if r.Done() {
			return nil, r.erf("unbalanced expression (missing %q)", untilValue)
		}
		p := r.Peek()
		if p.Type == untilType && p.Value == untilValue {
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

// wrap wraps err into a reader error.
func (r *Reader[K, V]) wrap(v V, err error) (V, bool, error) {
	if err != nil {
		var zero V
		return zero, true, r.erf("%w", err)
	}
	return v, true, nil
}
