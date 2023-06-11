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

func (svc *toDoService) GetTask(ctx context.Context, id int) (*restapi.TaskOut, error) {

	t, err := svc.da.GetTaskById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestTaskOut(t)

	return out, nil
}

func (svc *toDoService) DeleteTask(ctx context.Context, id int) error {

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		_, err := svc.da.GetTaskById(tx, id)
		if err != nil {
			return err
		}

		err = svc.da.DeleteTaskById(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
