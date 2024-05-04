package models

import (
	"encoding/json"
	"time"
)

type Course struct {
	ID                      int64    `json:"id" db:"courses.course_id"`
	Name                    string   `json:"name" db:"courses.course_name" validate:"required"`
	MonthlySubscriptionCost *float64 `json:"monthly_subscription_cost" db:"courses.monthly_subscription_cost" validate:"required,min=0"`
	Events                  []*Event `json:"events" db:"events" validate:"required,dive"`
}

type MyTime time.Time

type Event struct {
	ID             int64       `json:"id" db:"events.event_id"`
	Description    string      `json:"description" db:"events.event_description" validate:"required"`
	StartDate      *MyTime     `json:"start_date" db:"events.start_date" validate:"required"`
	RecurrentCount *int64      `json:"recurrent_count" db:"events.recurrent_count" validate:"required,min=1"`
	PeriodFreq     *int64      `json:"period_freq" db:"events.period_freq" validate:"required,min=1"`
	PeriodType     *PeriodType `json:"period_type" db:"events.period_type" validate:"required"`
	CourseID       int64       `json:"course_id" db:"events.course_id"`
}

func (mt *MyTime) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*mt = MyTime(t)
	return nil
}

func (mt *MyTime) MarshalJSON() ([]byte, error) {
	timestamp := time.Time(*mt)
	return json.Marshal(timestamp)
}

//go:generate ../../../tools/enumer -type=PeriodType -json -transform=snake
type PeriodType int

const (
	Day PeriodType = iota + 1
	Month
	Year
)
