//revive:disable:cognitive-complexity,cyclomatic,function-length,line-length-limit
package harp_test

import (
	"errors"
	"io"
	"testing"

	"github.com/mooss/bagend/go/assert"
	"github.com/mooss/harper/go/harp"
	"github.com/mooss/harper/go/harp/typ"
)

type PO = assert.Case[string, any]
type PA = assert.Case[string, []any]

func checkOne(a assert.A, c PO) {
	a.Helper()
	got, err := harp.ParseOne(c.Input)
	if c.WantErr {
		a.MustError(err)
		return
	}
	a.NoError(err)
	a.Equal(got, c.Want)
}

func checkAll(a assert.A, c PA) {
	a.Helper()
	got, err := harp.ParseAll(c.Input)
	if c.WantErr {
		a.MustError(err)
		return
	}
	a.NoError(err)
	a.Equal(got, c.Want)
}

func TestParseOneAtoms(t *testing.T) {
	a := assert.New(t)
	a.RunCases("atoms", checkOne,
		PO{Name: "zero", Input: "0", Want: typ.Int(0)},
		PO{Name: "int", Input: "42", Want: typ.Int(42)},
		PO{Name: "int leading zeros", Input: "007", Want: typ.Int(7)},
		PO{Name: "int large", Input: "1234567890", Want: typ.Int(1234567890)},
		PO{Name: "int surrounding space", Input: "  7  ", Want: typ.Int(7)},
		PO{Name: "int newline", Input: "\n99\n", Want: typ.Int(99)},

		PO{Name: "empty string", Input: `""`, Want: typ.String("")},
		PO{Name: "simple string", Input: `"hello"`, Want: typ.String("hello")},
		PO{Name: "string with spaces", Input: `"a b c"`, Want: typ.String("a b c")},
		PO{Name: "string escaped quote", Input: `"a\"b"`, Want: typ.String(`a"b`)},
		PO{Name: "string escaped backslash", Input: `"a\\b"`, Want: typ.String(`a\b`)},
		PO{Name: "string escaped newline", Input: `"a\nb"`, Want: typ.String("a\nb")},
		PO{Name: "string escaped tab", Input: `"a\tb"`, WantErr: true},

		PO{Name: "symbol", Input: "foo", Want: typ.Symbol("foo")},
		PO{Name: "symbol with dash", Input: "foo-bar", Want: typ.Symbol("foo-bar")},
		PO{Name: "symbol plus", Input: "+", Want: typ.Symbol("+")},
		PO{Name: "symbol star", Input: "*", Want: typ.Symbol("*")},
		PO{Name: "symbol with digits", Input: "a1", Want: typ.Symbol("a1")},
		PO{Name: "symbol true prefix", Input: "trueX", Want: typ.Symbol("trueX")},
		PO{Name: "symbol false prefix", Input: "falsey", Want: typ.Symbol("falsey")},
		PO{Name: "symbol question", Input: "foo?", Want: typ.Symbol("foo?")},

		PO{Name: "true", Input: "true", Want: typ.Bool(true)},
		PO{Name: "false", Input: "false", Want: typ.Bool(false)},

		PO{Name: "keyword", Input: ":foo", Want: typ.Keyword("foo")},
		PO{Name: "keyword dash", Input: ":foo-bar", Want: typ.Keyword("foo-bar")},
		PO{Name: "keyword single char", Input: ":a", Want: typ.Keyword("a")},
	)
}

