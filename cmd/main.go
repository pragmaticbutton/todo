package main

import (
	"fmt"
	stdhttp "net/http"
	"todo/pkg/todo/config"
	"todo/pkg/todo/dba"
	"todo/pkg/todo/errors"
	"todo/pkg/todo/middleware"
	"todo/pkg/todo/service"
	"todo/pkg/todo/transport/http"
)

var (
	errBadRequest = errors.ToDoError{ErrorCode: errors.ERROR_CODE_BAD_REQUEST, Text: "Bad request", HttpStatus: stdhttp.StatusBadRequest}
	errUnknown    = errors.ToDoError{ErrorCode: errors.ERROR_CODE_UNKNOWN_ERROR, Text: "Unknown error", HttpStatus: stdhttp.StatusInternalServerError}
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		panic(fmt.Sprintf("reading config failed: %v", err))
	}

	da, err := dba.NewDatabaseAccess(config.Dsn)
	if err != nil {
		panic(err)
	}

	svc := service.NewToDoService(da)

	httpHandler := http.NewHTTPHandler(
		svc,
		http.ErrorEncoder(errUnknown),
		middleware.OpenapiRequestValidatorMiddleware{ErrorEncode: http.ErrorEncoder(errBadRequest)}.Middleware,
	)

	stdhttp.ListenAndServe("localhost:8090", httpHandler)
}
