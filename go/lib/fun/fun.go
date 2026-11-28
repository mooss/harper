package fun

import "github.com/mooss/harper/go/lib"

type Signature = func(...any) (any, error)

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
