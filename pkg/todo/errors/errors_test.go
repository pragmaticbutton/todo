package errors

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var err error

	err = &ToDoError{ErrorCode: ERROR_CODE_BAD_REQUEST}
	s := err.Error()

	fmt.Println(s)

}
