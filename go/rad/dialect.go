package rad

import (
	"errors"
	"fmt"

	"github.com/mooss/harper/go/core"
)

type ruleOp int

const (
	opAtom ruleOp = iota
	opSequence
	opTable
	opPrefix
)

type Rule[K comparable, V any] struct {
	op     ruleOp
	kind   K
	open   string
	close  string
	sigil  string
	arity  int
	atomFn func(core.Token[K]) (V, error)
	build  func([]V) (V, error)
	expand func(V) (V, error)
}

type Dialect[K comparable, V any] struct {
	Rules []Rule[K, V]
}

type openSpec[K comparable, V any] struct {
	close   string
	kind    K
	isTable bool
	arity   int
	build   func([]V) (V, error)
}

type token[Kind comparable] struct {
	kind  Kind
	value string
}

type Spec[K comparable, V any] struct {
	atoms       map[K]func(core.Token[K]) (V, error)
	opens       map[token[K]]openSpec[K, V]
	prefixes    map[token[K]]func(V) (V, error)
	delimKinds  map[K]bool
	prefixKinds map[K]bool
	closes      map[string]bool
}

// Atom defines an atomic value rule for kind using fn to convert a token.
func Atom[K comparable, V any](kind K, fn func(core.Token[K]) (V, error)) Rule[K, V] {
	return Rule[K, V]{op: opAtom, kind: kind, atomFn: fn}
}

// Sequence defines a variable-length delimited sequence rule.
func Sequence[K comparable, V any](
	kind K, open, until string, build func([]V) (V, error),
) Rule[K, V] {
	return Rule[K, V]{op: opSequence, kind: kind, open: open, close: until, build: build}
}

// Table defines a fixed-arity delimited table rule.
func Table[K comparable, V any](
	kind K, open, until string, arity int, build func([]V) (V, error),
) Rule[K, V] {
	return Rule[K, V]{op: opTable, kind: kind, open: open, close: until, arity: arity, build: build}
}

// Prefix defines a sigil prefix rule using expand to transform the operand.
func Prefix[K comparable, V any](kind K, sigil string, expand func(V) (V, error)) Rule[K, V] {
	return Rule[K, V]{op: opPrefix, kind: kind, sigil: sigil, expand: expand}
}

// Compile validates the dialect rules and builds a Spec for parsing.
//
//revive:disable:cognitive-complexity,cyclomatic,line-length-limit
func (d Dialect[K, V]) Compile() (*Spec[K, V], error) {
	res := &Spec[K, V]{
		atoms:       make(map[K]func(core.Token[K]) (V, error)),
		opens:       make(map[token[K]]openSpec[K, V]),
		prefixes:    make(map[token[K]]func(V) (V, error)),
		delimKinds:  make(map[K]bool),
		prefixKinds: make(map[K]bool),
		closes:      make(map[string]bool),
	}
	for _, rule := range d.Rules {
		switch rule.op {
		case opAtom:
			if _, ok := res.atoms[rule.kind]; ok {
				return nil, fmt.Errorf("rad: duplicate atom kind %v", rule.kind)
			}
			res.atoms[rule.kind] = rule.atomFn
		case opSequence:
			if rule.open == "" {
				return nil, errors.New("rad: empty open")
			}
			if rule.close == "" {
				return nil, errors.New("rad: empty close")
			}
			key := token[K]{rule.kind, rule.open}
			if _, ok := res.opens[key]; ok {
				return nil, fmt.Errorf("rad: duplicate open (%v, %q)", rule.kind, rule.open)
			}
			res.opens[key] = openSpec[K, V]{close: rule.close, kind: rule.kind,
				isTable: false, arity: 0, build: rule.build}
			res.delimKinds[rule.kind] = true
			res.closes[rule.close] = true
		case opTable:
			if rule.open == "" {
				return nil, errors.New("rad: empty open")
			}
			if rule.close == "" {
				return nil, errors.New("rad: empty close")
			}
			if rule.arity <= 0 {
				return nil, fmt.Errorf("rad: invalid arity %d", rule.arity)
			}
			key := token[K]{rule.kind, rule.open}
			if _, ok := res.opens[key]; ok {
				return nil, fmt.Errorf("rad: duplicate open (%v, %q)", rule.kind, rule.open)
			}
			res.opens[key] = openSpec[K, V]{close: rule.close, kind: rule.kind,
				isTable: true, arity: rule.arity, build: rule.build}
			res.delimKinds[rule.kind] = true
			res.closes[rule.close] = true
		case opPrefix:
			if rule.sigil == "" {
				return nil, errors.New("rad: empty sigil")
			}
			key := token[K]{rule.kind, rule.sigil}
			if _, ok := res.prefixes[key]; ok {
				return nil, fmt.Errorf("rad: duplicate sigil (%v, %q)", rule.kind, rule.sigil)
			}
			res.prefixes[key] = rule.expand
			res.prefixKinds[rule.kind] = true
		default:
			return nil, fmt.Errorf("unknown rule op: %v", rule)
		}
	}
	return res, nil
}
