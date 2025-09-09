package task

import "time"

type Task struct {
	Id          uint32
	Description string
	Done        bool
	Created     time.Time
}

func New(id uint32, desc string) *Task {
	return &Task{
		Id:          id,
		Description: desc,
		Created:     time.Now(),
	}
}
