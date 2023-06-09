package dba

import (
	"database/sql"
	"testing"
	"todo/pkg/todo/errors"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInsertCategory_WithoutTransaction(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
	id, err := da.InsertCategory(nil, &c)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, id)
}

func TestInsertCategory_Transaction(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	var idToCheck int
	err := da.ExecuteInTransaction(func(tx *sqlx.Tx) error {
		c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
		id, err := da.InsertCategory(tx, &c)
		idToCheck = id

		assert.Nil(t, err)
		assert.NotEqual(t, 0, id)

		_, err = da.InsertCategory(tx, &c)
		assert.NotNil(t, err)

		return err
	})

	assert.NotNil(t, err)

	var cs []Category
	da.db.Select(&cs, "SELECT * FROM category WHERE id=?", idToCheck)
	assert.Empty(t, cs)
}

func TestGetCategoryById(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
	id, _ := da.InsertCategory(nil, &c)

	cat, err := da.GetCategoryById(nil, id)
	assert.Nil(t, err)
	assert.NotNil(t, cat)

	assert.Equal(t, c.Name, cat.Name)
	assert.Equal(t, c.Description, cat.Description)
}

func TestSearchCategory(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare categories
	c1 := Category{Name: "meditation1", Description: sql.NullString{Valid: true, String: "Category for meditation tasks."}}
	da.InsertCategory(nil, &c1)

	c2 := Category{Name: "meditation2", Description: sql.NullString{Valid: true, String: "Category for meditation tasks."}}
	da.InsertCategory(nil, &c2)

	n := "medita%"
	cs, err := da.SearchCategory(nil, &n)

	assert.Nil(t, err)
	assert.NotNil(t, cs)
	assert.NotEmpty(t, cs)
	assert.Len(t, cs, 2)
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
	id, _ := da.InsertCategory(nil, &c)

	err := da.DeleteCategoryById(nil, id)

	assert.Nil(t, err)

	c1, err := da.GetCategoryById(nil, id)
	assert.NotNil(t, err)
	assert.Nil(t, c1)

	toDoErr, ok := err.(errors.ToDoError)
	assert.True(t, ok)
	assert.NotNil(t, toDoErr)
	assert.Equal(t, ErrEntityNotFound.ErrorCode, toDoErr.ErrorCode)
}

func TestUpdateCategory(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	// prepare category
	c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
	id, _ := da.InsertCategory(nil, &c)
	c.Id = id

	c.Name = "new name"
	c.Description.String = "new description"

	err := da.UpdateCategory(nil, &c)
	assert.Nil(t, err)

	c1, err := da.GetCategoryById(nil, id)
	assert.Nil(t, err)
	assert.NotNil(t, c1)
	assert.Equal(t, c.Name, c1.Name)
	assert.Equal(t, c.Description, c1.Description)
}
