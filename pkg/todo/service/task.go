package service

import (
	"context"
	"todo/pkg/todo/dba"
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

func (svc *toDoService) GetTask(ctx context.Context, id int32) (*restapi.TaskOut, error) {

	t, err := svc.da.GetTaskById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestTaskOut(t)

	return out, nil
}

func (svc *toDoService) DeleteTask(ctx context.Context, id int32) error {

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

func (svc *toDoService) SearchTask(ctx context.Context, params *restapi.SearchTaskParams) (*restapi.SearchTaskOut, error) {

	var name *string
	if params.Name != nil {
		n := replaceWildCards(*params.Name)
		name = &n
	}

	var done *int8
	if params.Done != nil {
		d := int8(0)
		if *params.Done {
			d = int8(1)
		}
		done = &d

	}
	ts, err := svc.da.SearchTask(nil, name, params.CategoryId, taskPriorityForSearch(params.Priority), done)
	if err != nil {
		return nil, err
	}

	tsOut := make([]restapi.TaskOut, len(ts))
	for i, t := range ts {
		tsOut[i] = *dbaToRestTaskOut(&t)
	}
	out := restapi.SearchTaskOut{Tasks: &tsOut}

	return &out, nil
}

func (svc *toDoService) FinishTask(ctx context.Context, id int32) (*restapi.TaskOut, error) {

	var t *dba.Task
	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {
		var err error
		t, err = svc.da.GetTaskById(tx, id)
		if err != nil {
			return err
		}

		t.Done = 1 // task is finished
		err = svc.da.UpdateTask(tx, t)
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

func (svc *toDoService) UpdateTask(ctx context.Context, id int32, in *restapi.UpdateTaskIn) (*restapi.TaskOut, error) {

	var t *dba.Task
	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		var err error
		t, err = svc.da.GetTaskById(tx, id)
		if err != nil {
			return err
		}

		updateTaskWithValuesFromRequest(t, in)
		err = svc.da.UpdateTask(tx, t)
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

func taskPriorityForSearch(in *restapi.TaskPriority) *dba.TaskPriorityType {
	if in == nil {
		return nil
	}

	out := restToDbaTaskPriority(*in)
	return &out
}

func updateTaskWithValuesFromRequest(t *dba.Task, in *restapi.UpdateTaskIn) {
	if in.Name != nil {
		t.Name = *in.Name
	}

	if in.CategoryId != nil {
		t.FkCategory = *in.CategoryId
	}

	if in.Description != nil {
		t.Description = stringPToNullString(in.Description)
	}

	if in.Done != nil {
		if *in.Done {
			t.Done = int8(1)
		} else {
			t.Done = int8(0)
		}
	}

	if in.Priority != nil {
		t.Priority = restToDbaTaskPriority(*in.Priority)
	}
}
