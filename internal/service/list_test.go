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

func TestAddList(t *testing.T) {
	t.Run("successfully adds a list with next ID", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("NextListID").Return(uint32(1))
		mockListStorage.On("AddList", mock.MatchedBy(func(l *list.List) bool {
			return l.ID == 1 && l.Name == "shopping"
		})).Return(nil)

		err := svc.AddList(AddListInput{Name: "shopping"})

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("NextListID").Return(uint32(1))
		mockListStorage.On("AddList", mock.Anything).Return(errors.New("storage error"))

		err := svc.AddList(AddListInput{Name: "shopping"})

		assert.Error(t, err)
		assert.Equal(t, "storage error", err.Error())
	})
}

func TestListLists(t *testing.T) {
	t.Run("successfully lists all lists", func(t *testing.T) {
		lists := []list.List{
			{ID: 1, Name: "shopping", Created: fixedTimestamp},
			{ID: 2, Name: "work", Created: fixedTimestamp},
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("ListLists").Return(lists, nil)

		result, err := svc.ListLists()

		assert.NoError(t, err)
		assert.Equal(t, lists, result)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("ListLists").Return(nil, errors.New("storage error"))

		result, err := svc.ListLists()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "storage error", err.Error())
		mockListStorage.AssertExpectations(t)
	})
}

func TestGetList(t *testing.T) {
	t.Run("successfully gets a list by ID", func(t *testing.T) {
		expectedList := &list.List{ID: 1, Name: "shopping", Created: fixedTimestamp}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(expectedList, nil)

		result, err := svc.GetList(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedList, result)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when list not found", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(nil, errors.New("not found"))

		result, err := svc.GetList(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "not found", err.Error())
		mockListStorage.AssertExpectations(t)
	})
}

func TestDeleteList(t *testing.T) {
	t.Run("successfully deletes a list", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("DeleteList", uint32(1)).Return(nil)

		err := svc.DeleteList(1)

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails to delete", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("DeleteList", uint32(1)).Return(errors.New("storage error"))

		err := svc.DeleteList(1)

		assert.Error(t, err)
		assert.Equal(t, "storage error", err.Error())
		mockListStorage.AssertExpectations(t)
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("successfully updates list with both fields", func(t *testing.T) {
		existingList := &list.List{
			ID:          1,
			Name:        "shopping",
			Description: "old desc",
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(existingList, nil)
		mockListStorage.On("UpdateList", mock.MatchedBy(func(l *list.List) bool {
			return l.ID == 1 &&
				l.Name == "groceries" &&
				l.Description == "new desc" &&
				l.Updated != nil && !l.Updated.IsZero()
		})).Return(nil)

		err := svc.UpdateList(1, &UpdateListInput{
			Name:        utils.Ptr("groceries"),
			Description: utils.Ptr("new desc"),
		})

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("updates only name when description is nil", func(t *testing.T) {
		existingList := &list.List{
			ID:          1,
			Name:        "shopping",
			Description: "original desc",
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(existingList, nil)
		mockListStorage.On("UpdateList", mock.MatchedBy(func(l *list.List) bool {
			return l.Name == "new name" && l.Description == "original desc"
		})).Return(nil)

		err := svc.UpdateList(1, &UpdateListInput{Name: utils.Ptr("new name")})

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("updates only description when name is nil", func(t *testing.T) {
		existingList := &list.List{
			ID:          1,
			Name:        "shopping",
			Description: "old desc",
			Created:     fixedTimestamp,
			Updated:     nil,
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(existingList, nil)
		mockListStorage.On("UpdateList", mock.MatchedBy(func(l *list.List) bool {
			return l.Name == "shopping" && l.Description == "new desc"
		})).Return(nil)

		err := svc.UpdateList(1, &UpdateListInput{Description: utils.Ptr("new desc")})

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("sets Updated timestamp on successful update", func(t *testing.T) {
		existingList := &list.List{
			ID:      1,
			Name:    "shopping",
			Created: fixedTimestamp,
			Updated: nil, // Zero value - not yet updated
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(existingList, nil)
		mockListStorage.On("UpdateList", mock.MatchedBy(func(l *list.List) bool {
			return l.Updated != nil && !l.Updated.IsZero() && l.Updated.After(fixedTimestamp)
		})).Return(nil)

		err := svc.UpdateList(1, &UpdateListInput{Name: utils.Ptr("new")})

		assert.NoError(t, err)
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when GetList fails", func(t *testing.T) {
		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(nil, errors.New("not found"))

		err := svc.UpdateList(1, &UpdateListInput{Name: utils.Ptr("groceries")})

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
		mockListStorage.AssertExpectations(t)
	})

	t.Run("returns error when UpdateList fails", func(t *testing.T) {
		existingList := &list.List{
			ID:      1,
			Name:    "shopping",
			Created: fixedTimestamp,
		}

		svc, mockListStorage, _ := newListServiceWithMocks()
		mockListStorage.On("GetList", uint32(1)).Return(existingList, nil)
		mockListStorage.On("UpdateList", mock.Anything).Return(errors.New("update error"))

		err := svc.UpdateList(1, &UpdateListInput{Description: utils.Ptr("new desc")})

		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
		mockListStorage.AssertExpectations(t)
	})
}

func TestListTasks(t *testing.T) {
	t.Run("successfully lists tasks for a list", func(t *testing.T) {
		tasks := []task.Task{
			{ID: 1, ListID: utils.Ptr(uint32(1)), Description: "desc 1", Done: false, Priority: task.PriorityLow, Created: fixedTimestamp},
			{ID: 2, ListID: utils.Ptr(uint32(1)), Description: "desc 2", Done: true, Priority: task.PriorityHigh, Created: fixedTimestamp},
		}

		svc, _, mockTaskStorage := newListServiceWithMocks()
		mockTaskStorage.On("SearchTasks", utils.Ptr(uint32(1))).Return(tasks, nil)

		result, err := svc.ListTasks(1)

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		mockTaskStorage.AssertExpectations(t)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		svc, _, mockTaskStorage := newListServiceWithMocks()
		mockTaskStorage.On("SearchTasks", utils.Ptr(uint32(1))).Return(nil, errors.New("storage error"))

		result, err := svc.ListTasks(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "storage error", err.Error())
		mockTaskStorage.AssertExpectations(t)
	})
}
