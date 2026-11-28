package harp

import (
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
	"github.com/mooss/harper/go/lib/fun"
	"github.com/mooss/harper/go/lib/math"
)

// builtin converts an arbitrary variadic function to typ.Builtin.
func builtin[A, R any](fn func(...A) (R, error)) typ.Builtin {
	return fun.Variadic[typ.Builtin](fn)
}

// ii converts a function with two arguments and a single result to typ.Builtin.
func ii[A, B, R any](fn func(A, B) R) typ.Builtin { return fun.II[typ.Builtin](fn) }

// MkEnv creates a default Harp environment.
func MkEnv() *typ.Env {
	return lib.NewEnv(map[string]any{
		"/": builtin(math.Divide[typ.Int]),
		"-": builtin(math.Minus[typ.Int]),
		"+": builtin(math.Plus[typ.Int]),
		"*": builtin(math.Times[typ.Int]),

		"<":  ii(math.Lt[typ.Int, typ.Bool]),
		"<=": ii(math.Leq[typ.Int, typ.Bool]),
		">":  ii(math.Gt[typ.Int, typ.Bool]),
		">=": ii(math.Geq[typ.Int, typ.Bool]),
	})
}
