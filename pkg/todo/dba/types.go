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
	Id          int              `db:"id"`
	Name        string           `db:"name"`
	FkCategory  int              `db:"fk_category"`
	Priority    TaskPriorityType `db:"priority"`
	Done        int8             `db:"done"`
	Description sql.NullString   `db:"description"`
	Created     time.Time        `db:"created"`
	LastChanged time.Time        `db:"last_changed"`
}

// Category struct represents category table.
type Category struct {
	Id          int            `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Created     time.Time      `db:"created"`
	LastChanged time.Time      `db:"last_changed"`
}
