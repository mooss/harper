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

func writeRepr(sb *strings.Builder, ast any) error {
	switch typed := ast.(type) {
	case typ.Bool:
		sb.WriteString(strconv.FormatBool(bool(typed)))
	case typ.Int:
		sb.WriteString(strconv.FormatInt(int64(typed), 10))
	case typ.Keyword:
		sb.WriteByte(':')
		sb.WriteString(string(typed))
	case typ.List:
		return writeSliceRepr(sb, writeRepr, typed, '(', ' ', ')')
	case typ.String:
		sb.WriteString(lib.Quote(string(typed)))
	case typ.Symbol:
		sb.WriteString(string(typed))
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
