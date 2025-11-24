package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	l := New(1, "shopping", "grocery shopping list")

	assert.NotNil(t, l)
	assert.Equal(t, uint32(1), l.ID)
	assert.Equal(t, "shopping", l.Name)
	assert.Equal(t, "grocery shopping list", l.Description)
	assert.False(t, l.Created.IsZero())
	assert.Nil(t, l.Updated)
}
