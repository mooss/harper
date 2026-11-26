package harp

// Kind identifies a token's syntactic class.
type Kind int

const (
	// Lexer error, must stay to top so that the default value is invalid.
	KindError     Kind = iota
	KindDQS            // Double quoted string.
	KindInt            // Integer.
	KindDelimiter      // Balanced expression demimiter.
	KindSigil          // Special prefix wrapping a list application around what it qualifies.
	KindSymbol         // Raw symbol.
)

// String returns the short name of a token kind.
func (k Kind) String() string {
	switch k {
	case KindError:
		return "ERROR"
	case KindDQS:
		return "DQS"
	case KindInt:
		return "INT"
	case KindDelimiter:
		return "DELIMITER"
	case KindSigil:
		return "SIGIL"
	case KindSymbol:
		return "SYMBOL"
	default:
		return "UNKNOWN"
	}
}

// What returns a human-readable description of a token kind.
func (k Kind) What() string {
	switch k {
	case KindError:
		return "error"
	case KindDQS:
		return "double quoted string"
	case KindInt:
		return "integer"
	case KindDelimiter:
		return "delimiter"
	case KindSigil:
		return "sigil"
	case KindSymbol:
		return "symbol"
	default:
		return "unknown"
	}
}
