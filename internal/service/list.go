package service

// TODO: should this be in separate package?

import (
	"time"
	"todo/internal/domain/list"
	"todo/internal/domain/task"
	"todo/internal/storage"
	"todo/internal/utils"
)

type ListService struct {
	storage storage.Storage
}

func NewListService(s storage.Storage) *ListService {
	return &ListService{
		storage: s,
	}
}

type AddListInput struct {
	Name        string
	Description string
}

type UpdateListInput struct {
	Name        *string
	Description *string
}

func (l *ListService) AddList(input AddListInput) error {
	list := list.New(l.storage.NextListID(), input.Name, input.Description)
	err := l.storage.AddList(list)
	if err != nil {
		return err
	}
	return nil
}

func (l *ListService) ListLists() ([]list.List, error) {
	return l.storage.ListLists()
}

func (l *ListService) GetList(id uint32) (*list.List, error) {
	return l.storage.GetList(id)
}

func (l *ListService) DeleteList(id uint32) error {
	return l.storage.DeleteList(id)
}

func (l *ListService) UpdateList(id uint32, input *UpdateListInput) error {
	list, err := l.storage.GetList(id)
	if err != nil {
		return err
	}
	if input.Name != nil {
		list.Name = *input.Name
	}
	if input.Description != nil {
		list.Description = *input.Description
	}
	list.Updated = time.Now()
	err = l.storage.UpdateList(list)
	if err != nil {
		return err
	}
	return nil
}

func (l *ListService) ListTasks(id uint32) ([]task.Task, error) {
	ts, err := l.storage.SearchTasks(utils.Ptr(id))
	if err != nil {
		return nil, err
	}
	return ts, nil
}
