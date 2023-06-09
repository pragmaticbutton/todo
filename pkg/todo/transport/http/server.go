package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/pkg/todo"
	"todo/pkg/todo/errors"
	"todo/pkg/todo/restapi"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(svc todo.ToDoService) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/v1/category", createCategory(svc)).Methods("POST")
	r.HandleFunc("/v1/category/{id}", getCategory(svc)).Methods("GET")
	r.HandleFunc("/v1/category", searchCategory(svc)).Methods("GET")
	r.HandleFunc("/v1/category/{id}", deleteCategory(svc)).Methods("DELETE")
	r.HandleFunc("/v1/category/{id}", updateCategory(svc)).Methods("PATCH")

	return r
}

func createCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var in restapi.CreateCategoryIn
		d := json.NewDecoder(r.Body)
		d.Decode(&in)

		out, err := svc.CreateCategory(r.Context(), &in)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func getCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			encodeError(w, err)
			return
		}

		out, err := svc.GetCategory(r.Context(), id)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func searchCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := restapi.SearchCategoryParams{}

		name := r.FormValue("name")
		if name != "" {
			params.Name = &name
		}

		out, err := svc.SearchCategory(r.Context(), &params)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func deleteCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			encodeError(w, err)
			return
		}

		err = svc.DeleteCategory(r.Context(), id)
		if err != nil {
			encodeError(w, err)
			return
		}
	}
}

func updateCategory(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			encodeError(w, err)
			return
		}

		var in restapi.UpdateCategoryIn
		d := json.NewDecoder(r.Body)
		d.Decode(&in)

		out, err := svc.UpdateCategory(r.Context(), id, &in)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func encodeError(w http.ResponseWriter, e error) {
	if e == nil {
		return
	}
	enc := json.NewEncoder(w)
	toDoErr, ok := e.(errors.ToDoError)
	if !ok {
		enc.Encode(e)
		return
	}
	enc.Encode(toDoErr)
}

func encodeOutput(w http.ResponseWriter, out interface{}) {
	if out == nil {
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(out)
	if err != nil {
		encodeError(w, err)
	}
}
