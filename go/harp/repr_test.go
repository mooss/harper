//revive:disable:line-length-limit
package harp_test

import (
	"testing"

	"github.com/mooss/bagend/go/assert"
	"github.com/mooss/harper/go/harp"
	"github.com/mooss/harper/go/harp/typ"
)

type ca = assert.Case[any, string]

func TestReprBool(t *testing.T) {
	a := assert.New(t)
	a.RunCases("bool", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "true", Input: typ.Bool(true), Want: "true"},
		ca{Name: "false", Input: typ.Bool(false), Want: "false"},
	)
}

func TestReprInt(t *testing.T) {
	a := assert.New(t)
	a.RunCases("int", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "zero", Input: typ.Int(0), Want: "0"},
		ca{Name: "one", Input: typ.Int(1), Want: "1"},
		ca{Name: "negative one", Input: typ.Int(-1), Want: "-1"},
		ca{Name: "positive", Input: typ.Int(42), Want: "42"},
		ca{Name: "negative", Input: typ.Int(-42), Want: "-42"},
		ca{Name: "large positive", Input: typ.Int(1000000), Want: "1000000"},
		ca{Name: "large negative", Input: typ.Int(-1000000), Want: "-1000000"},
	)
}

func TestReprKeyword(t *testing.T) {
	a := assert.New(t)
	a.RunCases("keyword", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "simple", Input: typ.Keyword("foo"), Want: ":foo"},
		ca{Name: "single char", Input: typ.Keyword("a"), Want: ":a"},
		ca{Name: "with dash", Input: typ.Keyword("hello-world"), Want: ":hello-world"},
		ca{Name: "numeric", Input: typ.Keyword("123"), Want: ":123"},
		ca{Name: "empty", Input: typ.Keyword(""), Want: ":"},
		ca{Name: "with slash", Input: typ.Keyword("foo/bar"), Want: ":foo/bar"},
	)
}

func TestReprSymbol(t *testing.T) {
	a := assert.New(t)
	a.RunCases("symbol", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "simple", Input: typ.Symbol("foo"), Want: "foo"},
		ca{Name: "single char", Input: typ.Symbol("a"), Want: "a"},
		ca{Name: "with dash", Input: typ.Symbol("foo-bar"), Want: "foo-bar"},
		ca{Name: "plus", Input: typ.Symbol("+"), Want: "+"},
		ca{Name: "star", Input: typ.Symbol("*"), Want: "*"},
		ca{Name: "def", Input: typ.Symbol("def"), Want: "def"},
		ca{Name: "nil symbol", Input: typ.Symbol("nil"), Want: "nil"},
		ca{Name: "with question", Input: typ.Symbol("foo?"), Want: "foo?"},
		ca{Name: "with slash", Input: typ.Symbol("foo/bar"), Want: "foo/bar"},
	)
}

func TestReprString(t *testing.T) {
	a := assert.New(t)
	a.RunCases("string", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "empty", Input: typ.String(""), Want: `""`},
		ca{Name: "simple", Input: typ.String("hello"), Want: `"hello"`},
		ca{Name: "with space", Input: typ.String("hello world"), Want: `"hello world"`},
		ca{Name: "with escaped quote", Input: typ.String(`hello "world"`), Want: `"hello \"world\""`},
		ca{Name: "with backslash", Input: typ.String(`a\b`), Want: `"a\\b"`},
		ca{Name: "with newline", Input: typ.String("a\nb"), Want: `"a\nb"`},
		ca{Name: "single quote kept", Input: typ.String("it's"), Want: `"it's"`},
		ca{Name: "numeric string", Input: typ.String("123"), Want: `"123"`},
	)
}

