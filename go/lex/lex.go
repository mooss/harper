package lex

import (
	"errors"
	"fmt"

	"github.com/mooss/harper/go/core"
)

// ScanFun scans the lexer's input and returns the next token.
type ScanFun[T any] func(*Lexer) (core.Token[T], error)

// LexicalFailure describes why lexing failed.
type LexicalFailure string

const (
	InvalidStart         LexicalFailure = "invalid token start"
	InvalidTokenBoundary LexicalFailure = "invalid token boundary"
)

const EOF = -1

var ErrEndOfInput = errors.New("end of input")

///////////
// Lexer //

// Lexer holds the input and cursor positions for a scan.
type Lexer struct {
	Input []rune
	start int
	pos   int
}

// Current returns the text consumed since the start of the current token.
func (l *Lexer) Current() string {
	return string(l.Input[l.start:l.pos])
}

// Failure returns an error describing the current token and advances the token start.
func (l *Lexer) Failure(reason LexicalFailure) error {
	token := l.Produce("")
	return fmt.Errorf("lexer: %s at %#q", reason, token.Value)
}

// Failuref returns a formatted error describing the current token and advances the token start.
func (l *Lexer) Failuref(reason LexicalFailure, msg string, args ...any) error {
	token := l.Produce("")
	args = append([]any{reason, token.Value}, args...)
	return fmt.Errorf("lexer: %s at %#q: "+msg, args...)
}

// Forward advances the position of the current token by at most n (does not go beyond len(Input)).
func (l *Lexer) Forward(n int) {
	l.pos = min(l.pos+n, len(l.Input))
}

// MatchRunes returns true when the given runes match with the input at the current position.
func (l *Lexer) MatchRunes(runes []rune) bool {
	for i, r := range runes {
		if l.Peek(i) != r {
			return false
		}
	}
	return true
}

// Next advances the cursor and returns the current rune, or EOF.
func (l *Lexer) Next() rune {
	if l.Over() {
		return EOF
	}
	r := l.Input[l.pos]
	l.pos++
	return r
}

// Over reports whether all input has been consumed.
func (l *Lexer) Over() bool {
	return l.pos >= len(l.Input)
}

// Peek returns the rune at the given offset without advancing, or EOF.
func (l *Lexer) Peek(index int) rune {
	index += l.pos
	if index < 0 || index >= len(l.Input) {
		return EOF
	}
	return l.Input[index]
}

// Pos returns the current position of the cursor.
func (l *Lexer) Pos() int { return l.pos }

// Produce returns the current token and advances the token start.
func (l *Lexer) Produce[T any](typ T) core.Token[T] {
	token := core.Token[T]{Type: typ, Value: l.Current()}
	l.start = l.pos
	return token
}

// SkipWhitespace consumes whitespace and line comments (use 0 for no line comment).
func (l *Lexer) SkipWhitespace(isWS func(rune) bool, lineComment rune) {
	for !l.Over() {
		r := l.Peek(0)
		if isWS != nil && isWS(r) {
			l.Next()
			l.start = l.pos
			continue
		}
		if lineComment != 0 && r == lineComment {
			for r = l.Next(); r != '\n' && r != '\r' && r != EOF; {
				r = l.Next()
			}
			l.start = l.pos
			continue
		}
		return
	}
}

// Start returns the current start of the cursor.
func (l *Lexer) Start() int { return l.start }

// SuccessOnBoundary succeeds only when the next rune is a token boundary.
func (l *Lexer) SuccessOnBoundary[T any](typ T, isBoundary func(rune) bool) (core.Token[T], error) {
	next := l.Peek(0)
	if next != EOF && !isBoundary(next) {
		l.Next()
		msg := "%q is not a boundary character"
		return core.Token[T]{}, l.Failuref(InvalidTokenBoundary, msg, next)
	}
	return l.Produce(typ), nil
}
