package service

import (
	"database/sql"
	"todo/pkg/todo/dba"
	"todo/pkg/todo/restapi"
)

func restToDbaCreateCategoryIn(in *restapi.CreateCategoryIn) *dba.Category {
	if in == nil {
		return nil
	}

	out := dba.Category{
		Name:        in.Name,
		Description: stringPToNullString(in.Description),
	}

	return &out
}

func dbaToRestCategoryOut(in *dba.Category) *restapi.CategoryOut {
	if in == nil {
		return nil
	}

	out := restapi.CategoryOut{
		Id:          int32(in.Id),
		Name:        in.Name,
		Description: nullStringToStringP(in.Description),
		Created:     in.Created,
		LastChanged: in.LastChanged,
	}

	return &out
}

func nullStringToStringP(in sql.NullString) *string {
	if !in.Valid {
		return nil
	}

	return &in.String
}

func stringPToNullString(in *string) sql.NullString {
	if in == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: *in, Valid: true}
}
