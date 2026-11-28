package spe

import (
	"errors"

	"github.com/mooss/harper/go/harp"
	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
)

// Def introduces a new binding in the environment.
func Def(envr *typ.Env, key typ.Symbol, value any) (any, error) {
	evald, err := harp.Eval(value, envr)
	if err != nil {
		return nil, err
	}

	envr.Set(string(key), evald)
	return evald, nil
}

// Let creates a temporary environment containing the given bindings and evaluates form in this new
// environment.
func Let(envr *typ.Env, bindings any, form any) (any, error) {
	bindlist, err := lib.As[typ.List](bindings)
	if err != nil {
		return nil, err
	}

	if len(bindlist)%2 != 0 {
		return nil, errors.New("odd-numbered let bindings")
	}

	letenv := envr.Nest(nil)
	for i := 0; i < len(bindlist); i += 2 {
		key, err := lib.As[typ.Symbol](bindlist[i])
		if err != nil {
			return nil, err
		}

		value, err := harp.Eval(bindlist[i+1], letenv)
		if err != nil {
			return nil, err
		}

		letenv.Set(string(key), value)
	}

	//TODO: find out how to do TCO with this.
	return harp.Eval(form, letenv)
}

func IE[A, R any](fn func(*typ.Env, A) (R, error)) typ.Special {
	return func(env *typ.Env, args ...any) (any, error) {
		a, err := lib.As1[A](args)
		if err != nil {
			return nil, err
		}

		return fn(env, a)
	}
}

func IIE[A, B, R any](fn func(*typ.Env, A, B) (R, error)) typ.Special {
	return func(env *typ.Env, args ...any) (any, error) {
		a, b, err := lib.As2[A, B](args)
		if err != nil {
			return nil, err
		}

		return fn(env, a, b)
	}
}
