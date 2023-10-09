package errors

import "fmt"

type ErrorCode int

const (
	ERROR_CODE_ENTITY_NOT_FOUND ErrorCode = iota + 100
	ERROR_CODE_ENTITY_ALREADY_EXISTS
	ERROR_CODE_BAD_REQUEST
	ERROR_CODE_DATABASE_ERROR
	ERROR_CODE_CONFIGURATION_ERROR
	ERROR_CODE_UNAUTHORIZED
	ERROR_CODE_UNKNOWN_ERROR
)

type ToDoError struct {
	ErrorCode     ErrorCode
	Text          string
	HttpStatus    int
	ContextValues map[string]string
	Cause         error
}

func WithContextValue(e ToDoError, key string, value string) ToDoError {

	if len(e.ContextValues) == 0 {
		e.ContextValues = map[string]string{}
	}

	e.ContextValues[key] = value
	return e
}

func WithCause(e ToDoError, cause error) ToDoError {

	e.Cause = cause
	return e
}

func (e ToDoError) Error() string {
	//TODO
	return fmt.Sprintf("Error code: %d, error text: %s, http status: %d, context values: %v, cause: %v",
		e.ErrorCode, e.Text, e.HttpStatus, e.ContextValues, e.Cause)
}

func (e ToDoError) Unwrap() error {
	return e.Cause
}
