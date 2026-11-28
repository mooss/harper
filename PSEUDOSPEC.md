# Harp pseudo-spec

This document describes the syntax and behavior currently visible in Harp.
This is meant to be a quick and convenient reference, not a complete specification.

## Reading rules

| Token        | Syntax                                                                          |
|--------------|---------------------------------------------------------------------------------|
| *Whitespace* | Spaces, tabs, carriage returns, and newlines                                    |
| *Comments*   | `;` through the end of the line                                                 |
| *Lists*      | `(` followed by zero or more forms, then `)`                                    |
| *Integers*   | Signed decimal integers that fit the machine’s `int` size                       |
| *Strings*    | Double-quoted; `\"`, `\\`, and `\n` are the only supported escapes              |
| *Booleans*   | `true` and `false`                                                              |
| *Keywords*   | Tokens beginning with `:`                                                       |
| *Symbols*    | Start with one of `!&*+-./<=>?@A-Z_a-z`, the rest may include digits and colons |
| *Sigils*     | `'form` is read as `(quote form)`                                               |

## Values

The reader produces *integers*, *strings*, *booleans*, *keywords*, *symbols*, and *lists*.
At runtime, Harp also has builtins, special forms, and closures.

## Evaluation

- *Integers*, *strings*, *booleans*, and *keywords* evaluate to themselves.
- A *symbol* evaluates by looking up its name in the environment.
- *Builtins*, *special forms* and *closures* cannot be evaluated (they are only called when at the head of a list).
- *empty lists* also cannot be evaluated.
- *Non empty lists* evaluate the head and depending on its type:
  - A *builtin* or *closure* receives arguments evaluated from left to right.
  - A *special form* receives the environment and the unevaluated arguments.
  - Any other operator value is an error.

## Default environment

### Integer arithmetic

The operators fold left to right.

| Form          | Behavior                                 |
|---------------|------------------------------------------|
| `(+ n ...)`   | Addition, at least one argument          |
| `(* n ...)`   | Multiplication, at least one argument    |
| `(- n)`       | Negation                                 |
| `(- n m ...)` | Subtraction, at least two arguments      |
| `(/ n m ...)` | Integer division, at least two arguments |

Division by zero panics.

### Integer comparisons

| Form       | Behavior                            |
|------------|-------------------------------------|
| `(< a b)`  | `a` is less than `b`                |
| `(<= a b)` | `a` is less than or equal to `b`    |
| `(> a b)`  | `a` is greater than `b`             |
| `(>= a b)` | `a` is greater than or equal to `b` |

### Special forms

| Form                            | Behavior                                       |
|---------------------------------|------------------------------------------------|
| `(def symbol value)`            | Define a global variable                       |
| `(let (symbol value ...) body)` | Evaluate `body` with temporary local variables |

## Known gaps

- `quote` sigil wraps forms but `quote` is NYI so quoting is currently unusable.
- No conditionals, lambdas/closures construction, recursion, or TCO.
- Division by zero panics rather than erroring.
- No floats, nil, collections beyond lists.
