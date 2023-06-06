package dba

import (
	"database/sql"
	"testing"

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
