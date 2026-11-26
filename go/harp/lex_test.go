//revive:disable:line-length-limit
package harp

import (
	"testing"

	"github.com/mooss/bagend/go/assert"
)

func mustTokenize(a assert.A, program string) []Token {
	a.Helper()
	return a.Must(Tokenize(program))
}

type ic = assert.Case[string, struct{}] // Input case.
type tc = assert.Case[string, []Token]  // Token case.
type ec = assert.Case[string, string]   // Error case.

// Token utilities to reduce boilerplate.
func del(v string) Token         { return Token{Type: KindDelimiter, Value: v} }
func sym(v string) Token         { return Token{Type: KindSymbol, Value: v} }
func num(v string) Token         { return Token{Type: KindInt, Value: v} }
func str(v string) Token         { return Token{Type: KindDQS, Value: v} }
func sig(v string) Token         { return Token{Type: KindSigil, Value: v} }
func toks(args ...Token) []Token { return args }

func TestTokenizeEmpty(t *testing.T) {
	a := assert.New(t)
	a.RunCases("empty input yields no tokens", func(a assert.A, c ic) {
		got := mustTokenize(a, c.Input)
		a.Equal(len(got), 0)
	}, ic{Name: "empty string", Input: ""},
		ic{Name: "spaces", Input: "   "},
		ic{Name: "tabs and newlines", Input: " \t\n\r \n\t "},
		ic{Name: "only comment", Input: "; hello"},
		ic{Name: "only comment with newline", Input: "; hello\n"},
		ic{Name: "multiple comments", Input: "; one\n; two\n ; three"},
		ic{Name: "whitespace and comments", Input: "  ; c1 \n \t ; c2\r\n  "},
	)
}

func TestTokenizeDelimiters(t *testing.T) {
	a := assert.New(t)
	a.RunCases("delimiters", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "open", Input: "(", Want: toks(del("("))},
		tc{Name: "close", Input: ")", Want: toks(del(")"))},
		tc{Name: "pair", Input: "()", Want: toks(del("("), del(")"))},
		tc{Name: "pair with space", Input: "( )", Want: toks(del("("), del(")"))},
		tc{Name: "nested", Input: "((()))", Want: toks(del("("), del("("), del("("), del(")"), del(")"), del(")"))},
		tc{Name: "delimiters around symbol", Input: "(a)", Want: toks(del("("), sym("a"), del(")"))},
		tc{Name: "no space needed", Input: "(a(b)c)", Want: toks(del("("), sym("a"), del("("), sym("b"), del(")"), sym("c"), del(")"))},
	)
}

func TestTokenizeWhitespaceAndComments(t *testing.T) {
	a := assert.New(t)
	a.RunCases("whitespace skipping", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "space separated", Input: "a b", Want: toks(sym("a"), sym("b"))},
		tc{Name: "tab separated", Input: "a\tb", Want: toks(sym("a"), sym("b"))},
		tc{Name: "newline separated", Input: "a\nb", Want: toks(sym("a"), sym("b"))},
		tc{Name: "crlf separated", Input: "a\r\nb", Want: toks(sym("a"), sym("b"))},
		tc{Name: "leading trailing", Input: "  (a)  ", Want: toks(del("("), sym("a"), del(")"))},
		tc{Name: "comment to eol", Input: "a ; comment\nb", Want: toks(sym("a"), sym("b"))},
		tc{Name: "comment at eof no newline", Input: "a ; comment", Want: toks(sym("a"))},
		tc{Name: "comment between delimiters", Input: "(; comment\n)", Want: toks(del("("), del(")"))},
		tc{Name: "comment consumes parens", Input: "(a) ; )(((", Want: toks(del("("), sym("a"), del(")"))},
		tc{Name: "comment with CR", Input: "a ; comment\rb", Want: toks(sym("a"), sym("b"))},
		tc{Name: "empty comment then token", Input: ";\na", Want: toks(sym("a"))},
		tc{Name: "semicolon starts comment not symbol", Input: ";(a)", Want: nil},
	)
}

