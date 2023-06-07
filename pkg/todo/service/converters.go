package service

import (
	"database/sql"
	"todo/pkg/todo/dba"
	"todo/pkg/todo/restapi"
)

func dbaToRestCategoryOut(in *dba.Category) *restapi.CategoryOut {
	if in == nil {
		return nil
	}

	out := restapi.CategoryOut{
		Id:          int32(in.Id),
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
