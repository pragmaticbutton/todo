package http

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"todo/pkg/todo/errors"
	"todo/pkg/todo/restapi"
)

var (
	errUnknown = errors.ToDoError{ErrorCode: errors.ERROR_CODE_UNKNOWN_ERROR, Text: "Unknown error", HttpStatus: http.StatusInternalServerError}
)

type ErrorEncoderFunc func(http.ResponseWriter, error)

func ErrorEncoder(defErr errors.ToDoError) ErrorEncoderFunc {
	return func(w http.ResponseWriter, err error) {
		if err == nil {
			panic("error is nil")
		}

		var restError restapi.Error
		var httpStatus int
		var toDoErr errors.ToDoError
		ok := stderrors.As(err, &toDoErr)
		if !ok {
			restError = createRestError(errors.WithCause(errUnknown, err))
			httpStatus = errUnknown.HttpStatus
		} else {
			restError = createRestError(toDoErr)
			httpStatus = toDoErr.HttpStatus
		}

		w.WriteHeader(httpStatus)
		enc := json.NewEncoder(w)
		if err := enc.Encode(restError); err != nil {
			panic(err)
		}
	}
}

func encodeError1(w http.ResponseWriter, err error) {
	if err == nil {
		panic("error is nil")
	}

	var restError restapi.Error
	var httpStatus int
	var toDoErr errors.ToDoError
	ok := stderrors.As(err, &toDoErr)
	if !ok {
		restError = createRestError(errors.WithCause(errUnknown, err))
		httpStatus = errUnknown.HttpStatus
	} else {
		restError = createRestError(toDoErr)
		httpStatus = toDoErr.HttpStatus
	}

	w.WriteHeader(httpStatus)
	enc := json.NewEncoder(w)
	if err := enc.Encode(restError); err != nil {
		panic(err)
	}

}

func createRestError(e errors.ToDoError) restapi.Error {
	err := restapi.Error{
		ErrorCode: int32(e.ErrorCode),
		ErrorText: e.Text,
	}

	if e.Cause != nil {
		tmp := e.Cause.Error()
		err.Cause = &tmp
	}

	if e.ContextValues != nil {
		cvs := []restapi.ErrorContextValue{}
		for k, v := range e.ContextValues {
			cvs = append(cvs, restapi.ErrorContextValue{ContextLabel: k, ContextValue: v})
		}
		err.ContextValues = &cvs
	}

	return err
}
