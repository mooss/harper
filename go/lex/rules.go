package lex

import "github.com/mooss/bagend/go/fun/eager/lie"

// Rule scans a token when its leading input matches.
// Must return an error when the rule was activated but the input was invalid.
// Rules must always be called on a non-empty Lexer and can assume the first character is not EOF.
type Rule[T any] func(*Lexer[T]) (Token[T], bool, error)

// AnyOf matches any rune from chars.
func AnyOf[T any](chars string, typ T) Rule[T] {
	contains := InSet(chars)
	return func(l *Lexer[T]) (Token[T], bool, error) {
		if !contains(l.Peek(0)) {
			return Token[T]{}, false, nil
		}
		l.Next()
		return l.Produce(typ), true, nil
	}
}

// Quoted matches a token quoted between delimiter, honoring escape characters and failing with
// unterminated.
func Quoted[T any](delimiter, escape rune, unterminated LexicalFailure, typ T) Rule[T] {
	return func(l *Lexer[T]) (Token[T], bool, error) {
		if l.Peek(0) != delimiter {
			return Token[T]{}, false, nil
		}
		l.Next()
		for {
			switch l.Next() {
			case delimiter:
				return l.Produce(typ), true, nil
			case escape:
				if l.Next() == EOF {
					return Token[T]{}, true, l.Failure(unterminated)
				}
			case EOF:
				return Token[T]{}, true, l.Failure(unterminated)
			}
		}
	}
}

// FirstPrefix returns the first match of the given prefixes.
func FirstPrefix[T any](typ T, prefixes ...string) Rule[T] {
	runePrefixes := lie.Map(func(s string) []rune { return []rune(s) }, prefixes)
	return func(l *Lexer[T]) (Token[T], bool, error) {
		for _, runes := range runePrefixes {
			if len(runes) == 0 {
				continue
			}
			if l.MatchRunes(runes) {
				l.Forward(len(runes))
				return l.Produce(typ), true, nil
			}
		}

		return Token[T]{}, false, nil
	}
}

// SignedInt matches an optionally signed decimal integer.
func SignedInt[T any](isDigit func(rune) bool, isBoundary func(rune) bool, typ T) Rule[T] {
	isDigit = notEOFAnd(isDigit)
	return func(l *Lexer[T]) (Token[T], bool, error) {
		offset := 0
		if l.Peek(0) == '+' || l.Peek(0) == '-' {
			offset = 1
		}
		if !isDigit(l.Peek(offset)) {
			return Token[T]{}, false, nil
		}
		if offset != 0 {
			l.Next()
		}
		for isDigit(l.Peek(0)) {
			l.Next()
		}
		token, err := l.SuccessOnBoundary(typ, isBoundary)
		return token, true, err
	}
}

// Symbol matches a symbol with separate start and rest predicates.
func Symbol[T any](isStart, isRest func(rune) bool, isBoundary func(rune) bool, typ T) Rule[T] {
	isRest = notEOFAnd(isRest)
	return func(l *Lexer[T]) (Token[T], bool, error) {
		if !isStart(l.Peek(0)) {
			return Token[T]{}, false, nil
		}
		l.Next()
		for isRest(l.Peek(0)) {
			l.Next()
		}
		token, err := l.SuccessOnBoundary(typ, isBoundary)
		return token, true, err
	}
}

func notEOFAnd(pred func(rune) bool) func(rune) bool {
	return func(r rune) bool { return r != EOF && pred(r) }
}
