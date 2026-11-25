package lex

import (
	"errors"
	"iter"

	"github.com/mooss/harper/go/core"
)

// TokenIterator builds an iterable sequence of tokens and an error getter from input data and a
// scanner.
// The error getter should be checked after iteration.
func TokenIterator[T any](input string, scan ScanFun[T]) (iter.Seq[core.Token[T]], func() error) {
	l := &Lexer{Input: []rune(input)}
	var err error
	return func(yield func(core.Token[T]) bool) {
		if err != nil {
			return
		}
		for {
			token, scanErr := scan(l)
			if scanErr != nil {
				if errors.Is(scanErr, ErrEndOfInput) {
					return
				}
				err = scanErr
				return
			}
			if !yield(token) {
				return
			}
		}
	}, func() error { return err }
}
