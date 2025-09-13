package memory

import (
	"fmt"
	"sync"
	"todo/internal/domain/list"
	"todo/internal/domain/task"
)

type Memory struct {
	tasks map[uint32]*task.Task
	lists map[uint32]*list.List
}

var data *Memory // TODO: is singleton really needed here? Also, should this be protected by mutex or sync.Map?

var once sync.Once

func New() *Memory {
	if data == nil {
		once.Do(func() {
			data = &Memory{
				tasks: make(map[uint32]*task.Task),
				lists: make(map[uint32]*list.List),
			}
		})

	}
	return data
}

func (m *Memory) AddTask(t *task.Task) error {
	m.tasks[t.ID] = t
	return nil
}

func (m *Memory) ListTasks() ([]task.Task, error) {
	ts := make([]task.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		ts = append(ts, *t)

	}
	return ts, nil
}

func (m *Memory) GetTask(id uint32) (*task.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (m *Memory) DeleteTask(id uint32) error {
	_, ok := m.tasks[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(m.tasks, id)
	return nil
}

func (m *Memory) UpdateTask(t *task.Task) error {
	if _, ok := m.tasks[t.ID]; !ok {
		return fmt.Errorf("task with id %d not found", t.ID)
	}
	m.tasks[t.ID] = t
	return nil
}

func (m *Memory) NextTaskID() uint32 {
	return uint32(len(m.tasks) + 1)
}

func (m *Memory) AddList(l *list.List) error {
	m.lists[l.ID] = l
	return nil
}

func (m *Memory) ListLists() ([]list.List, error) {
	ls := make([]list.List, 0, len(m.lists))
	for _, l := range m.lists {
		ls = append(ls, *l)

	}
	return ls, nil
}

func (m *Memory) GetList(id uint32) (*list.List, error) {
	l, ok := m.lists[id]
	if !ok {
		return nil, fmt.Errorf("list with id %d not found", id)
	}
	return l, nil
}

func (m *Memory) DeleteList(id uint32) error {
	_, ok := m.lists[id]
	if !ok {
		return fmt.Errorf("list with id %d not found", id)
	}
	delete(m.lists, id)
	return nil
}

func (m *Memory) UpdateList(l *list.List) error {
	if _, ok := m.lists[l.ID]; !ok {
		return fmt.Errorf("list with id %d not found", l.ID)
	}
	m.lists[l.ID] = l
	return nil
}

func (m *Memory) NextListID() uint32 {
	return uint32(len(m.lists) + 1)
}
