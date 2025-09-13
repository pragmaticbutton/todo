package task

import "time"

type Priority int

const (
	PriorityUnknown Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

type Task struct {
	ID          uint32
	Description string
	Done        bool
	Priority    Priority
	Created     time.Time
	Updated     time.Time // TODO: what about nil time?
}

func New(id uint32, desc string, pr Priority) *Task {
	return &Task{
		ID:          id,
		Description: desc,
		Priority:    pr,
		Created:     time.Now(),
		Updated:     time.Time{}, // TODO: what about nil time?
	}
}
