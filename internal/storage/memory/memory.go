package memory

import (
	"fmt"
	"sync"
	"todo/internal/task"
)

type Memory struct {
	ts map[uint32]*task.Task
}

var data *Memory

var once sync.Once

func New() *Memory {
	if data == nil {
		once.Do(func() {
			data = &Memory{
				ts: make(map[uint32]*task.Task),
			}
		})

	}
	return data
}

func (m *Memory) AddTask(t *task.Task) error {
	m.ts[t.ID] = t
	return nil
}

func (m *Memory) ListTasks() ([]task.Task, error) {
	ts := make([]task.Task, 0, len(m.ts))
	for _, t := range m.ts {
		ts = append(ts, *t)

	}
	return ts, nil
}

func (m *Memory) GetTask(id uint32) (*task.Task, error) {
	t, ok := m.ts[id]
	if !ok {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (m *Memory) DeleteTask(id uint32) error {
	_, ok := m.ts[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(data.ts, id)
	return nil
}

func (m *Memory) UpdateTask(t *task.Task) error {
	if _, ok := m.ts[t.ID]; !ok {
		return fmt.Errorf("task with id %d not found", t.ID)
	}
	m.ts[t.ID] = t
	return nil
}

func (m *Memory) NextID() uint32 {
	return uint32(len(m.ts) + 1)
}
