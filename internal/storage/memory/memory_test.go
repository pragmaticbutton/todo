package memory

import (
	"testing"
	"time"
	"todo/internal/domain/list"
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

	t.Run("success, single", func(t *testing.T) {
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

	t.Run("success, multiple", func(t *testing.T) {
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

	t.Run("duplicate id error", func(t *testing.T) {
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
		mem := New()
		task := task.Task{ID: 1, Description: "one", Priority: task.PriorityLow, Created: time.Now()}
		require.NoError(t, mem.AddTask(&task))

		got, err := mem.GetTask(task.ID)
		require.NoError(t, err)
		assert.Equal(t, task, *got)
	})

	t.Run("not found", func(t *testing.T) {
		mem := New()
		_, err := mem.GetTask(666)
		require.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		mem := New()
		err := mem.DeleteTask(666)
		require.Error(t, err)
	})
}

func TestListTasks(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mem := New()
		list, err := mem.ListTasks()
		require.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("multiple", func(t *testing.T) {
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

func TestUpdateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		mem := New()
		nonexistent := task.Task{ID: 666, Description: "nope", Priority: task.PriorityLow, Created: time.Now()}

		err := mem.UpdateTask(&nonexistent)
		require.Error(t, err)
	})
}

func TestSearchTasks(t *testing.T) {
	t.Run("success, found", func(t *testing.T) {
		mem := New()

		// create list
		lst := list.List{ID: 1, Description: "my list", Created: time.Now()}
		err := mem.AddList(&lst)
		require.NoError(t, err)

		// add task in that list
		task := task.Task{ID: 1, Description: "eat", ListID: &lst.ID, Priority: task.PriorityLow, Created: time.Now()}
		require.NoError(t, mem.AddTask(&task))

		results, err := mem.SearchTasks(&lst.ID)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, task, results[0])
	})

	t.Run("success, not found", func(t *testing.T) {
		mem := New()

		// create list
		lst := list.List{ID: 1, Description: "my list", Created: time.Now()}
		err := mem.AddList(&lst)
		require.NoError(t, err)

		// add task in that list
		task := task.Task{ID: 1, Description: "eat", ListID: &lst.ID, Priority: task.PriorityLow, Created: time.Now()}
		require.NoError(t, mem.AddTask(&task))

		// search in a different list
		otherListID := uint32(999)
		results, err := mem.SearchTasks(&otherListID)
		require.NoError(t, err)
		assert.Len(t, results, 0)
	})
}

func TestNextTaskID(t *testing.T) {
	t.Run("sequential", func(t *testing.T) {
		mem := New()
		for i := 1; i <= 5; i++ {
			nextID := mem.NextTaskID()
			assert.Equal(t, uint32(i), nextID)
			task := task.Task{ID: nextID, Description: "task", Priority: task.PriorityLow, Created: time.Now()}
			require.NoError(t, mem.AddTask(&task))
		}
	})

	t.Run("after deletion", func(t *testing.T) {
		mem := New()
		for i := 1; i <= 3; i++ {
			task := task.Task{ID: uint32(i), Description: "task", Priority: task.PriorityLow, Created: time.Now()}
			require.NoError(t, mem.AddTask(&task))
		}
		require.NoError(t, mem.DeleteTask(2))

		nextID := mem.NextTaskID()
		assert.Equal(t, uint32(3), nextID) // still 3, since we had 3 tasks added
	})
}

func TestAddList(t *testing.T) {
	t.Run("success, single", func(t *testing.T) {
		mem := New()

		in := list.List{ID: 1, Description: "my list", Created: time.Now()}
		require.NoError(t, mem.AddList(&in))

		got, err := mem.GetList(in.ID)
		require.NoError(t, err)
		assert.Equal(t, in, *got)

		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("success, multiple", func(t *testing.T) {
		mem := New()

		inputs := []list.List{
			{ID: 1, Description: "list one", Created: time.Now()},
			{ID: 2, Description: "list two", Created: time.Now()},
			{ID: 3, Description: "list three", Created: time.Now()},
		}

		for _, v := range inputs {
			ls := v
			require.NoError(t, mem.AddList(&ls))
		}

		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, len(inputs))

		for _, v := range inputs {
			got, err := mem.GetList(v.ID)
			require.NoError(t, err)
			assert.Equal(t, v, *got)
		}
	})

	t.Run("duplicate id error", func(t *testing.T) {
		mem := New()

		first := list.List{ID: 1, Description: "orig", Created: time.Now()}
		dup := list.List{ID: 1, Description: "dup", Created: time.Now()}

		require.NoError(t, mem.AddList(&first)) // first add OK

		err := mem.AddList(&dup) // second add must fail
		require.Error(t, err)

		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, 1)

		got, err := mem.GetList(first.ID)
		require.NoError(t, err)
		assert.Equal(t, first, *got)
	})
}

func TestListLists(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mem := New()
		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("multiple", func(t *testing.T) {
		mem := New()
		lists := []list.List{
			{ID: 1, Description: "list one", Created: time.Now()},
			{ID: 2, Description: "list two", Created: time.Now()},
			{ID: 3, Description: "list three", Created: time.Now()},
		}
		for _, v := range lists {
			ls := v
			require.NoError(t, mem.AddList(&ls))
		}

		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, len(lists))
	})
}

func TestGetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mem := New()
		lst := list.List{ID: 1, Description: "my list", Created: time.Now()}
		require.NoError(t, mem.AddList(&lst))

		got, err := mem.GetList(lst.ID)
		require.NoError(t, err)
		assert.Equal(t, lst, *got)
	})

	t.Run("not found", func(t *testing.T) {
		mem := New()
		_, err := mem.GetList(666)
		require.Error(t, err)
	})
}

