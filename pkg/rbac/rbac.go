package rbac

import (
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

type RoleManager interface {
	CreateRole(name string, permissions ...*Permission) *Role
	SetPermissions(role *Role, permissions ...*Permission) error
	GetByName(name string) *Role
	IsGranted(roleName, method, route string) (bool, error)
}

type rbac struct {
	roles mapset.Set[*Role]
}

func NewRoleManager(roles ...*Role) RoleManager {
	return &rbac{
		roles: mapset.NewSet[*Role](roles...),
	}
}

func (r *rbac) CreateRole(name string, permissions ...*Permission) *Role {
	role := &Role{
		Name:        name,
		permissions: mapset.NewSet[*Permission](permissions...),
	}
	r.roles.Add(role)
	return role
}

func (r *rbac) SetPermissions(role *Role, permissions ...*Permission) error {
	if role == nil {
		return ErrRoleNotFound
	}
	role.SetPermissions(permissions...)

	return nil
}

func (r *rbac) GetByName(name string) *Role {
	for _, role := range r.roles.ToSlice() {
		if role.Name == name {
			return role
		}
	}
	return nil
}

func (r *rbac) IsGranted(roleName, method, route string) (bool, error) {
	role := r.GetByName(roleName)
	if role == nil {
		return false, ErrRoleNotFound
	}

	for _, perm := range role.permissions.ToSlice() {
		if methodsIsEqual(perm.Method, method) {
			return perm.routes.Contains(NormalizeRoute(route)), nil
		}
	}

	return false, nil
}

func methodsIsEqual(methodIn, method string) bool {
	return strings.Compare(strings.ToLower(methodIn), strings.ToLower(method)) == 0
}
