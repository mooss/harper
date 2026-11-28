package main

import (
	"fmt"
	"os"

	"github.com/mooss/bagend/go/damn"
	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/harper/go/harp"
	"github.com/mooss/harper/go/harp/typ"
)

var expected = "expected input from stdin xor a single filename"

func main() {
	args := os.Args[1:]
	program := damn.Must(sysu.ReadStdin())
	if len(program) == 0 {
		if len(args) != 1 {
			damn.Fatal(expected)
		}
		program = string(damn.Must(os.ReadFile(args[0]))) //nolint:gosec
	} else if len(args) != 0 {
		damn.Fatal(expected)
	}

	fmt.Println(damn.Must(rep(program, harp.MkEnv())))
}

// rep calls read, eval and print on its input.
func rep(input string, envr *typ.Env) (string, error) {
	prog, err := harp.ParseAll(input)
	if err != nil {
		return "", err
	}

	if len(prog) == 0 {
		return "", nil
	}

	for _, ast := range prog[:len(prog)-1] {
		_, err := harp.Eval(ast, envr)
		if err != nil {
			return "", err
		}
	}

	res, err := harp.Eval(prog[len(prog)-1], envr)
	if err != nil {
		return "", err
	}
	return harp.Repr(res)
}
