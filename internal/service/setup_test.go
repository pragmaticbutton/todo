package service

import (
	"time"

	"github.com/pragmaticbutton/todo/internal/domain/list"
	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/storage"

	"github.com/stretchr/testify/mock"
)

// Fixed timestamp for deterministic testing
var fixedTimestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type MockListStorage struct {
	mock.Mock
}

type MockTaskStorage struct {
	mock.Mock
}

var (
	_ storage.TaskStorage = (*MockTaskStorage)(nil)
	_ storage.ListStorage = (*MockListStorage)(nil)
)

func newListServiceWithMocks() (*ListService, *MockListStorage, *MockTaskStorage) {
	mockListStorage := new(MockListStorage)
	mockTaskStorage := new(MockTaskStorage)
	return NewListService(mockListStorage, mockTaskStorage), mockListStorage, mockTaskStorage
}

func newTaskServiceWithMocks() (*TaskService, *MockTaskStorage, *MockListStorage) {
	mockTaskStorage := new(MockTaskStorage)
	mockListStorage := new(MockListStorage)
	return NewTaskService(mockTaskStorage, mockListStorage), mockTaskStorage, mockListStorage
}

func (m *MockTaskStorage) NextTaskID() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockTaskStorage) AddTask(t *task.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTaskStorage) ListTasks() ([]task.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]task.Task), args.Error(1)
}

func (m *MockTaskStorage) GetTask(id uint32) (*task.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task.Task), args.Error(1)
}

func (m *MockTaskStorage) DeleteTask(id uint32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskStorage) UpdateTask(t *task.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTaskStorage) SearchTasks(listID *uint32) ([]task.Task, error) {
	args := m.Called(listID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]task.Task), args.Error(1)
}

func (m *MockListStorage) NextListID() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockListStorage) AddList(l *list.List) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *MockListStorage) ListLists() ([]list.List, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]list.List), args.Error(1)
}

func (m *MockListStorage) GetList(id uint32) (*list.List, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*list.List), args.Error(1)
}

func (m *MockListStorage) DeleteList(id uint32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockListStorage) UpdateList(l *list.List) error {
	args := m.Called(l)
	return args.Error(0)
}
