package memory

import (
	"testing"
	"time"
	"todo/internal/domain/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	m := New()

	require.NotNil(t, m)
	require.NotNil(t, m.tasks)
	require.NotNil(t, m.lists)
}

func TestAddTask(t *testing.T) {
	fixed := time.Date(2025, 10, 5, 0, 0, 0, 0, time.UTC)

	t.Run("single_success", func(t *testing.T) {
		resetForTest()
		mem := New()

		in := task.Task{ID: 1, Description: "one", Priority: task.PriorityLow, Created: fixed}
		require.NoError(t, mem.AddTask(&in))

		got, err := mem.GetTask(in.ID)
		require.NoError(t, err)
		assert.Equal(t, in, *got)

		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("multiple_success", func(t *testing.T) {
		resetForTest()
		mem := New()

		inputs := []task.Task{
			{ID: 1, Description: "one", Priority: task.PriorityLow, Created: fixed},
			{ID: 2, Description: "two", Priority: task.PriorityMedium, Created: fixed},
			{ID: 3, Description: "three", Priority: task.PriorityHigh, Created: fixed},
		}

		for _, v := range inputs {
			tsk := v
			require.NoError(t, mem.AddTask(&tsk))
		}

		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, len(inputs))

		for _, v := range inputs {
			got, err := mem.GetTask(v.ID)
			require.NoError(t, err)
			assert.Equal(t, v, *got)
		}
	})

	t.Run("duplicate_id_error", func(t *testing.T) {
		resetForTest()
		mem := New()

		first := task.Task{ID: 1, Description: "orig", Priority: task.PriorityLow, Created: fixed}
		dup := task.Task{ID: 1, Description: "dup", Priority: task.PriorityHigh, Created: fixed}

		require.NoError(t, mem.AddTask(&first)) // first add OK

		err := mem.AddTask(&dup) // second add must fail
		require.Error(t, err)

		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, 1)

		got, err := mem.GetTask(first.ID)
		require.NoError(t, err)
		assert.Equal(t, first, *got)
	})
}
