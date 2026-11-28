package lib

import (
	"fmt"
	"strconv"
)

////////////////////
// Error builders //

// ExpectAtLeastNArgs formats an error if got contains less than expected values.
// Otherwise it returns nil.
func ExpectAtLeastNArgs[T any](expected int, got []T) error {
	if len(got) >= expected {
		return nil
	}
	return FmtExpectedNArgs(fmt.Sprintf("at least %d", expected), len(got))
}

// ExpectNArgs formats an error if got does not contain exactly expected values.
// Otherwise it returns nil.
func ExpectNArgs[T any](expected int, got []T) error {
	if expected == len(got) {
		return nil
	}
	return FmtExpectedNArgs(strconv.Itoa(expected), len(got))
}

// Expected formats an error explaining that an unexpected type was given.
// When the unexpected type is already an error, it is returned as-is.
func Expected[T ~string](expected T, got any) error {
	if err, isErr := got.(error); isErr {
		return err
	}
	return &TypeError{Want: string(expected), Got: got}
}

// FmtExpectedNArgs format an error explaining that an unexpected number of arguments was given.
func FmtExpectedNArgs(expected string, got int) error {
	return &ArityError{Want: expected, Got: got}
}

/////////////////
// Error types //

type ArityError struct {
	Want string
	Got  int
}

func (e *ArityError) Error() string {
	return fmt.Sprintf("expected %s args, got %d", e.Want, e.Got)
}

type TypeError struct {
	Want string
	Got  any
}

func (e *TypeError) Error() string {
	if e.Got == nil {
		return fmt.Sprintf("expected %s, got nil (nil)", e.Want)
	}
	return fmt.Sprintf("expected %s, got %T (%v)", e.Want, e.Got, e.Got)
}
