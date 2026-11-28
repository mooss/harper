package math

import (
	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/bagend/go/fun/seq"
	"github.com/mooss/bagend/go/fun/seq2"
	"github.com/mooss/harper/go/lib"
)

func reduce[T any](fn func(T, T) T) lib.Fun[T, T] {
	return func(args ...T) (T, error) {
		if err := lib.ExpectAtLeastNArgs(1, args); err != nil {
			return fun.Zero[T](), err
		}

		typed, gerr := seq2.SplitErr(lib.AsSeq[T](args...))
		acc, _ := seq.First(typed)
		acc = seq.Fold(fn, acc, typed)
		return acc, gerr()
	}
}

type numbers interface{ ~int64 | ~float64 | ~int }

func div[T numbers](l, r T) T   { return l / r } // Will panic on r == 0.
func minus[T numbers](l, r T) T { return l - r }
func plus[T numbers](l, r T) T  { return l + r }
func times[T numbers](l, r T) T { return l * r }

// Divide reduces the arguments with the division operator.
// It requires at least two arguments and panics if any divisor is zero.
func Divide[N numbers](args ...N) (N, error) {
	if err := lib.ExpectAtLeastNArgs(2, args); err != nil {
		return fun.Zero[N](), err
	}

	return reduce(div[N])(args...)
}

// Minus subtracts the arguments from left to right.
// With a single argument, it returns its negation.
func Minus[N numbers](args ...N) (N, error) {
	switch len(args) {
	case 1:
		val, err := lib.As[N](args[0])
		return -val, err
	default:
		return reduce(minus[N])(args...)
	}
}

// Plus sums the arguments.
func Plus[N numbers](args ...N) (N, error) { return reduce(plus[N])(args...) }

// Times multiplies the arguments.
func Times[N numbers](args ...N) (N, error) { return reduce(times[N])(args...) }

// Lt reports whether l is strictly less than t.
func Lt[T numbers, B ~bool](l, t T) B { return l < t }

// Leq reports whether l is less than or equal to t.
func Leq[T numbers, B ~bool](l, t T) B { return l <= t }

// Gt reports whether l is strictly greater than t.
func Gt[T numbers, B ~bool](l, t T) B { return l > t }

// Geq reports whether l is greater than or equal to t.
func Geq[T numbers, B ~bool](l, t T) B { return l >= t }
