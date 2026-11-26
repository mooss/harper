package harp

import (
	"math/bits"
	"strconv"

	"github.com/mooss/bagend/go/damn"
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
	"github.com/mooss/harper/go/rad"
)

// reader is Harp's reader, the set of rules that specify how to build an AST from raw Harp tokens.
var reader = rad.Dialect[Kind, any]{Rules: []rad.Rule[Kind, any]{
	rad.Atom(KindInt, mkInt),
	rad.Atom(KindDQS, mkStr),
	rad.Atom(KindSymbol, mkSym),
	rad.Sequence(KindDelimiter, "(", ")", mkLst),
	rad.Prefix(KindSigil, "'", apply("quote")),
}}

var spec = damn.Must(reader.Compile())

func ReadStr(input string) (*rad.Reader[Kind, any], error) {
	toks, err := Tokenize(input)
	return spec.NewReader(toks), err
}

func ParseOne(input string) (any, error) {
	toks, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	return spec.ParseOne(toks)
}

func ParseAll(input string) ([]any, error) {
	toks, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	return spec.ParseAll(toks)
}

//////////////
// Builders //

// mkLst wraps a sequence of elements into a list.
func mkLst(elements []any) (any, error) { return typ.List(elements), nil }

// mkInt constructs an int from the given input token.
func mkInt(tok Token) (any, error) {
	i, err := strconv.ParseInt(tok.Value, 10, bits.UintSize)
	return typ.Int(i), err
}

// mkStr constructs a string from the given input token.
func mkStr(tok Token) (any, error) {
	s, err := lib.Unquote(tok.Value)
	return typ.String(s), err
}

// mkSym constructs a symbol from the given input token, handling true, false and keywords as
// special cases.
func mkSym(tok Token) (any, error) {
	switch tok.Value {
	case "true":
		return typ.Bool(true), nil
	case "false":
		return typ.Bool(false), nil
	}
	if len(tok.Value) > 0 && tok.Value[0] == ':' {
		return typ.Keyword(tok.Value[1:]), nil
	}
	return typ.Symbol(tok.Value), nil
}

// apply wraps the next form within a list application starting with the given symbol.
func apply(sym string) func(any) (any, error) {
	return func(n any) (any, error) {
		return typ.List{typ.Symbol(sym), n}, nil
	}
}
