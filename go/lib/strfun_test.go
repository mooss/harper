//revive:disable:line-length-limit
package lib_test

import (
	"testing"

	"github.com/mooss/bagend/go/assert"
	"github.com/mooss/harper/go/lib"
)

type st = assert.Case[string, string]

func TestQuote(t *testing.T) {
	a := assert.New(t)
	a.RunCases("Quote", func(a assert.A, c st) {
		got := lib.Quote(c.Input)
		a.Equal(got, c.Want)
		a.StringStartsWith(got, `"`)
		a.StringEndsWith(got, `"`)
	},
		st{Name: "empty", Input: "", Want: `""`},
		st{Name: "single", Input: "a", Want: `"a"`},
		st{Name: "plain", Input: "hi", Want: `"hi"`},
		st{Name: "quote", Input: "a\"b", Want: "\"a\\\"b\""},
		st{Name: "backslash", Input: "a\\b", Want: "\"a\\\\b\""},
		st{Name: "newline", Input: "a\nb", Want: "\"a\\nb\""},
		st{Name: "only quote", Input: `"`, Want: "\"\\\"\""},
		st{Name: "only backslash", Input: `\`, Want: "\"\\\\\""},
		st{Name: "only newline", Input: "\n", Want: "\"\\n\""},
		st{Name: "all specials", Input: "\"\\\n", Want: "\"\\\"\\\\\\n\""},
		st{Name: "tab left as-is", Input: "a\tb", Want: "\"a\tb\""},
		st{Name: "cr left as-is", Input: "a\rb", Want: "\"a\rb\""},
		st{Name: "single quotes left", Input: "a'b", Want: "\"a'b\""},
		st{Name: "parens semicolon", Input: "()'; ", Want: "\"()'; \""},
		st{Name: "multiline", Input: "line1\nline2", Want: "\"line1\\nline2\""},
		st{Name: "multiple quotes", Input: `""`, Want: "\"\\\"\\\"\""},
		st{Name: "multiple backslashes", Input: `\\`, Want: "\"\\\\\\\\\""},
		st{Name: "multiple newlines", Input: "\n\n", Want: "\"\\n\\n\""},
		st{Name: "mixed", Input: "x\"y\\z\nw", Want: "\"x\\\"y\\\\z\\nw\""},
		st{Name: "null byte left", Input: "a\x00b", Want: "\"a\x00b\""},
		st{Name: "utf8 left", Input: "héllo", Want: "\"héllo\""},
	)
}

func TestUnquote(t *testing.T) {
	a := assert.New(t)
	a.RunCases("Unquote", func(a assert.A, c st) {
		got, err := lib.Unquote(c.Input)
		if c.WantErr {
			a.ErrorContains(err, c.Want)
		} else {
			a.NoError(err)
			a.Equal(got, c.Want)
		}
	},
		st{Name: "empty quoted", Input: `""`, Want: ""},
		st{Name: "single", Input: `"a"`, Want: "a"},
		st{Name: "plain", Input: `"hi"`, Want: "hi"},
		st{Name: "escaped quote", Input: "\"a\\\"b\"", Want: "a\"b"},
		st{Name: "escaped backslash", Input: "\"a\\\\b\"", Want: "a\\b"},
		st{Name: "escaped newline", Input: "\"a\\nb\"", Want: "a\nb"},
		st{Name: "only escaped quote", Input: "\"\\\"\"", Want: `"`},
		st{Name: "only escaped backslash", Input: "\"\\\\\\\\\"", Want: `\\`},
		st{Name: "only escaped newline", Input: "\"\\n\"", Want: "\n"},
		st{Name: "all specials", Input: "\"\\\"\\\\\\n\"", Want: "\"\\\n"},
		st{Name: "tab left as-is", Input: "\"a\tb\"", Want: "a\tb"},
		st{Name: "cr left as-is", Input: "\"a\rb\"", Want: "a\rb"},
		st{Name: "single quotes left", Input: "\"a'b\"", Want: "a'b"},
		st{Name: "parens semicolon", Input: "\"()'; \"", Want: "()'; "},
		st{Name: "multiline", Input: "\"line1\\nline2\"", Want: "line1\nline2"},
		st{Name: "null byte left", Input: "\"a\x00b\"", Want: "a\x00b"},
		st{Name: "empty input", Input: "", Want: "missing surrounding quotes", WantErr: true},
		st{Name: "single quote alone", Input: `"`, Want: "missing surrounding quotes", WantErr: true},
		st{Name: "single char", Input: "a", Want: "missing surrounding quotes", WantErr: true},
		st{Name: "missing leading", Input: `a"`, Want: "missing surrounding quotes", WantErr: true},
		st{Name: "missing trailing", Input: `"a`, Want: "missing surrounding quotes", WantErr: true},
		st{Name: "single quotes", Input: "'a'", Want: "missing surrounding quotes", WantErr: true},
		st{Name: "unquoted", Input: "hi", Want: "missing surrounding quotes", WantErr: true},
		st{Name: "trailing backslash", Input: "\"a\\\"", Want: "invalid escape", WantErr: true},
		st{Name: "lone backslash", Input: "\"\\\"", Want: "invalid escape", WantErr: true},
		st{Name: "invalid t", Input: "\"\\t\"", Want: "invalid escape \\t", WantErr: true},
		st{Name: "invalid r", Input: "\"\\r\"", Want: "invalid escape \\r", WantErr: true},
		st{Name: "invalid x", Input: "\"\\x\"", Want: "invalid escape \\x", WantErr: true},
		st{Name: "invalid b", Input: "\"\\b\"", Want: "invalid escape \\b", WantErr: true},
		st{Name: "invalid 0", Input: "\"\\0\"", Want: "invalid escape \\0", WantErr: true},
		st{Name: "invalid single quote", Input: "\"\\'\"", Want: "invalid escape \\'", WantErr: true},
		st{Name: "invalid space", Input: "\"\\ \"", Want: "invalid escape \\ ", WantErr: true},
		st{Name: "invalid n upper", Input: "\"\\N\"", Want: "invalid escape \\N", WantErr: true},
	)
}

func TestQuoteUnquoteRoundTrip(t *testing.T) {
	a := assert.New(t)
	a.RunCases("QuoteUnquoteRoundTrip", func(a assert.A, c st) {
		quoted := lib.Quote(c.Input)
		a.Equal(quoted, c.Want)
		a.StringStartsWith(quoted, `"`)
		a.StringEndsWith(quoted, `"`)
		a.Equal(a.Must(lib.Unquote(quoted)), c.Input)
		a.Equal(a.Must(lib.Unquote(c.Want)), c.Input)
	},
		st{Name: "empty", Input: "", Want: `""`},
		st{Name: "hi", Input: "hi", Want: `"hi"`},
		st{Name: "all specials", Input: "\"\\\n", Want: "\"\\\"\\\\\\n\""},
		st{Name: "multiline", Input: "line1\nline2", Want: "\"line1\\nline2\""},
		st{Name: "parens", Input: "()'; ", Want: "\"()'; \""},
	)
}
