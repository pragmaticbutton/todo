package todo

import (
	"context"
	"todo/pkg/todo/restapi"
)

type ToDoService interface {
	CreateCategory(ctx context.Context, in *restapi.CreateCategoryIn) (*restapi.CategoryOut, error)
	GetCategory(ctx context.Context, id int) (*restapi.CategoryOut, error)
}
