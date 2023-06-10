package service

import (
	"context"
	"todo/pkg/todo/restapi"

	"github.com/jmoiron/sqlx"
)

func (svc *toDoService) CreateTask(ctx context.Context, in *restapi.CreateTaskIn) (*restapi.TaskOut, error) {

	t := restToDbaCreateTaskIn(in)

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		id, err := svc.da.InsertTask(tx, t)
		if err != nil {
			return err
		}

		t, err = svc.da.GetTaskById(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	out := dbaToRestTaskOut(t)

	return out, nil
}
