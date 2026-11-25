package lex

import (
	"errors"
	"iter"

	"github.com/mooss/harper/go/core"
)

// TokenIterator builds an iterable sequence of tokens.
func TokenIterator[T any](input string, scan ScanFun[T]) iter.Seq2[core.Token[T], error] {
	l := NewLexer(input)
	return func(yield func(core.Token[T], error) bool) {
		for {
			token, err := scan(l)
			if err != nil {
				if errors.Is(err, ErrEndOfInput) {
					return
				}
				yield(core.Token[T]{}, err)
				return
			}
			if !yield(token, nil) {
				return
			}
		}
	}
}
