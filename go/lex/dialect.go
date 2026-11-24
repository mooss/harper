package lex

// Dialect declares how to tokenize a Lisp dialect.
type Dialect[T any] struct {
	// Whitespace returns true when its argument is whitespace.
	Whitespace func(rune) bool
	// LineComment is the character starting a comment.
	LineComment rune
	// Rules is the list of lexing rules to try in order.
	//
	// Each Rule must obey the cursor contract:
	//  - When declining a match (returning matched == false), it must leave the Lexer cursor
	//    unmoved so later rules scan from the same position.
	//  - When reporting a match (returning matched == true), it must advance the cursor by
	//    consuming at least one rune so TokenIterator makes progress.
	Rules []Rule[T]
}

// Scanner builds a scanner for the dialect.
func (d Dialect[T]) Scanner() ScanFun[T] {
	return func(l *Lexer) (Token[T], error) {
		l.SkipWhitespace(d.Whitespace, d.LineComment)
		if l.Over() {
			return Token[T]{}, ErrEndOfInput
		}
		for _, rule := range d.Rules {
			token, matched, err := rule(l)
			if err != nil {
				return Token[T]{}, err
			}
			if matched {
				return token, nil
			}
		}
		l.Next()
		return Token[T]{}, l.Failure(InvalidStart)
	}
}
