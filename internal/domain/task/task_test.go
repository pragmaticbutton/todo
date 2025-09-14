package task

import (
	"testing"
	"todo/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	task := New(1, "cookies", PriorityMedium, utils.Ptr(uint32(1)))

	assert.NotNil(t, task)
	assert.Equal(t, uint32(1), task.ID)
	assert.Equal(t, "cookies", task.Description)
	assert.Equal(t, PriorityMedium, task.Priority)
	assert.Equal(t, utils.Ptr(uint32(1)), task.ListID)
	assert.False(t, task.Created.IsZero())
	assert.True(t, task.Updated.IsZero())
	assert.False(t, task.Done)

}
