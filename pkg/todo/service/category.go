package service

import (
	"context"
	"todo/pkg/todo/restapi"

	"github.com/jmoiron/sqlx"
)

func (svc *toDoService) CreateCategory(ctx context.Context, in *restapi.CreateCategoryIn) (*restapi.CategoryOut, error) {

	c := restToDbaCreateCategoryIn(in)

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		id, err := svc.da.InsertCategory(tx, c)
		if err != nil {
			return err
		}

		c, err = svc.da.GetCategoryById(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	out := dbaToRestCategoryOut(c)

	return out, nil
}

func (svc *toDoService) GetCategory(ctx context.Context, id int) (*restapi.CategoryOut, error) {

	c, err := svc.da.GetCategoryById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestCategoryOut(c)

	return out, nil
}
