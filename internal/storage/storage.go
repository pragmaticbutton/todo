package storage

import (
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
		Id:          id,
		Description: desc,
		Created:     time.Now(),
	}
	return nil
}

func (s *storage) ListTasks() ([]task.Task, error) {
	ts := []task.Task{}
	for _, t := range s.ts {
		ts = append(ts, *t)
	}
	return ts, nil
}

func generateID() uint32 {
	return uint32(len(data.ts) + 1)
}
