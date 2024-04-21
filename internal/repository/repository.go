package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("the user already exists")
	ErrUserNotFound      = errors.New("the user not found")
)
