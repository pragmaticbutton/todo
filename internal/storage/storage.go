package storage

import (
	"todo/internal/domain/list"
	"todo/internal/domain/task"
)

type TaskStorage interface {
	NextTaskID() uint32
	AddTask(t *task.Task) error
	ListTasks() ([]task.Task, error)
	GetTask(id uint32) (*task.Task, error)
	DeleteTask(id uint32) error
	UpdateTask(t *task.Task) error
	SearchTasks(listID *uint32) ([]task.Task, error)
}

type ListStorage interface {
	NextListID() uint32
	AddList(l *list.List) error
	ListLists() ([]list.List, error)
	GetList(id uint32) (*list.List, error)
	DeleteList(id uint32) error
	UpdateList(l *list.List) error
}

type Storage interface {
	TaskStorage
	ListStorage
}