func TestParseOneLists(t *testing.T) {
	a := assert.New(t)
	a.RunCases("lists", checkOne,
		PO{Name: "empty", Input: "()", Want: typ.List(nil)},
		PO{Name: "empty spaced", Input: "( )", Want: typ.List(nil)},
		PO{Name: "empty multiline", Input: "(\n)", Want: typ.List(nil)},
		PO{Name: "single int", Input: "(1)", Want: typ.List{typ.Int(1)}},
		PO{Name: "two ints", Input: "(1 2)", Want: typ.List{typ.Int(1), typ.Int(2)}},
		PO{Name: "many ints", Input: "(1 2 3 4)", Want: typ.List{typ.Int(1), typ.Int(2), typ.Int(3), typ.Int(4)}},
		PO{Name: "nested empty", Input: "(())", Want: typ.List{typ.List(nil)}},
		PO{Name: "nested", Input: "((1) (2 3))", Want: typ.List{typ.List{typ.Int(1)}, typ.List{typ.Int(2), typ.Int(3)}}},
		PO{Name: "deep nest", Input: "(((((1)))))", Want: typ.List{typ.List{typ.List{typ.List{typ.List{typ.Int(1)}}}}}},
		PO{Name: "mixed", Input: `(1 "a" foo :k true false)`, Want: typ.List{typ.Int(1), typ.String("a"), typ.Symbol("foo"), typ.Keyword("k"), typ.Bool(true), typ.Bool(false)}},
		PO{Name: "whitespace", Input: "  (  1   2  )  ", Want: typ.List{typ.Int(1), typ.Int(2)}},
		PO{Name: "newlines", Input: "(1\n2\n3)", Want: typ.List{typ.Int(1), typ.Int(2), typ.Int(3)}},
		PO{Name: "list of lists strings", Input: `(("a" "b") ("c"))`, Want: typ.List{typ.List{typ.String("a"), typ.String("b")}, typ.List{typ.String("c")}}},
		PO{Name: "single symbol", Input: "(foo)", Want: typ.List{typ.Symbol("foo")}},
		PO{Name: "quote inside", Input: "('a 'b)", Want: typ.List{typ.List{typ.Symbol("quote"), typ.Symbol("a")}, typ.List{typ.Symbol("quote"), typ.Symbol("b")}}},
	)
}

func TestParseQuote(t *testing.T) {
	a := assert.New(t)
	a.RunCases("quote", checkOne,
		PO{Name: "quote sym", Input: "'a", Want: typ.List{typ.Symbol("quote"), typ.Symbol("a")}},
		PO{Name: "quote int", Input: "'1", Want: typ.List{typ.Symbol("quote"), typ.Int(1)}},
		PO{Name: "quote string", Input: `'"hi"`, Want: typ.List{typ.Symbol("quote"), typ.String("hi")}},
		PO{Name: "quote true", Input: "'true", Want: typ.List{typ.Symbol("quote"), typ.Bool(true)}},
		PO{Name: "quote keyword", Input: "':k", Want: typ.List{typ.Symbol("quote"), typ.Keyword("k")}},
		PO{Name: "quote empty", Input: "'()", Want: typ.List{typ.Symbol("quote"), typ.List(nil)}},
		PO{Name: "quote list", Input: "'(1 2)", Want: typ.List{typ.Symbol("quote"), typ.List{typ.Int(1), typ.Int(2)}}},
		PO{Name: "double quote", Input: "''a", Want: typ.List{typ.Symbol("quote"), typ.List{typ.Symbol("quote"), typ.Symbol("a")}}},
		PO{Name: "triple quote", Input: "'''a", Want: typ.List{typ.Symbol("quote"), typ.List{typ.Symbol("quote"), typ.List{typ.Symbol("quote"), typ.Symbol("a")}}}},
		PO{Name: "quote with space", Input: "' a", Want: typ.List{typ.Symbol("quote"), typ.Symbol("a")}},
		PO{Name: "quote list spaced", Input: "' (1 2)", Want: typ.List{typ.Symbol("quote"), typ.List{typ.Int(1), typ.Int(2)}}},
		PO{Name: "quote nested list", Input: "'((1) 2)", Want: typ.List{typ.Symbol("quote"), typ.List{typ.List{typ.Int(1)}, typ.Int(2)}}},
		PO{Name: "quote in list", Input: "('a)", Want: typ.List{typ.List{typ.Symbol("quote"), typ.Symbol("a")}}},
		PO{Name: "list quoted tail", Input: "(a 'b c)", Want: typ.List{typ.Symbol("a"), typ.List{typ.Symbol("quote"), typ.Symbol("b")}, typ.Symbol("c")}},
	)
}

