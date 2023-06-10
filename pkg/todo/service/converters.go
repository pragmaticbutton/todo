package service

import (
	"database/sql"
	"strings"
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

func replaceWildCards(s string) string {
	return strings.ReplaceAll(s, "*", "%")
}

func restToDbaCreateTaskIn(in *restapi.CreateTaskIn) *dba.Task {
	if in == nil {
		return nil
	}

	out := dba.Task{
		Name:        in.Name,
		FkCategory:  int(in.CategoryId),
		Priority:    restToDbaTaskPriority(in.Priority),
		Description: stringPToNullString(in.Description),
	}

	return &out
}

func restToDbaTaskPriority(in restapi.TaskPriority) dba.TaskPriorityType {
	switch in {
	case restapi.HIGH:
		return dba.TASK_PRIORITY_HIGH
	case restapi.MEDIUM:
		return dba.TASK_PRIORITY_MEDIUM
	case restapi.LOW:
		return dba.TASK_PRIORITY_LOW
	}

	panic("unsupported task priority: " + in)
}

func dbaToRestTaskOut(in *dba.Task) *restapi.TaskOut {
	if in == nil {
		return nil
	}

	out := restapi.TaskOut{
		Id:          int32(in.Id),
		Name:        in.Name,
		CategoryId:  int32(in.FkCategory),
		Priority:    dbaToRestTaskPriority(in.Priority),
		Done:        in.Done != 0,
		Description: nullStringToStringP(in.Description),
		Created:     in.Created,
		LastChanged: in.LastChanged,
	}

	return &out
}

func dbaToRestTaskPriority(in dba.TaskPriorityType) restapi.TaskPriority {
	switch in {
	case dba.TASK_PRIORITY_HIGH:
		return restapi.HIGH
	case dba.TASK_PRIORITY_MEDIUM:
		return restapi.MEDIUM
	case dba.TASK_PRIORITY_LOW:
		return restapi.LOW
	}

	panic("unsupported task priority: " + in)
}
