package http

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/pkg/todo"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(svc todo.ToDoService) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/v1/category/{id}", getCategory(svc)).Methods("GET")

	return r
}

func getCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}

		out, err := svc.GetCategory(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}

		fmt.Fprintf(w, "%v", out)
	}
}
