package task

import "time"

type Task struct {
	ID          uint32
	Description string
	Done        bool
	Created     time.Time
	Updated     time.Time // TODO: what about nil time?
}

func New(id uint32, desc string) *Task {
	return &Task{
		ID:          id,
		Description: desc,
		Created:     time.Now(),
		Updated:     time.Time{}, // TODO: what about nil time?
	}
}
