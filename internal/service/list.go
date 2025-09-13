package service

import "todo/internal/storage"

type ListService struct {
	storage storage.Storage
}

func NewListService(s storage.Storage) *ListService {
	return &ListService{
		storage: s,
	}
}

type CreateListInput struct {
	Description string
}

func (s *ListService) CreateList(input CreateListInput) error {
	return nil
}
