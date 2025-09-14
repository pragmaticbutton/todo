package list

import "time"

type List struct {
	ID          uint32
	Name        string
	Description string
	Created     time.Time
	Updated     time.Time // TODO: what about nil time?
}

func New(id uint32, name, desc string) *List {
	return &List{
		ID:          id,
		Name:        name,
		Description: desc,
		Created:     time.Now(),
		Updated:     time.Time{}, // TODO: what about nil time?
	}
}
