package models

type User struct {
	ID         int64  `json:"id" db:"personal_info_id"`
	FirstName  string `json:"first_name" db:"name" validate:"required"`
	MiddleName string `json:"middle_name" db:"middle_name" validate:"required"`
	Surname    string `json:"surname" db:"surname" validate:"required"`
	Email      string `json:"email" db:"email" validate:"required,email"`
	Password   string `db:"password" validate:"required,min=8"`
	Phone      string `json:"phone" db:"phone" validate:"required,e164"`
	Role       Role   `json:"role,omitempty" db:"roles_id"`
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
