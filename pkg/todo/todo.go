package todo

import (
	"context"
	"todo/pkg/todo/restapi"
)

type ToDoService interface {

	// category services
	CreateCategory(ctx context.Context, in *restapi.CreateCategoryIn) (*restapi.CategoryOut, error)
	GetCategory(ctx context.Context, id int32) (*restapi.CategoryOut, error)
	SearchCategory(ctx context.Context, params *restapi.SearchCategoryParams) (*restapi.SearchCategoryOut, error)
	DeleteCategory(ctx context.Context, id int32) error
	UpdateCategory(ctx context.Context, id int32, in *restapi.UpdateCategoryIn) (*restapi.CategoryOut, error)

	// task services
	CreateTask(ctx context.Context, in *restapi.CreateTaskIn) (*restapi.TaskOut, error)
	GetTask(ctx context.Context, id int32) (*restapi.TaskOut, error)
	DeleteTask(ctx context.Context, id int32) error
	SearchTask(ctx context.Context, params *restapi.SearchTaskParams) (*restapi.SearchTaskOut, error)
	FinishTask(ctx context.Context, id int32) (*restapi.TaskOut, error)
	UpdateTask(ctx context.Context, id int32, in *restapi.UpdateTaskIn) (*restapi.TaskOut, error)

	// user services
	CreateUser(ctx context.Context, in *restapi.CreateUserIn) (*restapi.UserOut, error)
}
