package list

import "time"

type List struct {
	ID          uint32
	Description string
	Created     time.Time
	Updated     time.Time // TODO: what about nil time?
}

func New(id uint32, desc string) *List {
	return &List{
		ID:          id,
		Description: desc,
		Created:     time.Now(),
		Updated:     time.Time{}, // TODO: what about nil time?
	}
}
