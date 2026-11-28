package harp_test

import (
	"testing"

	"github.com/mooss/harper/go/harp/env"
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib/fun"
	"github.com/mooss/harper/go/lib/math"
)

var sink any

// Baseline: raw generic binary function, no interpreter machinery.
func BenchmarkPlusDirect(b *testing.B) {
	plus := math.Add[typ.Int]
	b.ReportAllocs()
	for b.Loop() {
		r := plus(1, 2)
		sink = r
	}
}

// Current representation path: folding binary add.
func BenchmarkPlusBuiltin(b *testing.B) {
	plus := fun.Fold[typ.Builtin](math.Add[typ.Int], 1)
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
	envr := env.MkEnv()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = envr.Get("+")
	}
}

func BenchmarkEnvGetNested(b *testing.B) {
	envr := env.MkEnv().Nest(nil).Nest(nil)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = envr.Get("+")
	}
}

// Full path: env lookup + builtin call.
func BenchmarkPlusViaEnv(b *testing.B) {
	envr := env.MkEnv()
	val, _ := envr.Need("+")
	plus := val.(typ.Builtin) //nolint:revive
	args := []any{typ.Int(1), typ.Int(2)}
	b.ReportAllocs()
	for b.Loop() {
		r, _ := plus(args...)
		sink = r
	}
}
