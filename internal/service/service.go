package service

import (
	"todo/internal/storage"
	"todo/internal/task"
)

type Service struct {
	storage *storage.Storage
}

func New(s *storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) AddTask(desc string) (*task.Task, error) {
	t := task.New(s.storage.NextID(), desc)
	err := s.storage.AddTask(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTasks() ([]task.Task, error) {
	return s.storage.ListTasks()
}

func (s *Service) GetTask(id uint32) (*task.Task, error) {
	return s.storage.GetTask(id)
}

func (s *Service) DeleteTask(id uint32) error {
	return s.storage.DeleteTask(id)
}

func (s *Service) CompleteTask(id uint32) error {
	t, err := s.storage.GetTask(id)
	if err != nil {
		return err
	}
	t.Done = true
	return s.storage.UpdateTask(t)
}

func (s *Service) ReopenTask(id uint32) error {
	t, err := s.storage.GetTask(id)
	if err != nil {
		return err
	}
	t.Done = false
	return s.storage.UpdateTask(t)
}
