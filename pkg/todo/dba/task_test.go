package dba

import (
	"database/sql"
	"testing"
	"time"
	"todo/pkg/todo/errors"

	"github.com/stretchr/testify/assert"
)

func TestInsertTask(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "school", Description: sql.NullString{String: "School related tasks", Valid: true}}
	cId, _ := da.InsertCategory(nil, &c)

	task := Task{
		Name:        "geography homework",
		FkCategory:  cId,
		Priority:    TASK_PRIORITY_MEDIUM,
		Done:        0,
		Description: sql.NullString{String: "Finish geography homework", Valid: true},
	}

	tId, err := da.InsertTask(nil, &task)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, tId)
}

func TestGetTaskById(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "school", Description: sql.NullString{String: "School related tasks", Valid: true}}
	cId, _ := da.InsertCategory(nil, &c)

	// prepare task
	task := Task{
		Name:        "geography homework",
		FkCategory:  cId,
		Priority:    TASK_PRIORITY_MEDIUM,
		Done:        0,
		Description: sql.NullString{String: "Finish geography homework", Valid: true},
	}
	tId, _ := da.InsertTask(nil, &task)

	// TODO: check why is this necessary...
	time.Sleep(time.Nanosecond * 10)

	task1, err := da.GetTaskById(nil, tId)

	assert.Nil(t, err)
	assert.NotNil(t, task1)
	assert.Equal(t, task.Name, task1.Name)
	assert.Equal(t, task.Description, task1.Description)
	assert.Equal(t, task.Done, task1.Done)
	assert.Equal(t, task.Priority, task1.Priority)

}

func TestDeleteTaskById(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "school", Description: sql.NullString{String: "School related tasks", Valid: true}}
	cId, _ := da.InsertCategory(nil, &c)

	// prepare task
	task := Task{
		Name:        "geography homework",
		FkCategory:  cId,
		Priority:    TASK_PRIORITY_MEDIUM,
		Done:        0,
		Description: sql.NullString{String: "Finish geography homework", Valid: true},
	}
	tId, _ := da.InsertTask(nil, &task)

	// TODO: check why is this necessary...
	time.Sleep(time.Nanosecond * 5)

	err := da.DeleteCategoryById(nil, tId)

	// TODO: check why is this necessary...
	time.Sleep(time.Nanosecond * 5)

	assert.Nil(t, err)
	ta, err := da.GetTaskById(nil, tId)
	assert.NotNil(t, err)
	assert.Nil(t, ta)
	toDoErr := err.(errors.ToDoError)
	assert.NotNil(t, toDoErr)
	assert.Equal(t, ErrEntityNotFound.ErrorCode, toDoErr.ErrorCode)

}

func TestSearchTask(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "school", Description: sql.NullString{String: "School related tasks", Valid: true}}
	cId, _ := da.InsertCategory(nil, &c)

	// prepare tasks
	task1 := Task{
		Name:        "geography homework",
		FkCategory:  cId,
		Priority:    TASK_PRIORITY_MEDIUM,
		Done:        0,
		Description: sql.NullString{String: "Finish geography homework", Valid: true},
	}
	da.InsertTask(nil, &task1)

	// prepare tasks
	task2 := Task{
		Name:        "history homework",
		FkCategory:  cId,
		Priority:    TASK_PRIORITY_LOW,
		Done:        0,
		Description: sql.NullString{String: "Finish history homework", Valid: true},
	}
	da.InsertTask(nil, &task2)

	time.Sleep(time.Nanosecond * 5)

	n := "%homework"
	p := TASK_PRIORITY_MEDIUM
	ts, err := da.SearchTask(nil, &n, nil, &p, nil)

	assert.Nil(t, err)
	assert.Len(t, ts, 1)
}
