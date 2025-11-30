package memory

import (
	"fmt"
	"sync"

	"github.com/pragmaticbutton/todo/internal/domain/list"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

type memory struct {
	mu    sync.RWMutex
	tasks map[uint32]*task.Task
	lists map[uint32]*list.List
	// nextTaskID/nextListID are monotonic counters to avoid collisions even after deletions.
	nextTaskID uint32
	nextListID uint32
}

func New() *memory {
	return &memory{
		tasks:      make(map[uint32]*task.Task),
		lists:      make(map[uint32]*list.List),
		nextTaskID: 1,
		nextListID: 1,
	}
}

func (m *memory) AddTask(t *task.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID == 0 {
		t.ID = m.nextTaskID
		m.nextTaskID++
	} else if t.ID >= m.nextTaskID {
		m.nextTaskID = t.ID + 1
	}

	if _, ok := m.tasks[t.ID]; ok {
		return fmt.Errorf("task with id %d already exists", t.ID)
	}
	m.tasks[t.ID] = t
	return nil
}

func (m *memory) ListTasks() ([]task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ts := make([]task.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		ts = append(ts, *t)

	}
	return ts, nil
}

func (m *memory) GetTask(id uint32) (*task.Task, error) {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (m *memory) DeleteTask(id uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.tasks[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(m.tasks, id)
	return nil
}

func (m *memory) UpdateTask(t *task.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tasks[t.ID]; !ok {
		return fmt.Errorf("task with id %d not found", t.ID)
	}
	m.tasks[t.ID] = t
	return nil
}

func (m *memory) NextTaskID() uint32 {
	m.mu.Lock()
	id := m.nextTaskID
	m.nextTaskID++
	m.mu.Unlock()
	return id
}

func (m *memory) SearchTasks(listID *uint32) ([]task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ts := []task.Task{}

	for _, t := range m.tasks {
		if listID != nil && t.ListID != nil && *t.ListID == *listID {
			ts = append(ts, *t)
		}
	}

	return ts, nil
}

func (m *memory) AddList(l *list.List) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if l.ID == 0 {
		l.ID = m.nextListID
		m.nextListID++
	} else if l.ID >= m.nextListID {
		m.nextListID = l.ID + 1
	}

	if _, ok := m.lists[l.ID]; ok {
		return fmt.Errorf("list with id %d already exists", l.ID)
	}
	m.lists[l.ID] = l
	return nil
}

func (m *memory) ListLists() ([]list.List, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ls := make([]list.List, 0, len(m.lists))
	for _, l := range m.lists {
		ls = append(ls, *l)

	}
	return ls, nil
}

func (m *memory) GetList(id uint32) (*list.List, error) {
	m.mu.RLock()
	l, ok := m.lists[id]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("list with id %d not found", id)
	}
	return l, nil
}

func (m *memory) DeleteList(id uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.lists[id]
	if !ok {
		return fmt.Errorf("list with id %d not found", id)
	}
	delete(m.lists, id)
	return nil
}

func (m *memory) UpdateList(l *list.List) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.lists[l.ID]; !ok {
		return fmt.Errorf("list with id %d not found", l.ID)
	}
	m.lists[l.ID] = l
	return nil
}

func (m *memory) NextListID() uint32 {
	m.mu.Lock()
	id := m.nextListID
	m.nextListID++
	m.mu.Unlock()
	return id
}
