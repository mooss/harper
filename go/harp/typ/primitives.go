package typ

// Bool is a Harp boolean.
type Bool bool

// Builtin is a Harp primitive function implemented in Go.
type Builtin func(...any) (any, error)

// Closure is a Harp function implemented in Harp.
type Closure func(...any) (any, error)

// Int is a Harp integer.
type Int int

// Keyword is a Harp keyword, starting with a colon.
type Keyword string

// List is a Harp list containing arbitrary Harp values.
type List []any

// String is a Harp double-quoted string.
type String string

// Special is a Harp special form, evaluated with its raw (unevaluated) arguments.
type Special func(*Env, ...any) (any, error)

// Symbol is a Harp symbol taking a value in an environment.
type Symbol string
