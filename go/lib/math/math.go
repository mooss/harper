package math

type numbers interface{ ~int64 | ~float64 | ~int }

func Add[T numbers](l, r T) T { return l + r }
func Div[T numbers](l, r T) T { return l / r } // Will panic on r == 0.
func Mul[T numbers](l, r T) T { return l * r }
func Sub[T numbers](l, r T) T { return l - r }
func Neg[T numbers](x T) T    { return -x }

// Lt reports whether l is strictly less than t.
func Lt[T numbers, B ~bool](l, t T) B { return l < t }

// Leq reports whether l is less than or equal to t.
func Leq[T numbers, B ~bool](l, t T) B { return l <= t }

// Gt reports whether l is strictly greater than t.
func Gt[T numbers, B ~bool](l, t T) B { return l > t }

// Geq reports whether l is greater than or equal to t.
func Geq[T numbers, B ~bool](l, t T) B { return l >= t }
