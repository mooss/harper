package rad

import (
	"errors"
	"fmt"
	"io"

	"github.com/mooss/harper/go/lex"
)

func (s *Spec[K, V]) ParseOne(tokens []lex.Token[K]) (V, error) {
	r := s.NewReader(tokens)
	v, err := r.ReadForm()
	if err != nil {
		return v, err
	}
	if !r.Done() {
		var zero V
		return zero, fmt.Errorf("reader: leftover tokens remain: %v", r.Rest())
	}
	return v, nil
}

func (s *Spec[K, V]) ParseAll(tokens []lex.Token[K]) ([]V, error) {
	r := s.NewReader(tokens)
	var out []V
	for {
		v, err := r.ReadForm()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		out = append(out, v)
	}
	if !r.Done() {
		return nil, fmt.Errorf("reader: leftover tokens remain: %v", r.Rest())
	}
	return out, nil
}
