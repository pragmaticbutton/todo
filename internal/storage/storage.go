package storage

import (
	"fmt"
	"sync"
	"time"
	"todo/internal/task"
)

type storage struct {
	ts map[uint32]*task.Task
}

var data *storage

var once sync.Once

func New() *storage {
	if data == nil {
		once.Do(func() {
			data = &storage{
				ts: make(map[uint32]*task.Task),
			}
		})

	}
	return data
}

func (s *storage) AddTask(desc string) error {
	id := generateID()
	s.ts[id] = &task.Task{
		ID:          id,
		Description: desc,
		Created:     time.Now(),
	}
	return nil
}

func (s *storage) ListTasks() ([]task.Task, error) {
	ts := make([]task.Task, 0, len(s.ts))
	for _, t := range s.ts {
		ts = append(ts, *t)

	}
	return ts, nil
}

func (s *storage) GetTask(id uint32) (*task.Task, error) {
	t, ok := s.ts[id]
	if !ok {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (s *storage) DeleteTask(id uint32) error {
	_, ok := s.ts[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(data.ts, id)
	return nil
}

func (s *storage) UpdateTask(t *task.Task) error {
	if _, ok := s.ts[t.ID]; !ok {
		return fmt.Errorf("task with id %d not found", t.ID)
	}
	s.ts[t.ID] = t
	return nil
}

func generateID() uint32 {
	return uint32(len(data.ts) + 1)
}
