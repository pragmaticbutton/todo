package service

import (
	"time"
	"todo/internal/storage"
	"todo/internal/task"
)

type Service struct {
	storage storage.Storage
}

func New(s storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) AddTask(desc string, pr task.Priority) (*task.Task, error) {
	t := task.New(s.storage.NextID(), desc, pr)
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
	t.Updated = time.Now()
	return s.storage.UpdateTask(t)
}

func (s *Service) ReopenTask(id uint32) error {
	t, err := s.storage.GetTask(id)
	if err != nil {
		return err
	}
	t.Done = false
	t.Updated = time.Now()
	return s.storage.UpdateTask(t)
}

func (s *Service) PercentDone() (uint8, error) { // TODO: what type should be returned here?
	tasks, err := s.storage.ListTasks()
	if err != nil {
		return 0, err
	}
	if len(tasks) == 0 {
		return 0, nil
	}

	var doneCount int
	for _, t := range tasks {
		if t.Done {
			doneCount++
		}
	}
	return uint8((doneCount * 100) / len(tasks)), nil
}
