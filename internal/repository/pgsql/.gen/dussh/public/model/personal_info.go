//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type PersonalInfo struct {
	PersonalInfoID int32 `sql:"primary_key"`
	CredsID        int32
	Name           string
	MiddleName     *string
	Surname        string
	Email          string
	RolesID        int32
	Phone          *string
}