func TestParseOneErrors(t *testing.T) {
	a := assert.New(t)
	a.RunCases("parse one errors", checkOne,
		PO{Name: "empty", Input: "", WantErr: true},
		PO{Name: "blank", Input: "   ", WantErr: true},
		PO{Name: "newline only", Input: "\n", WantErr: true},
		PO{Name: "unbalanced open", Input: "(", WantErr: true},
		PO{Name: "unbalanced body", Input: "(1 2", WantErr: true},
		PO{Name: "unbalanced nested", Input: "((1)", WantErr: true},
		PO{Name: "unexpected close", Input: ")", WantErr: true},
		PO{Name: "close in list?", Input: "(1 2))", WantErr: true},
		PO{Name: "extra close", Input: "((1)))", WantErr: true},
		PO{Name: "dangling quote", Input: "'", WantErr: true},
		PO{Name: "dangling quote space", Input: "'   ", WantErr: true},
		PO{Name: "dangling quote open", Input: "'(", WantErr: true},
		PO{Name: "quote unbalanced", Input: "'(1 2", WantErr: true},
		PO{Name: "leftover two forms", Input: "1 2", WantErr: true},
		PO{Name: "leftover list sym", Input: "(1) a", WantErr: true},
		PO{Name: "unterminated string", Input: `"hi`, WantErr: true},
		PO{Name: "single dquote", Input: `"`, WantErr: true},
	)
}

func TestParseAll(t *testing.T) {
	a := assert.New(t)
	a.RunCases("parse all", checkAll,
		PA{Name: "empty", Input: "", Want: nil},
		PA{Name: "blank", Input: "   \n\t  ", Want: nil},
		PA{Name: "one int", Input: "1", Want: []any{typ.Int(1)}},
		PA{Name: "two ints", Input: "1 2", Want: []any{typ.Int(1), typ.Int(2)}},
		PA{Name: "newlines", Input: "1\n2\n3", Want: []any{typ.Int(1), typ.Int(2), typ.Int(3)}},
		PA{Name: "mixed", Input: `foo 1 "a" :k true`, Want: []any{typ.Symbol("foo"), typ.Int(1), typ.String("a"), typ.Keyword("k"), typ.Bool(true)}},
		PA{Name: "two lists", Input: "(1) (2)", Want: []any{typ.List{typ.Int(1)}, typ.List{typ.Int(2)}}},
		PA{Name: "list and atom", Input: "(1 2) 3 'a", Want: []any{typ.List{typ.Int(1), typ.Int(2)}, typ.Int(3), typ.List{typ.Symbol("quote"), typ.Symbol("a")}}},
		PA{Name: "quotes", Input: "'a 'b", Want: []any{typ.List{typ.Symbol("quote"), typ.Symbol("a")}, typ.List{typ.Symbol("quote"), typ.Symbol("b")}}},
		PA{Name: "nested lists", Input: "((1) (2)) ((3))", Want: []any{typ.List{typ.List{typ.Int(1)}, typ.List{typ.Int(2)}}, typ.List{typ.List{typ.Int(3)}}}},
	)
}

func TestParseAllErrors(t *testing.T) {
	a := assert.New(t)
	a.RunCases("parse all errors", checkAll,
		PA{Name: "unbalanced", Input: "(", WantErr: true},
		PA{Name: "unbalanced trailing", Input: "1 (2", WantErr: true},
		PA{Name: "unexpected close", Input: ")", WantErr: true},
		PA{Name: "unexpected close after", Input: "1 )", WantErr: true},
		PA{Name: "dangling quote", Input: "1 '", WantErr: true},
		PA{Name: "unterminated string", Input: `"oops`, WantErr: true},
		PA{Name: "extra close", Input: "(1))", WantErr: true},
	)
}

func TestQuoteEqualsExplicit(t *testing.T) {
	a := assert.New(t)
	inputs := [][2]string{
		{"'a", "(quote a)"},
		{"'1", "(quote 1)"},
		{`'"x"`, `(quote "x")`},
		{"'(1 2)", "(quote (1 2))"},
		{"''a", "(quote (quote a))"},
	}
	for _, pair := range inputs {
		quoted, err1 := harp.ParseOne(pair[0])
		a.NoError(err1)
		explicit, err2 := harp.ParseOne(pair[1])
		a.NoError(err2)
		a.Equal(quoted, explicit)
	}
}

