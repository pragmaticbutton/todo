package storage

import "todo/internal/task"

type Storage interface {
	NextID() uint32
	AddTask(t *task.Task) error
	ListTasks() ([]task.Task, error)
	GetTask(id uint32) (*task.Task, error)
	DeleteTask(id uint32) error
	UpdateTask(t *task.Task) error
}
