package dba

import (
	"testing"
	"todo/pkg/todo/errors"

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

func TestGetUser(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	user := User{
		Username: "jelena",
		Password: "somepassword",
	}

	id, err := da.InsertUser(nil, &user)

	assert.Nil(t, err)
	assert.NotZero(t, id)

	u, err := da.GetUserById(nil, id)
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, id, u.Id)
	assert.Equal(t, user.Username, u.Username)
}

func TestDeleteUser(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	user := User{
		Username: "jelena",
		Password: "somepassword",
	}

	id, err := da.InsertUser(nil, &user)

	assert.Nil(t, err)
	assert.NotZero(t, id)

	err = da.DeleteUserById(nil, id)
	assert.Nil(t, err)

	_, err = da.GetUserById(nil, id)
	assert.NotNil(t, err)
	todoErr, ok := err.(errors.ToDoError)
	assert.True(t, ok)
	assert.NotNil(t, todoErr)
	assert.Equal(t, errors.ERROR_CODE_ENTITY_NOT_FOUND, todoErr.ErrorCode)
}

func TestSearchUser(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase()

	user := User{
		Username: "jelena",
		Password: "somepassword",
	}

	id, err := da.InsertUser(nil, &user)

	assert.Nil(t, err)
	assert.NotZero(t, id)

	us, err := da.SearchUser(nil, nil)
	assert.Nil(t, err)
	assert.Len(t, us, 1)
}
