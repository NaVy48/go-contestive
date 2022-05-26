package entity

import (
	"time"
)

// Contest db model
type Contest struct {
	Entity
	ProblemID   int64
	AuthorID    int64
	ContestName string
	StartTime   time.Time
	EndTime     time.Time
	Problems    []int64
	Users       []int64
}
