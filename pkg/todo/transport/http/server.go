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

func NewHTTPHandler(svc todo.ToDoService, mws ...mux.MiddlewareFunc) http.Handler {

	r := mux.NewRouter()

	for _, mw := range mws {
		r.Use(mw)
	}

	// category
	r.HandleFunc("/v1/category", createCategory(svc)).Methods("POST")
	r.HandleFunc("/v1/category/{id}", getCategory(svc)).Methods("GET")
	r.HandleFunc("/v1/category", searchCategory(svc)).Methods("GET")
	r.HandleFunc("/v1/category/{id}", deleteCategory(svc)).Methods("DELETE")
	r.HandleFunc("/v1/category/{id}", updateCategory(svc)).Methods("PATCH")

	// task
	r.HandleFunc("/v1/task", createTask(svc)).Methods("POST")
	r.HandleFunc("/v1/task/{id}", getTask(svc)).Methods("GET")
	r.HandleFunc("/v1/task/{id}", deleteTask(svc)).Methods("DELETE")
	r.HandleFunc("/v1/task", searchTask(svc)).Methods("GET")

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

		out, err := svc.GetCategory(r.Context(), int32(id))
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

		err = svc.DeleteCategory(r.Context(), int32(id))
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

		out, err := svc.UpdateCategory(r.Context(), int32(id), &in)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func createTask(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var in restapi.CreateTaskIn
		d := json.NewDecoder(r.Body)
		d.Decode(&in)

		out, err := svc.CreateTask(r.Context(), &in)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func getTask(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			encodeError(w, err)
			return
		}

		out, err := svc.GetTask(r.Context(), int32(id))
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func deleteTask(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			encodeError(w, err)
			return
		}

		err = svc.DeleteTask(r.Context(), int32(id))
		if err != nil {
			encodeError(w, err)
			return
		}
	}
}

func searchTask(svc todo.ToDoService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := restapi.SearchTaskParams{}

		name := r.FormValue("name")
		if name != "" {
			params.Name = &name
		}

		categoryId := r.FormValue("categoryId")
		if categoryId != "" {
			cId, err := strconv.Atoi(categoryId)
			if err != nil {
				encodeError(w, err)
				return
			}
			c := int32(cId)
			params.CategoryId = &c
		}

		priority := r.FormValue("priority")
		if priority != "" {
			p := restapi.TaskPriority(priority)
			params.Priority = &p
		}

		done := r.FormValue("done")
		if done != "" && (done == "false" || done == "true") {
			d := done == "true"
			params.Done = &d
		}

		out, err := svc.SearchTask(r.Context(), &params)
		if err != nil {
			encodeError(w, err)
			return
		}

		encodeOutput(w, out)
	}
}

func encodeOutput(w http.ResponseWriter, out interface{}) {
	if out == nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err := enc.Encode(out)
	if err != nil {
		panic(err)
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
	w.WriteHeader(toDoErr.HttpStatus)
	enc.Encode(toDoErr)
}
