package main

import (
	"fmt"
	"net/http"
	"todo/pkg/config"
	"todo/pkg/dba"
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

	fmt.Println(da)

	http.ListenAndServe("localhost:8090", nil)
}
