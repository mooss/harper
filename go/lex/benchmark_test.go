//nolint:gosec
package lex

import (
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

// lispBenchKind is the token type used by the benchmark Lisp dialect.
type lispBenchKind int

const (
	lispBenchLParen lispBenchKind = iota
	lispBenchRParen
	lispBenchQuote
	lispBenchNumber
	lispBenchSymbol
	lispBenchString
)

func lispBenchIsBoundary(r rune) bool {
	return unicode.IsSpace(r) || strings.ContainsRune("()'`\",;", r)
}

func lispBenchIsSymbolStart(r rune) bool {
	return unicode.IsLetter(r) || strings.ContainsRune("!$%&*+-./:<=>?@^_~", r)
}

func lispBenchIsSymbolRest(r rune) bool {
	return unicode.IsLetter(r) || IsDigitASCII(r) || strings.ContainsRune("!$%&*+-./:<=>?@^_~", r)
}

// lispBenchDialect is a plausible Lisp dialect exercising every lexer
// primitive: single-char tokens, quoted strings with escapes, signed
// integers with boundary checks, symbols with start/rest predicates,
// whitespace skipping and line comments.
var lispBenchDialect = Dialect[lispBenchKind]{
	Whitespace:  unicode.IsSpace,
	LineComment: ';',
	Rules: []Rule[lispBenchKind]{
		AnyOf("(", lispBenchLParen),
		AnyOf(")", lispBenchRParen),
		AnyOf("'`,", lispBenchQuote),
		Quoted('"', '\\', InvalidTokenBoundary, lispBenchString),
		SignedInt(IsDigitASCII, lispBenchIsBoundary, lispBenchNumber),
		Symbol(lispBenchIsSymbolStart, lispBenchIsSymbolRest, lispBenchIsBoundary, lispBenchSymbol),
	},
}

var lispBenchHeads = []string{
	"defun", "let", "lambda", "if", "cond", "car", "cdr", "cons",
	"list", "map", "apply", "funcall", "quote", "quasiquote",
	"+", "-", "*", "/", "=", "<", ">", "<=", ">=",
	"foo", "bar", "baz", "x", "y", "n", "acc", "items",
}

var lispBenchSeparators = []string{" ", " ", " ", "  ", "\n", "\n  ", "\t", "\r\n"}

var lispBenchCommentWords = []string{
	"hello", "world", "loop", "recurse", "fixme", "todo", "note", "check",
}

func lispBenchRandomSymbol(rng *rand.Rand) string {
	base := lispBenchHeads[rng.Intn(len(lispBenchHeads))]
	n := rng.Intn(4)
	if n == 0 {
		return base
	}
	var sb strings.Builder
	sb.WriteString(base)
	for range n {
		if rng.Intn(2) == 0 {
			sb.WriteByte('-')
		}
		c := rng.Intn(36)
		if c < 26 {
			sb.WriteByte(byte('a' + c))
		} else {
			sb.WriteByte(byte('0' + c - 26))
		}
	}
	s := sb.String()
	// Avoid a signed-int prefix (e.g. "-1j"): SignedInt would match the
	// leading sign+digits then fail on the invalid token boundary
	// instead of falling through to Symbol. Break the prefix so the
	// whole text lexes as a symbol.
	if len(s) >= 2 && (s[0] == '+' || s[0] == '-') && s[1] >= '0' && s[1] <= '9' {
		return s[:1] + "a" + s[1:]
	}
	return s
}

func lispBenchRandomInt(rng *rand.Rand) string {
	var sb strings.Builder
	if rng.Intn(8) == 0 {
		if rng.Intn(2) == 0 {
			sb.WriteByte('+')
		} else {
			sb.WriteByte('-')
		}
	}
	digits := 1 + rng.Intn(6)
	for range digits {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	return sb.String()
}

func lispBenchRandomString(rng *rand.Rand) string {
	var sb strings.Builder
	sb.WriteByte('"')
	n := rng.Intn(32)
	for range n {
		switch rng.Intn(16) {
		case 0:
			sb.WriteString("\\\"")
		case 1:
			sb.WriteString("\\\\")
		case 2:
			sb.WriteString("\\n")
		case 3:
			sb.WriteByte(' ')
		default:
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

func lispBenchRandomAtom(rng *rand.Rand) string {
	switch r := rng.Intn(100); {
	case r < 55:
		return lispBenchRandomSymbol(rng)
	case r < 80:
		return lispBenchRandomInt(rng)
	case r < 90:
		return lispBenchRandomString(rng)
	default:
		return "'" + lispBenchRandomSymbol(rng)
	}
}

func lispBenchRandomExpr(rng *rand.Rand, depth int) string {
	if depth <= 0 || rng.Intn(4) == 0 {
		return lispBenchRandomAtom(rng)
	}
	var sb strings.Builder
	sb.WriteByte('(')
	sb.WriteString(lispBenchHeads[rng.Intn(len(lispBenchHeads))])
	elems := rng.Intn(5)
	for range elems {
		sb.WriteByte(' ')
		if rng.Intn(2) == 0 {
			sb.WriteString(lispBenchRandomExpr(rng, depth-1))
		} else {
			sb.WriteString(lispBenchRandomAtom(rng))
		}
	}
	sb.WriteByte(')')
	return sb.String()
}

// GeneratePseudorandomLisp generates deterministic, plausible pseudorandom
// Lisp code of approximately approxBytes bytes using seed for reproducibility.
// Output mixes nested forms, atoms, strings with escapes, varied whitespace
// and line comments so benchmarks exercise all lexer paths.
func GeneratePseudorandomLisp(seed int64, approxBytes int) string {
	rng := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.Grow(approxBytes + 1024)
	for sb.Len() < approxBytes {
		if rng.Intn(10) == 0 {
			sb.WriteString(";")
			for i := 0; i < 1+rng.Intn(6); i++ {
				sb.WriteByte(' ')
				sb.WriteString(lispBenchCommentWords[rng.Intn(len(lispBenchCommentWords))])
			}
			sb.WriteByte('\n')
			continue
		}
		sb.WriteString(lispBenchRandomExpr(rng, 5))
		sb.WriteString(lispBenchSeparators[rng.Intn(len(lispBenchSeparators))])
		if rng.Intn(20) == 0 {
			sb.WriteString("; trailing comment\n")
		}
	}
	return sb.String()
}

// HandCraftedLispInputs is a set of hand-crafted sequences covering edge cases:
// empty input, whitespace, single tokens, classic forms, escapes, comments,
// deep nesting, very long symbols and dense numeric/string runs.
var HandCraftedLispInputs = []struct {
	Name  string
	Input string
}{
	{"Empty", ""},
	{"WhitespaceOnly", "  \n\t  ; just a comment\n   \n"},
	{"SingleSymbol", "foobar"},
	{"SimpleAdd", "(+ 1 2)"},
	{"FibDefun", "(defun fib (n) (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2))))"},
	{"LetLambda", "(let ((x 10) (y -20)) ((lambda (a b) (+ a b)) x y))"},
	{"QuoteQuasiquote", "'(1 2 foo) `(a ,b ,@c) 'bar"},
	{"StringsEscapes", "(display \"hello \\\"world\\\"\\n\") \"back\\\\slash\" \"\" \"plain\""},
	{"Comments", "; leading comment\n(defun foo (x) ; trailing\n (+ x 1)) ; end\n"},
	{"DeepNesting", strings.Repeat("(", 256) + "x" + strings.Repeat(")", 256)},
	{"LongSymbol", "(foo " + strings.Repeat("a", 4096) + ")"},
	{"ManyIntegers", strings.Repeat("1234567890 ", 2000)},
	{"ManyStrings", strings.Repeat("\"hello \\\"world\\\"\" ", 500)},
	{"WhitespaceHeavy", "(  foo\n\t  bar   \n\n ; comment\n  baz  )"},
}

var (
	lispBenchSynthetic1K   = GeneratePseudorandomLisp(1, 1024)
	lispBenchSynthetic64K  = GeneratePseudorandomLisp(2, 64*1024)
	lispBenchSynthetic512K = GeneratePseudorandomLisp(3, 512*1024)
)

// lispBenchSink prevents the compiler from eliminating the scan loop.
var lispBenchSink int

func lispBenchDrain(b *testing.B, input string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	scan := lispBenchDialect.Scanner()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seq, errFn := TokenIterator(input, scan)
		n := 0
		for range seq {
			n++
		}
		if err := errFn(); err != nil {
			b.Fatalf("scan failed: %v", err)
		}
		lispBenchSink += n
	}
}

func BenchmarkLexerSynthetic1K(b *testing.B) {
	lispBenchDrain(b, lispBenchSynthetic1K)
}

func BenchmarkLexerSynthetic64K(b *testing.B) {
	lispBenchDrain(b, lispBenchSynthetic64K)
}

func BenchmarkLexerSynthetic512K(b *testing.B) {
	lispBenchDrain(b, lispBenchSynthetic512K)
}

func BenchmarkLexerHandCrafted(b *testing.B) {
	for _, tc := range HandCraftedLispInputs {
		b.Run(tc.Name, func(b *testing.B) {
			lispBenchDrain(b, tc.Input)
		})
	}
}

func BenchmarkLexerThroughputBySize(b *testing.B) {
	sizes := []struct {
		Name string
		Size int
	}{
		{"4K", 4 * 1024},
		{"64K", 64 * 1024},
		{"256K", 256 * 1024},
	}
	for _, s := range sizes {
		input := GeneratePseudorandomLisp(42, s.Size)
		b.Run(s.Name, func(b *testing.B) {
			lispBenchDrain(b, input)
		})
	}
}
