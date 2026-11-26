// This file implements string manipulation functions.

package lib

import (
	"errors"
	"fmt"
)

// Quote returns s as a double quoted string, including the quotes, escaping backslashes, double
// quotes, and newlines.
//
// The escaping rules are:
//   - `\` is escaped as `\\`.
//   - `"` is escaped as `\"`.
//   - newline is escaped as `\n`.
//
// All other bytes are left as-is.
//
// The result always starts and ends with a double quote (`"`), so passing it to Unquote recovers
// the original string.
func Quote(s string) string {
	out := make([]byte, 0, len(s)+2)
	out = append(out, '"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		default:
			out = append(out, s[i])
		}
	}
	out = append(out, '"')
	return string(out)
}

// Unquote takes a double quoted string, including the quotes, and returns its unquoted contents
// with escape sequences resolved.
//
// The input must start and end with a double quote (`"`).
//
// The resolution rules are:
//   - `\"` resolves to `"`.
//   - `\\` resolves to `\`.
//   - `\n` resolves to a newline.
//
// All other bytes are left as-is.
//
// A backslash followed by any other byte or followed by nothing returns an `invalid escape` error.
//
//revive:disable:cognitive-complexity
func Unquote(quoted string) (string, error) {
	if len(quoted) < 2 || quoted[0] != '"' || quoted[len(quoted)-1] != '"' {
		return "", errors.New("missing surrounding quotes")
	}
	inner := quoted[1 : len(quoted)-1]
	out := make([]byte, 0, len(inner))
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' {
			if i+1 >= len(inner) {
				return "", errors.New("invalid escape \\")
			}
			next := inner[i+1]
			switch next {
			case 'n':
				out = append(out, '\n')
				i++
				continue
			case '"', '\\':
				out = append(out, next)
				i++
				continue
			default:
				return "", fmt.Errorf("invalid escape \\%c", next)
			}
		}
		out = append(out, inner[i])
	}
	return string(out), nil
}
