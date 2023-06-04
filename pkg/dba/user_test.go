package dba

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	// Given
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	fkRole := prepareRole()

	userToInsert := User{Username: "someusername", FkRole: fkRole}
	id, err := da.InsertUser(nil, &userToInsert)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, id)
}
