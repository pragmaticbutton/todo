package http

import (
	"fmt"
	"net/http"
	"todo/pkg/todo"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(svc *todo.ToDoService) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/v1/category", getCategory(svc)).Methods("GET")

	return r
}

func getCategory(svc *todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", (*svc).GetCategory("string"))
	}
}
