package lib

import (
	"fmt"
	"iter"
)

///////////////////
// Fixed lengths //

// As asserts that val is of type T.
func As[T any](val any) (T, error) {
	res, isT := val.(T)
	if !isT {
		return res, Expected(TypeOf(res), val)
	}
	return res, nil
}

// As0 asserts that args contains exactly zero arguments.
func As0(args []any) (err error) {
	return ExpectNArgs(0, args)
}

// As1 asserts that args contains exactly one argument of type A.
func As1[A any](args []any) (a A, err error) {
	if err = ExpectNArgs(1, args); err != nil {
		return
	}

	a, err = As[A](args[0])
	if err != nil {
		return
	}

	return a, nil
}

// As2 asserts that args contains exactly two arguments of types A and B.
func As2[A, B any](args []any) (a A, b B, err error) {
	if err = ExpectNArgs(2, args); err != nil {
		return
	}

	a, err = As[A](args[0])
	if err != nil {
		return
	}
	b, err = As[B](args[1])
	if err != nil {
		return
	}

	return a, b, nil
}

// As3 asserts that args contains exactly three arguments of types A, B and C.
func As3[A, B, C any](args []any) (a A, b B, c C, err error) {
	if err = ExpectNArgs(3, args); err != nil {
		return
	}

	a, err = As[A](args[0])
	if err != nil {
		return
	}
	b, err = As[B](args[1])
	if err != nil {
		return
	}
	c, err = As[C](args[2])
	if err != nil {
		return
	}

	return a, b, c, nil
}

/////////////////////
// Variable length //

// AsSeq returns a sequence of values from args, each asserted to be of type T.
// The sequence stops with an error on the first value that is not of type T.
func AsSeq[T any, U any](values ...U) iter.Seq2[T, error] {
	i := 0
	return func(yield func(T, error) bool) {
		for i < len(values) {
			val := values[i]
			typed, isT := any(val).(T)
			if !isT {
				yield(typed, Expected(TypeOf(typed), val))
				i = len(values)
				return
			}
			i++
			if !yield(typed, nil) {
				return
			}
		}
	}
}

///////////
// Other //

// TypeOf returns the type name of value, as formatted by fmt's %T verb.
func TypeOf(value any) string { return fmt.Sprintf("%T", value) }