func TestTokenizeIntegers(t *testing.T) {
	a := assert.New(t)
	a.RunCases("integers", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "zero", Input: "0", Want: toks(num("0"))},
		tc{Name: "simple", Input: "42", Want: toks(num("42"))},
		tc{Name: "leading zeros", Input: "007", Want: toks(num("007"))},
		tc{Name: "plus", Input: "+42", Want: toks(num("+42"))},
		tc{Name: "minus", Input: "-42", Want: toks(num("-42"))},
		tc{Name: "plus zero", Input: "+0", Want: toks(num("+0"))},
		tc{Name: "minus zero", Input: "-0", Want: toks(num("-0"))},
		tc{Name: "in list", Input: "(1 2 3)", Want: toks(del("("), num("1"), num("2"), num("3"), del(")"))},
		tc{Name: "boundary paren", Input: "123)", Want: toks(num("123"), del(")"))},
		tc{Name: "boundary open", Input: "(123", Want: toks(del("("), num("123"))},
		tc{Name: "boundary comment", Input: "123;hi", Want: toks(num("123"))},
		tc{Name: "boundary whitespace", Input: "  9  ", Want: toks(num("9"))},
		tc{Name: "lone plus is symbol", Input: "+", Want: toks(sym("+"))},
		tc{Name: "lone minus is symbol", Input: "-", Want: toks(sym("-"))},
		tc{Name: "signed symbol", Input: "+foo", Want: toks(sym("+foo"))},
		tc{Name: "signed symbol minus", Input: "-foo", Want: toks(sym("-foo"))},
		tc{Name: "plus minus symbols", Input: "+ -", Want: toks(sym("+"), sym("-"))},
	)
}

func TestTokenizeSymbols(t *testing.T) {
	a := assert.New(t)
	a.RunCases("symbols", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "simple", Input: "foo", Want: toks(sym("foo"))},
		tc{Name: "define", Input: "define", Want: toks(sym("define"))},
		tc{Name: "with digits", Input: "a1", Want: toks(sym("a1"))},
		tc{Name: "operators", Input: "<=>", Want: toks(sym("<=>"))},
		tc{Name: "punctuation", Input: "!&*+-./:<=>?@", Want: toks(sym("!&*+-./:<=>?@"))},
		tc{Name: "underscore", Input: "_foo_bar_", Want: toks(sym("_foo_bar_"))},
		tc{Name: "slash", Input: "a/b", Want: toks(sym("a/b"))},
		tc{Name: "question", Input: "null?", Want: toks(sym("null?"))},
		tc{Name: "bang", Input: "set!", Want: toks(sym("set!"))},
		tc{Name: "upper", Input: "ABC", Want: toks(sym("ABC"))},
		tc{Name: "mixed case", Input: "AbC123", Want: toks(sym("AbC123"))},
		tc{Name: "adjacent to parens", Input: "foo(bar)baz", Want: toks(sym("foo"), del("("), sym("bar"), del(")"), sym("baz"))},
		tc{Name: "symbol then int", Input: "x 1", Want: toks(sym("x"), num("1"))},
		tc{Name: "int then symbol", Input: "1 x", Want: toks(num("1"), sym("x"))},
	)
}

func TestTokenizeSigil(t *testing.T) {
	a := assert.New(t)
	a.RunCases("sigil quote", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "lone", Input: "'", Want: toks(sig("'"))},
		tc{Name: "double", Input: "''", Want: toks(sig("'"), sig("'"))},
		tc{Name: "quote symbol", Input: "'a", Want: toks(sig("'"), sym("a"))},
		tc{Name: "quote list", Input: "'(a)", Want: toks(sig("'"), del("("), sym("a"), del(")"))},
		tc{Name: "quote int", Input: "'42", Want: toks(sig("'"), num("42"))},
		tc{Name: "quote string", Input: "'\"hi\"", Want: toks(sig("'"), str("\"hi\""))},
		tc{Name: "space between", Input: "' 'a", Want: toks(sig("'"), sig("'"), sym("a"))},
	)
}

func TestTokenizeStrings(t *testing.T) {
	a := assert.New(t)
	a.RunCases("double quoted strings", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{Name: "empty", Input: "\"\"", Want: toks(str("\"\""))},
		tc{Name: "simple", Input: "\"hi\"", Want: toks(str("\"hi\""))},
		tc{Name: "with spaces", Input: "\"hello world\"", Want: toks(str("\"hello world\""))},
		tc{Name: "with parens", Input: "\"()\"", Want: toks(str("\"()\""))},
		tc{Name: "with semicolon", Input: "\"; not a comment\"", Want: toks(str("\"; not a comment\""))},
		tc{Name: "with single quote", Input: "\"'\"", Want: toks(str("\"'\""))},
		tc{Name: "escaped quote", Input: "\"a\\\"b\"", Want: toks(str("\"a\\\"b\""))},
		tc{Name: "escaped backslash", Input: "\"a\\\\b\"", Want: toks(str("\"a\\\\b\""))},
		tc{Name: "multiline", Input: "\"a\nb\"", Want: toks(str("\"a\nb\""))},
		tc{Name: "adjacent strings", Input: "\"a\"\"b\"", Want: toks(str("\"a\""), str("\"b\""))},
		tc{Name: "string in list", Input: "(\"a\" b)", Want: toks(del("("), str("\"a\""), sym("b"), del(")"))},
		tc{Name: "unicode content", Input: "\"héllo\"", Want: toks(str("\"héllo\""))},
	)
}

