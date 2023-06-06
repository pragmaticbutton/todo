package dba

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertCategory(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	c := Category{Name: "health", Description: sql.NullString{Valid: true, String: "Category for health tasks."}}
	id, err := da.InsertCategory(nil, &c)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, id)
}
