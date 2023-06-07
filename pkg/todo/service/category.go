package service

import (
	"context"
	"fmt"
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
		c.Id = id

		return nil
	})
	if err != nil {
		return nil, err
	}

	out := dbaToRestCategoryOut(c)

	return out, nil
}

func (svc *toDoService) GetCategory(ctx context.Context, id int) (*restapi.CategoryOut, error) {

	if err := validateGetCategoryRequest(id); err != nil {
		return nil, err
	}

	c, err := svc.da.GetCategoryById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestCategoryOut(c)

	return out, nil
}

func validateGetCategoryRequest(id int) error {

	if id < 0 {
		return fmt.Errorf("id must be a positive number")
	}

	return nil
}
