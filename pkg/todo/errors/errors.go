package errors

type ErrorCode int

const (
	ERROR_CODE_ENTITY_NOT_FOUND      ErrorCode = 100
	ERROR_CODE_ENTITY_ALREADY_EXISTS ErrorCode = 101
	ERROR_CODE_BAD_REQUEST           ErrorCode = 102
	ERROR_CODE_DATABASE_ERROR        ErrorCode = 103
	ERROR_CODE_CONFIGURATION_ERROR   ErrorCode = 104
)

type ToDoError struct {
	ErrorCode     ErrorCode
	Text          string
	HttpStatus    int
	ContextValues map[string]interface{}
	Cause         error
}

func WithContextValue(e ToDoError, key string, value interface{}) ToDoError {

	if len(e.ContextValues) == 0 {
		e.ContextValues = map[string]interface{}{}
	}

	e.ContextValues[key] = value
	return e
}

func WithCause(e ToDoError, cause error) ToDoError {

	e.Cause = cause
	return e
}

func (e ToDoError) Error() string {
	return "error"
}

func (e ToDoError) Unwrap() error {
	return e.Cause
}
