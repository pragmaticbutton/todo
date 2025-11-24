package service

import (
	"errors"
	"testing"
	"todo/internal/domain/list"
	"todo/internal/domain/task"
	"todo/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddTask(t *testing.T) {
	t.Run("successfully adds a task with next ID", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("NextTaskID").Return(uint32(1))
		mockTaskStorage.On("AddTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ID == 1 && tk.Description == "cookies" && tk.Priority == task.PriorityMedium
		})).Return(nil)

		result, err := svc.AddTask(AddTaskInput{
			Description: "cookies",
			Priority:    task.PriorityMedium,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint32(1), result.ID)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("successfully adds a task with valid list", func(t *testing.T) {
		svc, mockTaskStorage, mockListStorage := newTaskServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(&list.List{ID: 1, Name: "shopping"}, nil)
		mockTaskStorage.On("NextTaskID").Return(uint32(1))
		mockTaskStorage.On("AddTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ID == 1 && tk.Description == "cookies" && tk.ListID != nil && *tk.ListID == 1
		})).Return(nil)

		result, err := svc.AddTask(AddTaskInput{
			Description: "cookies",
			Priority:    task.PriorityHigh,
			ListID:      utils.Ptr(uint32(1)),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint32(1), *result.ListID)
		mockListStorage.AssertExpectations(t)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when list does not exist", func(t *testing.T) {
		svc, _, mockListStorage := newTaskServiceWithMocks()
		mockListStorage.On("GetList", uint32(999)).Return(nil, errors.New("list not found"))

		result, err := svc.AddTask(AddTaskInput{
			Description: "cookies",
			Priority:    task.PriorityMedium,
			ListID:      utils.Ptr(uint32(999)),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "list not found", err.Error())
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails to add task", func(t *testing.T) {
		svc, mockTaskStorage, mockListStorage := newTaskServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(&list.List{ID: 1}, nil)
		mockTaskStorage.On("NextTaskID").Return(uint32(1))
		mockTaskStorage.On("AddTask", mock.Anything).Return(errors.New("storage error"))

		result, err := svc.AddTask(AddTaskInput{
			Description: "cookies",
			Priority:    task.PriorityMedium,
			ListID:      utils.Ptr(uint32(1)),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "storage error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestTaskServiceListTasks(t *testing.T) {
	t.Run("successfully lists all tasks", func(t *testing.T) {
		tasks := []task.Task{
			{ID: 1, Description: "task 1", Priority: task.PriorityLow, Done: false, Created: fixedTimestamp},
			{ID: 2, Description: "task 2", Priority: task.PriorityHigh, Done: true, Created: fixedTimestamp},
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(tasks, nil)

		result, err := svc.ListTasks()

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		assert.Equal(t, 2, len(result))
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns empty list when no tasks exist", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return([]task.Task{}, nil)

		result, err := svc.ListTasks()

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(nil, errors.New("storage error"))

		result, err := svc.ListTasks()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "storage error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestGetTask(t *testing.T) {
	t.Run("successfully gets a task by ID", func(t *testing.T) {
		expectedTask := &task.Task{
			ID:          1,
			Description: "cookies",
			Priority:    task.PriorityHigh,
			Done:        false,
			Created:     fixedTimestamp,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(expectedTask, nil)

		result, err := svc.GetTask(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when task not found", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(999)).Return(nil, errors.New("task not found"))

		result, err := svc.GetTask(999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "task not found", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("successfully deletes a task", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("DeleteTask", uint32(1)).Return(nil)

		err := svc.DeleteTask(1)

		assert.NoError(t, err)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails to delete", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("DeleteTask", uint32(1)).Return(errors.New("storage error"))

		err := svc.DeleteTask(1)

		assert.Error(t, err)
		assert.Equal(t, "storage error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestCompleteTask(t *testing.T) {
	t.Run("successfully marks task as complete", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "cookies",
			Priority:    task.PriorityHigh,
			Done:        false,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ID == 1 && tk.Done && tk.Updated != nil && !tk.Updated.IsZero()
		})).Return(nil)

		err := svc.CompleteTask(1)

		assert.NoError(t, err)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when task not found", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(999)).Return(nil, errors.New("task not found"))

		err := svc.CompleteTask(999)

		assert.Error(t, err)
		assert.Equal(t, "task not found", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when update fails", func(t *testing.T) {
		existingTask := &task.Task{
			ID:      1,
			Done:    false,
			Created: fixedTimestamp,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.Anything).Return(errors.New("update error"))

		err := svc.CompleteTask(1)

		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestReopenTask(t *testing.T) {
	t.Run("successfully marks task as reopened", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "cookies",
			Priority:    task.PriorityHigh,
			Done:        true,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ID == 1 && !tk.Done && tk.Updated != nil && !tk.Updated.IsZero()
		})).Return(nil)

		err := svc.ReopenTask(1)

		assert.NoError(t, err)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when task not found", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(999)).Return(nil, errors.New("task not found"))

		err := svc.ReopenTask(999)

		assert.Error(t, err)
		assert.Equal(t, "task not found", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when update fails", func(t *testing.T) {
		existingTask := &task.Task{
			ID:      1,
			Done:    true,
			Created: fixedTimestamp,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.Anything).Return(errors.New("update error"))

		err := svc.ReopenTask(1)

		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("successfully updates all task fields", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "old description",
			Priority:    task.PriorityLow,
			Done:        false,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ID == 1 &&
				tk.Description == "new description" &&
				tk.Priority == task.PriorityHigh &&
				tk.Done &&
				tk.Updated != nil && !tk.Updated.IsZero()
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Description: utils.Ptr("new description"),
			Priority:    utils.Ptr(task.PriorityHigh),
			Done:        utils.Ptr(true),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new description", result.Description)
		assert.Equal(t, task.PriorityHigh, result.Priority)
		assert.True(t, result.Done)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("updates only description when other fields are nil", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "old description",
			Priority:    task.PriorityLow,
			Done:        false,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.Description == "new description" &&
				tk.Priority == task.PriorityLow &&
				!tk.Done
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Description: utils.Ptr("new description"),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new description", result.Description)
		assert.Equal(t, task.PriorityLow, result.Priority)
		assert.False(t, result.Done)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("updates only priority when other fields are nil", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "description",
			Priority:    task.PriorityLow,
			Done:        false,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.Description == "description" &&
				tk.Priority == task.PriorityHigh &&
				!tk.Done
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Priority: utils.Ptr(task.PriorityHigh),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "description", result.Description)
		assert.Equal(t, task.PriorityHigh, result.Priority)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("updates only done status when other fields are nil", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "description",
			Priority:    task.PriorityMedium,
			Done:        false,
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.Description == "description" &&
				tk.Priority == task.PriorityMedium &&
				tk.Done
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Done: utils.Ptr(true),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Done)
		assert.Equal(t, task.PriorityMedium, result.Priority)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("updates ListID when provided", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "description",
			Priority:    task.PriorityMedium,
			Done:        false,
			ListID:      utils.Ptr(uint32(1)),
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.ListID != nil && *tk.ListID == 2
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			ListID: utils.Ptr(uint32(2)),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint32(2), *result.ListID)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("sets Updated timestamp on update", func(t *testing.T) {
		existingTask := &task.Task{
			ID:          1,
			Description: "description",
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.MatchedBy(func(tk *task.Task) bool {
			return tk.Updated != nil && !tk.Updated.IsZero() && tk.Updated.After(fixedTimestamp)
		})).Return(nil)

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Done: utils.Ptr(true),
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.Updated)
		assert.False(t, result.Updated.IsZero())
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when task not found", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(999)).Return(nil, errors.New("task not found"))

		result, err := svc.UpdateTask(999, UpdateTaskInput{
			Done: utils.Ptr(true),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "task not found", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails to update", func(t *testing.T) {
		existingTask := &task.Task{
			ID:      1,
			Created: fixedTimestamp,
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("GetTask", uint32(1)).Return(existingTask, nil)
		mockTaskStorage.On("UpdateTask", mock.Anything).Return(errors.New("update error"))

		result, err := svc.UpdateTask(1, UpdateTaskInput{
			Done: utils.Ptr(true),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "update error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}

func TestPercentDone(t *testing.T) {
	t.Run("returns 0 when no tasks exist", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return([]task.Task{}, nil)

		result, err := svc.PercentDone()

		assert.NoError(t, err)
		assert.Equal(t, uint8(0), result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns 0 when no tasks are done", func(t *testing.T) {
		tasks := []task.Task{
			{ID: 1, Description: "task 1", Done: false, Created: fixedTimestamp},
			{ID: 2, Description: "task 2", Done: false, Created: fixedTimestamp},
			{ID: 3, Description: "task 3", Done: false, Created: fixedTimestamp},
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(tasks, nil)

		result, err := svc.PercentDone()

		assert.NoError(t, err)
		assert.Equal(t, uint8(0), result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns 100 when all tasks are done", func(t *testing.T) {
		tasks := []task.Task{
			{ID: 1, Description: "task 1", Done: true, Created: fixedTimestamp},
			{ID: 2, Description: "task 2", Done: true, Created: fixedTimestamp},
			{ID: 3, Description: "task 3", Done: true, Created: fixedTimestamp},
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(tasks, nil)

		result, err := svc.PercentDone()

		assert.NoError(t, err)
		assert.Equal(t, uint8(100), result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("calculates correct percentage with mixed completion", func(t *testing.T) {
		// 2 out of 4 done = 50%
		tasks := []task.Task{
			{ID: 1, Description: "task 1", Done: true, Created: fixedTimestamp},
			{ID: 2, Description: "task 2", Done: false, Created: fixedTimestamp},
			{ID: 3, Description: "task 3", Done: true, Created: fixedTimestamp},
			{ID: 4, Description: "task 4", Done: false, Created: fixedTimestamp},
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(tasks, nil)

		result, err := svc.PercentDone()

		assert.NoError(t, err)
		assert.Equal(t, uint8(50), result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("calculates correct percentage with one task done", func(t *testing.T) {
		// 1 out of 3 done = 33%
		tasks := []task.Task{
			{ID: 1, Description: "task 1", Done: true, Created: fixedTimestamp},
			{ID: 2, Description: "task 2", Done: false, Created: fixedTimestamp},
			{ID: 3, Description: "task 3", Done: false, Created: fixedTimestamp},
		}

		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(tasks, nil)

		result, err := svc.PercentDone()

		assert.NoError(t, err)
		assert.Equal(t, uint8(33), result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		svc, mockTaskStorage, _ := newTaskServiceWithMocks()
		mockTaskStorage.On("ListTasks").Return(nil, errors.New("storage error"))

		result, err := svc.PercentDone()

		assert.Error(t, err)
		assert.Equal(t, uint8(0), result)
		assert.Equal(t, "storage error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}
