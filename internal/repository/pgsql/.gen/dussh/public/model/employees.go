//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type Employees struct {
	EmployeeID     int32 `sql:"primary_key"`
	PersonalInfoID int32
	PositionID     int32
	CourseID       *int32
	DiplomaID      int32
	DegreeID       *int32
	TitleID        *int32
}
