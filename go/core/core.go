package core

// Token is a token of type T and its source text.
type Token[T any] struct {
	Type  T
	Value string
}
