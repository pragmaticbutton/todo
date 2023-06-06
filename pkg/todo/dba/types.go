package dba

import (
	"database/sql"
	"time"
)

type TaskPriorityType string

const (
	TASK_PRIORITY_HIGH   TaskPriorityType = "HIGH"
	TASK_PRIORITY_MEDIUM TaskPriorityType = "MEDIUM"
	TASK_PRIORITY_LOW    TaskPriorityType = "LOW"
)

// Task struct represents task table.
type Task struct {
	Id          int
	Name        string
	FkCateogory int
	Priority    TaskPriorityType
	Done        int8
	Created     time.Time
	LastChanged time.Time
}

// Category struct represents category table.
type Category struct {
	Id          int
	Name        string
	Description sql.NullString
	Created     time.Time
	LastChanged time.Time
}