func TestTokenizeErrors(t *testing.T) {
	a := assert.New(t)
	a.RunCases("lexical failures", func(a assert.A, c ec) {
		_, err := Tokenize(c.Input)
		if c.Want == "" {
			a.NoError(err)
		} else {
			a.ErrorContains(err, c.Want)
		}
	}, ec{Name: "hash", Input: "#", Want: "invalid token start"},
		ec{Name: "hash in list", Input: "(#)", Want: "invalid token start"},
		ec{Name: "comma", Input: ",", Want: "invalid token start"},
		ec{Name: "bracket open", Input: "[", Want: "invalid token start"},
		ec{Name: "bracket close", Input: "]", Want: "invalid token start"},
		ec{Name: "brace", Input: "{a}", Want: "invalid token start"},
		ec{Name: "backtick", Input: "`", Want: "invalid token start"},
		ec{Name: "backslash", Input: "\\", Want: "invalid token start"},
		ec{Name: "pipe", Input: "|", Want: "invalid token start"},
		ec{Name: "non ascii start", Input: "é", Want: "invalid token start"},
		ec{Name: "digit start is int", Input: "0"},
		ec{Name: "unterminated", Input: "\"abc", Want: "unbalanced double quote"},
		ec{Name: "unterminated empty", Input: "\"", Want: "unbalanced double quote"},
		ec{Name: "unterminated after escape", Input: "\"abc\\", Want: "unbalanced double quote"},
		ec{Name: "unterminated multiline", Input: "\"a\nb", Want: "unbalanced double quote"},
		ec{Name: "int boundary letter", Input: "123abc", Want: "invalid token boundary"},
		ec{Name: "int boundary quote", Input: "123\"a\"", Want: "invalid token boundary"},
		ec{Name: "int boundary sigil", Input: "123'a", Want: "invalid token boundary"},
		ec{Name: "symbol boundary quote", Input: "foo\"bar\"", Want: "invalid token boundary"},
		ec{Name: "symbol boundary sigil", Input: "foo'bar", Want: "invalid token boundary"},
	)

	// The "digit start is int" entry above is not an error; verify separately
	// that digits never produce invalid start.
	a.Run("digits are ints not errors", func(a assert.A) {
		for _, tc := range []string{"0", "1", "9", "+1", "-1"} {
			_, err := Tokenize(tc)
			a.NoError(err)
		}
	})

	a.Run("error mentions offending text", func(a assert.A) {
		_, err := Tokenize("#oops")
		a.ErrorContains(err, "invalid token start")
		a.StringContains(err.Error(), "#")
	})

	a.Run("boundary error mentions char", func(a assert.A) {
		_, err := Tokenize("12x")
		a.ErrorContains(err, "invalid token boundary")
	})
}

func TestTokenizePrograms(t *testing.T) {
	a := assert.New(t)
	a.RunCases("full programs", func(a assert.A, c tc) {
		a.Equal(mustTokenize(a, c.Input), c.Want)
	}, tc{
		Name:  "define",
		Input: "(define x 42)",
		Want:  toks(del("("), sym("define"), sym("x"), num("42"), del(")")),
	},
		tc{
			Name:  "lambda with comment",
			Input: "(lambda (x) ; identity\n  x)",
			Want:  toks(del("("), sym("lambda"), del("("), sym("x"), del(")"), sym("x"), del(")")),
		},
		tc{
			Name:  "quoted list",
			Input: "'(1 2)",
			Want:  toks(sig("'"), del("("), num("1"), num("2"), del(")")),
		},
		tc{
			Name:  "string and int and symbol",
			Input: "(print \"hi\" -7 foo!)",
			Want:  toks(del("("), sym("print"), str("\"hi\""), num("-7"), sym("foo!"), del(")")),
		},
		tc{
			Name:  "multiline arithmetic",
			Input: "(+ 1\n   (* 2 3))",
			Want:  toks(del("("), sym("+"), num("1"), del("("), sym("*"), num("2"), num("3"), del(")"), del(")")),
		},
	)

	a.Run("token values preserve source", func(a assert.A) {
		got := mustTokenize(a, "(foo 123 \"bar\" 'baz)")
		wantVals := []string{"(", "foo", "123", "\"bar\"", "'", "baz", ")"}
		a.Equal(len(got), len(wantVals))
		for i, w := range wantVals {
			a.Equal(got[i].Value, w)
		}
	})

	a.Run("token types in order", func(a assert.A) {
		got := mustTokenize(a, "(foo 123 \"bar\" 'baz)")
		wantTypes := []Kind{KindDelimiter, KindSymbol, KindInt, KindDQS, KindSigil, KindSymbol, KindDelimiter}
		a.Equal(len(got), len(wantTypes))
		for i, w := range wantTypes {
			a.Equal(got[i].Type, w)
		}
	})
}
