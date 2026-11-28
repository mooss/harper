package harp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mooss/harper/go/harp/typ"
	"github.com/mooss/harper/go/lib"
)

// Repr serializes a Harp value as its representation.
func Repr(ast any) (string, error) {
	var sb strings.Builder
	if err := writeRepr(&sb, ast); err != nil {
		return "", err
	}
	return sb.String(), nil
}

//revive:disable:cyclomatic
func writeRepr(sb *strings.Builder, ast any) error {
	switch node := ast.(type) {
	case typ.Bool:
		sb.WriteString(strconv.FormatBool(bool(node)))
	case typ.Builtin:
		sb.WriteString("<builtin function>")
	case typ.Closure:
		sb.WriteString("<function>")
	case typ.Int:
		sb.WriteString(strconv.FormatInt(int64(node), 10))
	case typ.Keyword:
		sb.WriteByte(':')
		sb.WriteString(string(node))
	case typ.List:
		return writeSliceRepr(sb, writeRepr, node, '(', ' ', ')')
	case typ.String:
		sb.WriteString(lib.Quote(string(node)))
	case typ.Special:
		sb.WriteString("<internal special form>")
	case typ.Symbol:
		sb.WriteString(string(node))
	default:
		return fmt.Errorf("repr: unsupported type %#T", ast)
	}
	return nil
}

func writeSliceRepr(
	sb *strings.Builder, fun func(*strings.Builder, any) error, source []any, left, sep, right byte,
) error {
	sb.WriteByte(left)
	for i, el := range source {
		if i > 0 {
			sb.WriteByte(sep)
		}
		if err := fun(sb, el); err != nil {
			return err
		}
	}
	sb.WriteByte(right)
	return nil
}
