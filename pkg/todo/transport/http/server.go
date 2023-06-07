package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo/pkg/todo"
	"todo/pkg/todo/restapi"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(svc todo.ToDoService) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/v1/category", createCategory(svc)).Methods("POST")
	r.HandleFunc("/v1/category/{id}", getCategory(svc)).Methods("GET")

	return r
}

func createCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var in restapi.CreateCategoryIn
		d := json.NewDecoder(r.Body)
		d.Decode(&in)

		out, err := svc.CreateCategory(r.Context(), &in)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}

		encodeOutput(w, out)
	}
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

		encodeOutput(w, out)
	}
}

func encodeOutput(w http.ResponseWriter, out interface{}) {
	e := json.NewEncoder(w)
	err := e.Encode(out)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
}
