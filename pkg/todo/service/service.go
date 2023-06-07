package service

import (
	"todo/pkg/todo/dba"
)

type toDoService struct {
	da *dba.DatabaseAccess
}

func NewToDoService(da *dba.DatabaseAccess) *toDoService {
	return &toDoService{
		da: da,
	}
}
