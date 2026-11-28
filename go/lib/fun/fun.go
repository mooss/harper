package fun

import "github.com/mooss/harper/go/lib"

type Signature = func(...any) (any, error)

// Fold adapts an associative binary function into a Signature that left-folds over its arguments
// without intermediate allocations.
// minargs is the minimum required amount of arguments, which should never be less than one or it
// will crash.
func Fold[S ~Signature, A any](fn func(A, A) A, minargs int) S {
	return func(args ...any) (any, error) {
		if err := lib.ExpectAtLeastNArgs(minargs, args); err != nil {
			return nil, err
		}

		acc, err := lib.As[A](args[0])
		if err != nil {
			return nil, err
		}

		for _, arg := range args[1:] {
			val, err := lib.As[A](arg)
			if err != nil {
				return nil, err
			}
			acc = fn(acc, val)
		}
		return acc, nil
	}
}

func UnaryOrFold[S ~Signature, A any](unary func(A) A, binary func(A, A) A, minargs int) S {
	fold := Fold[S](binary, minargs)
	single := I[S](unary)
	return func(args ...any) (any, error) {
		if len(args) == 1 {
			return single(args...)
		}
		return fold(args...)
	}
}

func I[S ~Signature, A, R any](fun func(A) R) S {
	return func(args ...any) (any, error) {
		a, err := lib.As1[A](args)
		if err != nil {
			return nil, err
		}

		return fun(a), nil
	}
}

func II[S ~Signature, A, B, R any](fun func(A, B) R) S {
	return func(args ...any) (any, error) {
		a, b, err := lib.As2[A, B](args)
		if err != nil {
			return nil, err
		}

		return fun(a, b), nil
	}
}

func Variadic[S ~Signature, A, R any](fun func(...A) (R, error)) S {
	return func(args ...any) (any, error) {
		typed, err := lib.AsSlice[A](args...)
		if err != nil {
			return nil, err
		}
		return fun(typed...)
	}
}
