package service

import (
	"context"
	"database/sql"
	"todo/pkg/todo/dba"
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

func (svc *toDoService) GetCategory(ctx context.Context, id int32) (*restapi.CategoryOut, error) {

	c, err := svc.da.GetCategoryById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestCategoryOut(c)

	return out, nil
}

func (svc *toDoService) SearchCategory(ctx context.Context, params *restapi.SearchCategoryParams) (*restapi.SearchCategoryOut, error) {

	var name *string
	if params.Name != nil {
		n := replaceWildCards(*params.Name)
		name = &n
	}

	cs, err := svc.da.SearchCategory(nil, name, paginationForCategory(params))
	if err != nil {
		return nil, err
	}

	csOut := make([]restapi.CategoryOut, len(cs))
	for i, c := range cs {
		csOut[i] = *dbaToRestCategoryOut(&c)
	}
	count, err := svc.da.CountCategory(nil, name)
	if err != nil {
		return nil, err
	}

	out := &restapi.SearchCategoryOut{Categories: &csOut, StartIndex: params.StartIndex, TotalRecords: count}

	return out, nil
}

func (svc *toDoService) DeleteCategory(ctx context.Context, id int32) error {

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		_, err := svc.da.GetCategoryById(tx, id)
		if err != nil {
			return err
		}

		err = svc.da.DeleteCategoryById(tx, id)
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

func (svc *toDoService) UpdateCategory(ctx context.Context, id int32, in *restapi.UpdateCategoryIn) (*restapi.CategoryOut, error) {

	var c *dba.Category
	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		var err error
		c, err = svc.da.GetCategoryById(tx, id)
		if err != nil {
			return err
		}

		updateCategoryWithValuesFromRequest(c, in)
		err = svc.da.UpdateCategory(tx, c)
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

func updateCategoryWithValuesFromRequest(c *dba.Category, in *restapi.UpdateCategoryIn) {
	if in.Name != nil {
		c.Name = *in.Name
	}

	if in.Description != nil {
		c.Description = sql.NullString{String: *in.Description, Valid: true}
	}
}

func paginationForCategory(params *restapi.SearchCategoryParams) *dba.Pagination {
	opts := []dba.PaginationOption{}

	if params.OrderBy != nil {
		ob := restToDbaCategoryOrderBy(*params.OrderBy)
		opts = append(opts, dba.WithOrderBy(&ob))
	}
	if params.OrderDirection != nil {
		od := restToDbaOrderDirection(*params.OrderDirection)
		opts = append(opts, dba.WithOrderDirection(&od))
	}
	if params.StartIndex != nil {
		opts = append(opts, dba.WithStartIndex(params.StartIndex))
	}
	if params.RecordsPerPage != nil {
		opts = append(opts, dba.WithRecordsPerPage(params.RecordsPerPage))
	}

	return dba.NewPagination(opts...)
}
