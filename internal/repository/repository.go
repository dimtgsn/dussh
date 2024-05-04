package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrCourseNotFound    = errors.New("course not found")
	ErrEventsRequired    = errors.New("events required")
)
