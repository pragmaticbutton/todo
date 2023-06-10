package middleware

import (
	"context"
	"net/http"
	"todo/pkg/todo/errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

var (
	ErrBadRequest = errors.ToDoError{ErrorCode: errors.ERROR_CODE_BAD_REQUEST, Text: "Bad request", HttpStatus: http.StatusBadRequest}
)

func OpenapiRequestValidatorMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		doc, err := openapi3.NewLoader().LoadFromFile("../restapi/todo.yml")
		if err != nil {
			panic(err)
		}

		router, err := gorillamux.NewRouter(doc)
		if err != nil {
			panic(err)
		}

		route, pathParams, err := router.FindRoute(r)
		if err != nil {
			panic(err)
		}

		rvi := &openapi3filter.RequestValidationInput{
			Request:     r,
			Route:       route,
			PathParams:  pathParams,
			QueryParams: r.Form,
		}
		err = openapi3filter.ValidateRequest(context.Background(), rvi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//TODO: error
			return
		}

		h.ServeHTTP(w, r)
	})
}
