package role

import (
	"dussh/pkg/rbac"
	"errors"
	"flag"
	"os"
)

var (
	ErrForbidden = errors.New("access is denied")
)

func MustNew() rbac.RoleManager {
	path := fetchRolesPath()

	roles, err := rbac.GenerateRolesFromFile(path)
	if err != nil {
		panic(err)
	}

	return rbac.NewRoleManager(roles...)
}

func fetchRolesPath() string {
	var path string

	flag.StringVar(&path, "roles", "", "path to roles file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("ROLES_PATH")
	}

	return path
}
