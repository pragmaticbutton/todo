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

func TestGetTask(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		resetForTest()
		mem := New()
		task := task.Task{ID: 1, Description: "one", Priority: task.PriorityLow, Created: time.Now()}
		require.NoError(t, mem.AddTask(&task))

		got, err := mem.GetTask(task.ID)
		require.NoError(t, err)
		assert.Equal(t, task, *got)
	})

	t.Run("not found", func(t *testing.T) {
		resetForTest()
		mem := New()
		_, err := mem.GetTask(666)
		require.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetForTest()
		mem := New()
		task := task.Task{ID: 1, Description: "one", Priority: task.PriorityLow, Created: time.Now()}
		require.NoError(t, mem.AddTask(&task))

		err := mem.DeleteTask(task.ID)
		require.NoError(t, err)

		_, err = mem.GetTask(task.ID)
		require.Error(t, err)

		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("not found", func(t *testing.T) {
		resetForTest()
		mem := New()
		err := mem.DeleteTask(666)
		require.Error(t, err)
	})
}

func TestListTasks(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		resetForTest()
		mem := New()
		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("multiple", func(t *testing.T) {
		resetForTest()
		mem := New()
		tasks := []task.Task{
			{ID: 1, Description: "one", Priority: task.PriorityLow, Created: time.Now()},
			{ID: 2, Description: "two", Priority: task.PriorityMedium, Created: time.Now()},
			{ID: 3, Description: "three", Priority: task.PriorityHigh, Created: time.Now()},
		}
		for _, v := range tasks {
			tsk := v
			require.NoError(t, mem.AddTask(&tsk))
		}

		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, len(tasks))
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetForTest()
		mem := New()
		orig := task.Task{ID: 1, Description: "one", Priority: task.PriorityLow, Created: time.Now()}
		updated := task.Task{ID: 1, Description: "updated", Priority: task.PriorityHigh, Created: orig.Created}

		require.NoError(t, mem.AddTask(&orig))

		err := mem.UpdateTask(&updated)
		require.NoError(t, err)

		got, err := mem.GetTask(orig.ID)
		require.NoError(t, err)
		assert.Equal(t, updated, *got)
	})

	t.Run("not found", func(t *testing.T) {
		resetForTest()
		mem := New()
		nonexistent := task.Task{ID: 666, Description: "nope", Priority: task.PriorityLow, Created: time.Now()}

		err := mem.UpdateTask(&nonexistent)
		require.Error(t, err)
	})
}
