package harp

import (
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/harper/go/core"
	"github.com/mooss/harper/go/lex"
)

type Token = core.Token[Kind]

// lexer is Harp's lexer.
var lexer = lex.Dialect[Kind]{
	Whitespace:  lex.InSet(" \t\n\r"),
	LineComment: ';',
	Rules: []lex.Rule[Kind]{
		lex.AnyOf("()", KindDelimiter),
		lex.Quoted('"', '\\', KindDQS),
		lex.FirstPrefix(KindSigil, "'"),
		lex.SignedInt(lex.IsDigitASCII, boundary, KindInt),
		lex.Symbol(symbolStart, symbolRest, boundary, KindSymbol),
	},
}

// Tokenize splits a Harp program into its atomic parts.
func Tokenize(program string) ([]Token, error) {
	return lie.CollectE(lex.TokenIterator(program, lexer.Scanner()))
}

//revive:disable:line-length-limit
func symbolStart(r rune) bool {
	return strings.ContainsRune("!&*+-./:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", r)
}

//revive:disable:line-length-limit
func symbolRest(r rune) bool {
	return strings.ContainsRune("!&*+-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", r)
}

func boundary(r rune) bool {
	return strings.ContainsRune("() \t\n\r;", r)
}
