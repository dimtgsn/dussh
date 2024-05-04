package models

type User struct {
	ID         int64  `json:"id" db:"personal_info.personal_info_id"`
	FirstName  string `json:"first_name" db:"personal_info.name" validate:"required"`
	MiddleName string `json:"middle_name" db:"personal_info.middle_name" validate:"required"`
	Surname    string `json:"surname" db:"personal_info.surname" validate:"required"`
	Email      string `json:"email" db:"personal_info.email" validate:"required,email"`
	Password   string `json:"password" db:"personal_info.password" validate:"required,min=8"`
	Phone      string `json:"phone" db:"personal_info.phone" validate:"required,e164"`
	Role       Role   `json:"role,omitempty" db:"personal_info.roles_id"`
}

type UserInfo struct {
	ID         int64  `json:"id" db:"personal_info.personal_info_id"`
	FirstName  string `json:"first_name" db:"personal_info.name" `
	MiddleName string `json:"middle_name" db:"personal_info.middle_name" `
	Surname    string `json:"surname" db:"personal_info.surname" `
	Email      string `json:"email" db:"personal_info.email"`
	Phone      string `json:"phone" db:"personal_info.phone"`
	Role       Role   `json:"role,omitempty" db:"personal_info.roles_id"`
}

//go:generate ../../../tools/enumer -type=Role -json -transform=snake
type Role int

const (
	Unspecific Role = iota
	Guest
	Student
	Employee
	Admin
)
