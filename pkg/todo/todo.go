package todo

import (
	"context"
	"todo/pkg/todo/restapi"
)

type ToDoService interface {
	GetCategory(ctx context.Context, id int) (*restapi.CategoryOut, error)
}
