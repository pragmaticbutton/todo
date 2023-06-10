package todo

import (
	"context"
	"todo/pkg/todo/restapi"
)

type ToDoService interface {

	// category services
	CreateCategory(ctx context.Context, in *restapi.CreateCategoryIn) (*restapi.CategoryOut, error)
	GetCategory(ctx context.Context, id int) (*restapi.CategoryOut, error)
	SearchCategory(ctx context.Context, params *restapi.SearchCategoryParams) (*restapi.SearchCategoryOut, error)
	DeleteCategory(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, id int, in *restapi.UpdateCategoryIn) (*restapi.CategoryOut, error)

	// task services
	CreateTask(ctx context.Context, in *restapi.CreateTaskIn) (*restapi.TaskOut, error)
}
