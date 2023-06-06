package service

import (
	"todo/pkg/todo"
	"todo/pkg/todo/dba"
)

type toDoService struct {
	da *dba.DatabaseAccess
}

func NewToDoService(da *dba.DatabaseAccess) todo.ToDoService {
	return &toDoService{
		da: da,
	}
}
