package repository

import "errors"

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrEnrollmentAlreadyExists = errors.New("enrollment already exists")
	ErrUserNotFound            = errors.New("user not found")
	ErrPositionsNotFound       = errors.New("employees positions not found")
	ErrCourseNotFound          = errors.New("course not found")
	ErrEventsRequired          = errors.New("events required")
	ErrEmployeesRequired       = errors.New("employees required")
)
