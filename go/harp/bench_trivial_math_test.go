package harp_test

import (
	"testing"

	"github.com/mooss/harper/go/harp"
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib/fun"
	"github.com/mooss/harper/go/lib/math"
)

var sink any

// Baseline: raw generic function, no interpreter machinery.
func BenchmarkPlusDirect(b *testing.B) {
	plus := math.Plus[typ.Int]
	b.ReportAllocs()
	for b.Loop() {
		r, _ := plus(1, 2)
		sink = r
	}
}

// Current representation path: Variadic adapter + boxing.
func BenchmarkPlusBuiltin(b *testing.B) {
	plus := fun.Variadic[typ.Builtin](math.Plus[typ.Int])
	args := []any{typ.Int(1), typ.Int(2)}
	b.ReportAllocs()
	for b.Loop() {
		r, _ := plus(args...)
		sink = r
	}
}

// Comparison operators go through fun.II.
func BenchmarkLtBuiltin(b *testing.B) {
	lt := fun.II[typ.Builtin](math.Lt[typ.Int, typ.Bool])
	args := []any{typ.Int(1), typ.Int(2)}
	b.ReportAllocs()
	for b.Loop() {
		r, _ := lt(args...)
		sink = r
	}
}

func BenchmarkEnvGetRoot(b *testing.B) {
	env := harp.MkEnv()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = env.Get("+")
	}
}

func BenchmarkEnvGetNested(b *testing.B) {
	env := harp.MkEnv().Nest().Nest()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = env.Get("+")
	}
}

// Full path: env lookup + builtin call.
func BenchmarkPlusViaEnv(b *testing.B) {
	env := harp.MkEnv()
	val, _ := env.Need("+")
	plus := val.(typ.Builtin)
	args := []any{typ.Int(1), typ.Int(2)}
	b.ReportAllocs()
	for b.Loop() {
		r, _ := plus(args...)
		sink = r
	}
}
