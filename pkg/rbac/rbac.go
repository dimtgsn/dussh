package rbac

import (
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

type RoleManger interface {
	CreateRole(id int, name string, permissions ...*Permission) *Role
	SetPermissions(roleID int, permissions ...*Permission) error
	GetByName(name string) *Role
	GetByID(id int) *Role
	IsGranted(roleID int, permName, route string) (bool, error)
}

type rbac struct {
	roles mapset.Set[*Role]
}

func NewRoleManager(roles ...*Role) RoleManger {
	return &rbac{
		roles: mapset.NewSet[*Role](roles...),
	}
}

func (r *rbac) CreateRole(id int, name string, permissions ...*Permission) *Role {
	role := &Role{
		ID:          id,
		Name:        name,
		permissions: mapset.NewSet[*Permission](permissions...),
	}
	r.roles.Add(role)
	return role
}

func (r *rbac) SetPermissions(id int, permissions ...*Permission) error {
	role := r.GetByID(id)
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

func (r *rbac) GetByID(id int) *Role {
	for _, role := range r.roles.ToSlice() {
		if role.ID == id {
			return role
		}
	}
	return nil
}

func (r *rbac) IsGranted(roleID int, permName, route string) (bool, error) {
	role := r.GetByID(roleID)
	if role == nil {
		return false, ErrRoleNotFound
	}

	permName = strings.ToLower(permName)
	route = NormalizeRoute(route)
	for _, perm := range role.permissions.ToSlice() {
		if permName != "" {
			if strings.ToLower(perm.Name) == permName {
				return perm.routes.Contains(route), nil
			}
		} else {
			if perm.routes.Contains(route) {
				return true, nil
			}
		}
	}

	if permName != "" {
		return false, ErrPermissionNotFound
	}

	return false, nil
}
