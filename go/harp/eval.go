package harp

import (
	"errors"
	"fmt"

	"github.com/mooss/bagend/go/damn"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
)

// Eval evaluates a Harp value in the given environment, returning the result.
//
//revive:disable:cyclomatic
func Eval(ast any, envr *typ.Env) (any, error) {
	switch node := ast.(type) {
	case typ.Bool, typ.Int, typ.Keyword, typ.String:
		return node, nil
	case typ.Builtin, typ.Closure, typ.Special:
		return nil, fmt.Errorf("cannot evaluate %#q", damn.Must(Repr(node)))
	case typ.List:
		if len(node) == 0 {
			return nil, errors.New("cannot evaluate the empty list")
		}
		op, err := Eval(node[0], envr)
		if err != nil {
			return nil, err
		}

		switch op := op.(type) {
		case typ.Builtin:
			return evalFun(envr, op, node[1:])
		case typ.Closure:
			return evalFun(envr, op, node[1:])
		case typ.Special:
			return op(envr, node[1:]...)
		default:
			return nil, lib.Expected("callable", op)
		}
	case typ.Symbol:
		return envr.Need(string(node))
	default:
	}
	return nil, fmt.Errorf("cannot evaluate %T (unknown type)", ast)
}

func evalFun(envr *typ.Env, fn func(...any) (any, error), args []any) (any, error) {
	evald, err := EvalArgs(envr, args...)
	if err != nil {
		return nil, err
	}
	return fn(evald...)
}

// EvalArgs evaluates all its arguments in order, returning them as a slice.
// The first error encountered triggers a short-circuit and is returned as an error.
func EvalArgs(envr *typ.Env, args ...any) ([]any, error) {
	return lie.MapE(func(arg any) (any, error) { return Eval(arg, envr) }, args)
}
