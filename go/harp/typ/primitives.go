package typ

// Bool is a Harp boolean.
type Bool bool

// Int is a Harp integer.
type Int int

// Keyword is a Harp keyword, starting with a colon.
type Keyword string

// List is a Harp list containing arbitrary Harp values.
type List []any

// String is a Harp double-quoted string.
type String string

// Symbol is a Harp symbol taking a value in an environment.
type Symbol string