func TestReprList(t *testing.T) {
	a := assert.New(t)
	a.RunCases("list", func(a assert.A, c ca) {
		got, err := harp.Repr(c.Input)
		a.NoError(err)
		a.Equal(got, c.Want)
	},
		ca{Name: "empty", Input: typ.List{}, Want: "()"},
		ca{Name: "nil list", Input: typ.List(nil), Want: "()"},
		ca{Name: "single int", Input: typ.List{typ.Int(1)}, Want: "(1)"},
		ca{Name: "multiple ints", Input: typ.List{typ.Int(1), typ.Int(2), typ.Int(3)}, Want: "(1 2 3)"},
		ca{Name: "negative ints", Input: typ.List{typ.Int(-1), typ.Int(-2)}, Want: "(-1 -2)"},
		ca{Name: "bools", Input: typ.List{typ.Bool(true), typ.Bool(false)}, Want: "(true false)"},
		ca{Name: "keywords", Input: typ.List{typ.Keyword("a"), typ.Keyword("b")}, Want: "(:a :b)"},
		ca{Name: "symbols", Input: typ.List{typ.Symbol("a"), typ.Symbol("b")}, Want: "(a b)"},
		ca{Name: "strings", Input: typ.List{typ.String("a"), typ.String("b")}, Want: `("a" "b")`},
		ca{Name: "mixed", Input: typ.List{typ.Int(1), typ.Bool(true), typ.Keyword("kw"), typ.String("str"), typ.Symbol("sym")}, Want: `(1 true :kw "str" sym)`},
		ca{Name: "nested empty", Input: typ.List{typ.List{}}, Want: "(())"},
		ca{Name: "nested single", Input: typ.List{typ.List{typ.Int(1)}}, Want: "((1))"},
		ca{Name: "nested multiple", Input: typ.List{typ.Int(1), typ.List{typ.Int(2), typ.Int(3)}}, Want: "(1 (2 3))"},
		ca{Name: "deeply nested", Input: typ.List{typ.List{typ.List{typ.Int(1)}}}, Want: "(((1)))"},
		ca{Name: "two nested", Input: typ.List{typ.List{typ.Int(1)}, typ.List{typ.Int(2)}}, Want: "((1) (2))"},
		ca{Name: "empty nested twice", Input: typ.List{typ.List{}, typ.List{}}, Want: "(() ())"},
		ca{Name: "string with space in list", Input: typ.List{typ.String("hello world"), typ.Int(1)}, Want: `("hello world" 1)`},
		ca{Name: "keyword and symbol", Input: typ.List{typ.Keyword("foo"), typ.Symbol("foo")}, Want: "(:foo foo)"},
		ca{Name: "list starting with list", Input: typ.List{typ.List{typ.Symbol("a")}, typ.Symbol("b"), typ.Int(3)}, Want: "((a) b 3)"},
	)
}

func TestReprErrors(t *testing.T) {
	a := assert.New(t)
	a.RunCases("errors", func(a assert.A, c ca) {
		_, err := harp.Repr(c.Input)
		a.MustError(err)
		a.ErrorContains(err, "repr:")
	},
		ca{Name: "nil", Input: nil},
		ca{Name: "go string", Input: "hello"},
		ca{Name: "go int", Input: 42},
		ca{Name: "go bool", Input: true},
		ca{Name: "go float64", Input: float64(1.5)},
		ca{Name: "go float32", Input: float32(1.5)},
		ca{Name: "go slice", Input: []any{typ.Int(1)}},
		ca{Name: "go map", Input: map[string]any{}},
		ca{Name: "go struct", Input: struct{}{}},
		ca{Name: "list with nil", Input: typ.List{nil}},
		ca{Name: "list with go int", Input: typ.List{42}},
		ca{Name: "list with go string", Input: typ.List{"hello"}},
		ca{Name: "list with float", Input: typ.List{typ.Int(1), float64(2.5)}},
		ca{Name: "nested error", Input: typ.List{typ.List{typ.Int(1), nil}}},
		ca{Name: "deeply nested error", Input: typ.List{typ.Symbol("a"), typ.List{typ.Keyword("b"), typ.List{struct{}{}}}}},
	)
}

func TestReprErrorMessage(t *testing.T) {
	a := assert.New(t)
	a.Run("unsupported mentions type", func(a assert.A) {
		_, err := harp.Repr(float64(1.5))
		a.MustError(err)
		a.ErrorContains(err, "repr: unsupported type")
	})
	a.Run("nested error propagates", func(a assert.A) {
		_, err := harp.Repr(typ.List{typ.Int(1), typ.List{float64(2)}})
		a.MustError(err)
		a.ErrorContains(err, "repr:")
	})
}
