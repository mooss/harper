package lib

import "fmt"

// Env is a lexical environment mapping names to values.
type Env[Value any] struct {
	outer *Env[Value]
	data  map[string]Value
}

// NewEnv returns a root environment containing the given associations (nil will create a valid
// env).
func NewEnv[V any](assoc map[string]V) *Env[V] {
	if assoc == nil {
		assoc = map[string]V{}
	}
	return &Env[V]{outer: nil, data: assoc}
}

// Get looks up a symbol and returns false if it is not found.
func (e *Env[V]) Get(symbol string) (V, bool) {
	for env := e; env != nil; env = env.outer {
		if val, ok := env.data[symbol]; ok {
			return val, true
		}
	}
	var zero V
	return zero, false
}

// Need looks up a symbol and return an error if it is not found.
func (e *Env[V]) Need(symbol string) (V, error) {
	res, found := e.Get(symbol)
	if !found {
		return res, fmt.Errorf("symbol %#q not found", symbol)
	}
	return res, nil
}

// Set binds name to val in the current environment.
func (e *Env[V]) Set(name string, val V) {
	e.data[name] = val
}

// Nest returns a new empty child environment for a nested scope.
func (e *Env[V]) Nest(assoc map[string]V) *Env[V] {
	if assoc == nil {
		assoc = map[string]V{}
	}
	return &Env[V]{outer: e, data: assoc}
}
