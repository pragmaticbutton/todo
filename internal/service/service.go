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