func TestReadStrReader(t *testing.T) {
	a := assert.New(t)

	a.Run("empty done", func(a assert.A) {
		r, err := harp.ReadStr("")
		a.NoError(err)
		if r == nil {
			a.Fatal("nil reader")
		}
		if !r.Done() {
			a.Fatal("expected Done")
		}
		if got := r.Rest(); len(got) != 0 {
			a.Fatalf("expected no rest, got %#v", got)
		}
		_, err = r.ReadForm()
		if !errors.Is(err, io.EOF) {
			a.Fatalf("expected EOF, got %v", err)
		}
	})

	a.Run("single form", func(a assert.A) {
		r := a.Must(harp.ReadStr("42"))
		if r.Done() {
			a.Fatal("unexpected Done")
		}
		v := a.Must(r.ReadForm())
		a.Equal(v, typ.Int(42))
		if !r.Done() {
			a.Fatal("expected Done after read")
		}
		_, err := r.ReadForm()
		if !errors.Is(err, io.EOF) {
			a.Fatalf("expected EOF, got %v", err)
		}
	})

	a.Run("rest shrinks", func(a assert.A) {
		r := a.Must(harp.ReadStr("1 2 3"))
		if l := len(r.Rest()); l != 3 {
			a.Fatalf("rest len = %d, want 3", l)
		}
		v1 := a.Must(r.ReadForm())
		a.Equal(v1, typ.Int(1))
		if l := len(r.Rest()); l != 2 {
			a.Fatalf("rest len = %d, want 2", l)
		}
		v2 := a.Must(r.ReadForm())
		a.Equal(v2, typ.Int(2))
		v3 := a.Must(r.ReadForm())
		a.Equal(v3, typ.Int(3))
		if !r.Done() {
			a.Fatal("expected Done")
		}
	})

	a.Run("peek next", func(a assert.A) {
		r := a.Must(harp.ReadStr("1 2"))
		p1 := r.Peek()
		if p1.Value != "1" {
			a.Fatalf("peek = %q, want %q", p1.Value, "1")
		}
		if r.Done() {
			a.Fatal("peek advanced input")
		}
		n1 := r.Next()
		a.Equal(n1, p1)
		p2 := r.Peek()
		if p2.Value != "2" {
			a.Fatalf("peek = %q, want %q", p2.Value, "2")
		}
		v := a.Must(r.ReadForm())
		a.Equal(v, typ.Int(2))
		if !r.Done() {
			a.Fatal("expected Done")
		}
	})

	a.Run("tokenize error", func(a assert.A) {
		_, err := harp.ReadStr(`"unterminated`)
		a.MustError(err)
	})

	a.Run("reader sees tokens", func(a assert.A) {
		r := a.Must(harp.ReadStr("(a 'b)"))
		v := a.Must(r.ReadForm())
		a.Equal(v, typ.List{typ.Symbol("a"), typ.List{typ.Symbol("quote"), typ.Symbol("b")}})
	})

	a.Run("error contains reader", func(a assert.A) {
		r := a.Must(harp.ReadStr(")"))
		_, err := r.ReadForm()
		a.MustError(err)
		a.ErrorContains(err, "reader")
		_, err = harp.ParseOne("(")
		a.ErrorContains(err, "reader")
		_, err = harp.ParseOne("'")
		a.ErrorContains(err, "reader")
	})
}

func TestParseOneVsParseAll(t *testing.T) {
	a := assert.New(t)
	a.Run("one equals first of all", func(a assert.A) {
		for _, input := range []string{"1", "foo", `"hi"`, "(1 2)", "'a", ":k", "true"} {
			one := a.Must(harp.ParseOne(input))
			all := a.Must(harp.ParseAll(input))
			if len(all) != 1 {
				a.Fatalf("ParseAll(%q) len = %d, want 1", input, len(all))
			}
			a.Equal(one, all[0])
		}
	})
}
