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
	ListID      *uint32 // TODO: what about optional fields?
	Created     time.Time
	Updated     time.Time // TODO: what about nil time?
}

func New(id uint32, desc string, pr Priority, lID *uint32) *Task {
	return &Task{
		ID:          id,
		Description: desc,
		Priority:    pr,
		ListID:      lID,
		Created:     time.Now(),
		Updated:     time.Time{}, // TODO: what about nil time?
	}
}
