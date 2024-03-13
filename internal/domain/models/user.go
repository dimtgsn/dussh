package models

type User struct {
	Email string `json:"email" db:"email" validate:"required,email"`
}
