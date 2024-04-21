package rbac

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGenerateRolesFromFile(t *testing.T) {
	rolesPath := "pkg/rbac/roles.json"
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

	expectedRolesActions := []string{"/auth/logout/", "/auth/register/"}

	var actions []string
	for i, r := range roles {
		assert.Equal(t, r.ID, expectedRoles[i].ID)
		assert.Equal(t, r.Name, expectedRoles[i].Name)

		for _, p := range r.permissions.ToSlice() {
			actions = append(actions, p.actions.ToSlice()...)
		}
	}

	assert.Equal(t, actions, expectedRolesActions)
}
