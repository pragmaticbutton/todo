package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {

	e := errors.New("first")
	e1 := errors.New("second")
	e2 := errors.Join(e, e1)

	fmt.Println(errors.Unwrap(e2))
}
