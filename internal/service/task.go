package service

import (
	"time"
	"todo/internal/domain/task"
	"todo/internal/storage"
	"todo/internal/utils"
)

type TaskService struct {
	storage storage.Storage
}

type UpdateTaskInput struct {
	Description *string
	Priority    *task.Priority
	Done        *bool
}

func NewTaskService(s storage.Storage) *TaskService {
	return &TaskService{
		storage: s,
	}
}

func (s *TaskService) AddTask(desc string, pr task.Priority) (*task.Task, error) {
	t := task.New(s.storage.NextTaskID(), desc, pr)
	err := s.storage.AddTask(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TaskService) ListTasks() ([]task.Task, error) {
	return s.storage.ListTasks()
}

func (s *TaskService) GetTask(id uint32) (*task.Task, error) {
	return s.storage.GetTask(id)
}

func (s *TaskService) DeleteTask(id uint32) error {
	return s.storage.DeleteTask(id)
}

func (s *TaskService) CompleteTask(id uint32) error {
	_, err := s.UpdateTask(id, UpdateTaskInput{Done: utils.Ptr(true)})
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) ReopenTask(id uint32) error {
	_, err := s.UpdateTask(id, UpdateTaskInput{Done: utils.Ptr(false)})
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) PercentDone() (uint8, error) { // TODO: what type should be returned here?
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

func (s *TaskService) UpdateTask(id uint32, input UpdateTaskInput) (*task.Task, error) {
	t, err := s.storage.GetTask(id)
	if err != nil {
		return nil, err
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if input.Priority != nil {
		t.Priority = *input.Priority
	}
	if input.Done != nil {
		t.Done = *input.Done
	}
	t.Updated = time.Now()
	err = s.storage.UpdateTask(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
