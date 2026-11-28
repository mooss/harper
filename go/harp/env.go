package harp

import (
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
	"github.com/mooss/harper/go/lib/fun"
	"github.com/mooss/harper/go/lib/math"
)

// ii converts a function with two arguments and a single result to typ.Builtin.
func ii[A, B, R any](fn func(A, B) R) typ.Builtin { return fun.II[typ.Builtin](fn) }

// MkEnv creates a default Harp environment.
func MkEnv() *typ.Env {
	return lib.NewEnv(map[string]any{
		"+": fun.Fold[typ.Builtin](math.Add[typ.Int], 1),
		"/": fun.Fold[typ.Builtin](math.Div[typ.Int], 2),
		"*": fun.Fold[typ.Builtin](math.Mul[typ.Int], 1),
		"-": fun.UnaryOrFold[typ.Builtin](math.Neg[typ.Int], math.Sub[typ.Int], 2),

		"<":  ii(math.Lt[typ.Int, typ.Bool]),
		"<=": ii(math.Leq[typ.Int, typ.Bool]),
		">":  ii(math.Gt[typ.Int, typ.Bool]),
		">=": ii(math.Geq[typ.Int, typ.Bool]),
	})
}
