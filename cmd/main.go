package main

import (
	"fmt"
	stdhttp "net/http"
	"todo/pkg/todo/config"
	"todo/pkg/todo/dba"
	"todo/pkg/todo/middleware"
	"todo/pkg/todo/service"
	"todo/pkg/todo/transport/http"
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
		middleware.OpenapiRequestValidatorMiddleware,
	)

	stdhttp.ListenAndServe("localhost:8090", httpHandler)
}
