//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Events struct {
	EventID          int32 `sql:"primary_key"`
	EventDescription string
	StartDate        time.Time
	RecurrentCount   int32
	PeriodFreq       int32
	PeriodType       string
	CourseID         int32 `sql:"primary_key"`
}