package lex

import "strings"

// IsDigitASCII reports whether r is an ASCII decimal digit.
func IsDigitASCII(r rune) bool {
	return '0' <= r && r <= '9'
}

// InSet returns a predicate that reports whether a rune occurs in s.
func InSet(s string) func(rune) bool {
	return func(r rune) bool {
		return strings.ContainsRune(s, r)
	}
}
