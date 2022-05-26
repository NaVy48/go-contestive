package entity

import "time"

type ID = int64

// Entity base with mandatory DB fields
type Entity struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
}
