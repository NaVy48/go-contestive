package entity

import (
	"database/sql/driver"
	"errors"
)

type SubmissionStatus string

const (
	SubmissionStatusPending SubmissionStatus = "pending"
	SubmissionStatusJudging SubmissionStatus = "judging"
	SubmissionStatusDone    SubmissionStatus = "done"
)

func (s SubmissionStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SubmissionStatus) Scan(value interface{}) error {
	if value == nil {
		return errors.New("SubmissionStatus should not be nil")
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			*s = SubmissionStatus(v)
			return nil
		}
	}
	return errors.New("failed to scan SubmissionStatus")
}

// Submission db model
type Submission struct {
	Entity
	ProblemID    int64
	ProblemRevID int64
	ContestID    int64
	AuthorID     int64
	Language     string
	SourceCode   string
	Status       SubmissionStatus // pending, judging, done
	Result       string           // TLE, WA, RT, CE, AC
	Details      string
}