func TestDeleteList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mem := New()
		lst := list.List{ID: 1, Description: "my list", Created: time.Now()}
		require.NoError(t, mem.AddList(&lst))

		err := mem.DeleteList(lst.ID)
		require.NoError(t, err)
		_, err = mem.GetList(lst.ID)
		require.Error(t, err)

		list, err := mem.ListLists()
		require.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("not found", func(t *testing.T) {
		mem := New()
		err := mem.DeleteList(666)
		require.Error(t, err)
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mem := New()
		orig := list.List{ID: 1, Description: "my list", Created: time.Now()}
		updated := list.List{ID: 1, Description: "updated list", Created: orig.Created}

		require.NoError(t, mem.AddList(&orig))
		err := mem.UpdateList(&updated)
		require.NoError(t, err)

		got, err := mem.GetList(orig.ID)
		require.NoError(t, err)
		assert.Equal(t, updated, *got)
	})

	t.Run("not found", func(t *testing.T) {
		mem := New()
		nonexistent := list.List{ID: 666, Description: "nope", Created: time.Now()}

		err := mem.UpdateList(&nonexistent)
		require.Error(t, err)
	})
}

func TestNextListID(t *testing.T) {
	t.Run("sequential", func(t *testing.T) {
		mem := New()
		for i := 1; i <= 5; i++ {
			nextID := mem.NextListID()
			assert.Equal(t, uint32(i), nextID)
			lst := list.List{ID: nextID, Description: "list", Created: time.Now()}
			require.NoError(t, mem.AddList(&lst))
		}
	})

	t.Run("after deletion", func(t *testing.T) {
		mem := New()
		for i := 1; i <= 3; i++ {
			lst := list.List{ID: uint32(i), Description: "list", Created: time.Now()}
			require.NoError(t, mem.AddList(&lst))
		}
		require.NoError(t, mem.DeleteList(2))

		nextID := mem.NextListID()
		assert.Equal(t, uint32(3), nextID) // still 3, since we had 3 lists added
	})
}
