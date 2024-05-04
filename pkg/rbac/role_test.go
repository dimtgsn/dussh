package rbac

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGenerateRolesFromFile(t *testing.T) {
	rolesPath := "roles.json"
	roles, err := GenerateRolesFromFile(rolesPath)
	if err != nil {
		t.Error(err)
		return
	}

	expectedRoles := []*Role{
		{
			ID:   1,
			Name: "guest",
		},
	}

	expectedRolesRoutes := []string{"apiv1users"}

	var routes []string
	for i, r := range roles {
		assert.Equal(t, r.ID, expectedRoles[i].ID)
		assert.Equal(t, r.Name, expectedRoles[i].Name)

		for _, p := range r.permissions.ToSlice() {
			routes = append(routes, p.routes.ToSlice()...)
		}
	}

	assert.Equal(t, routes, expectedRolesRoutes)
}
