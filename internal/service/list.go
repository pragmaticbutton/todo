package service

// TODO: should this be in separate package?

import (
	"time"

	"github.com/pragmaticbutton/todo/internal/domain/list"
	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/storage"
	"github.com/pragmaticbutton/todo/internal/utils"
)

type ListService struct {
	listStorage storage.ListStorage
	taskStorage storage.TaskStorage
}

func NewListService(ls storage.ListStorage, ts storage.TaskStorage) *ListService {
	return &ListService{
		listStorage: ls,
		taskStorage: ts,
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
	list := list.New(l.listStorage.NextListID(), input.Name, input.Description)
	err := l.listStorage.AddList(list)
	if err != nil {
		return err
	}
	return nil
}

func (l *ListService) ListLists() ([]list.List, error) {
	return l.listStorage.ListLists()
}

func (l *ListService) GetList(id uint32) (*list.List, error) {
	return l.listStorage.GetList(id)
}

func (l *ListService) DeleteList(id uint32) error {
	return l.listStorage.DeleteList(id)
}

func (l *ListService) UpdateList(id uint32, input *UpdateListInput) error {
	list, err := l.listStorage.GetList(id)
	if err != nil {
		return err
	}
	if input.Name != nil {
		list.Name = *input.Name
	}
	if input.Description != nil {
		list.Description = *input.Description
	}
	list.Updated = utils.Ptr(time.Now())
	err = l.listStorage.UpdateList(list)
	if err != nil {
		return err
	}
	return nil
}

func (l *ListService) ListTasks(id uint32) ([]task.Task, error) {
	ts, err := l.taskStorage.SearchTasks(utils.Ptr(id))
	if err != nil {
		return nil, err
	}
	return ts, nil
}
