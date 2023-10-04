package dba

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	user := User{
		Username: "jelena",
		Password: "somepassword",
	}

	id, err := da.InsertUser(nil, &user)

	assert.Nil(t, err)
	assert.NotEqual(t, int32(0), id)
}
