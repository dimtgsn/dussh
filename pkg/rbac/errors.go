package rbac

import "errors"

var (
	ErrRoleNotFound          = errors.New("role not found")
	ErrPermissionNotFound    = errors.New("permission not found")
	ErrRolesFileDoesNotExist = errors.New("roles file does not exist")
	ErrCanNotRolesFile       = errors.New("cannot read roles file")
)
