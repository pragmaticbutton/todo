package storage

import (
	"fmt"
	"sync"
	"todo/internal/task"
)

type Storage struct {
	ts map[uint32]*task.Task
}

var data *Storage

var once sync.Once

func New() *Storage {
	if data == nil {
		once.Do(func() {
			data = &Storage{
				ts: make(map[uint32]*task.Task),
			}
		})

	}
	return data
}

func (s *Storage) AddTask(t *task.Task) error {
	s.ts[t.ID] = t
	return nil
}

func (s *Storage) ListTasks() ([]task.Task, error) {
	ts := make([]task.Task, 0, len(s.ts))
	for _, t := range s.ts {
		ts = append(ts, *t)

	}
	return ts, nil
}

func (s *Storage) GetTask(id uint32) (*task.Task, error) {
	t, ok := s.ts[id]
	if !ok {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (s *Storage) DeleteTask(id uint32) error {
	_, ok := s.ts[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(data.ts, id)
	return nil
}

func (s *Storage) UpdateTask(t *task.Task) error {
	if _, ok := s.ts[t.ID]; !ok {
		return fmt.Errorf("task with id %d not found", t.ID)
	}
	s.ts[t.ID] = t
	return nil
}

func (s *Storage) NextID() uint32 {
	return uint32(len(s.ts) + 1)
}
