package rbac

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
)

type Role struct {
	Name        string `json:"name" yaml:"name"`
	permissions mapset.Set[*Permission]
}

type StorageRoles struct {
	Roles []*Role `json:"roles" yaml:"roles"`
}

func (r *Role) String() string {
	return r.Name
}

func (r *Role) SetPermissions(permissions ...*Permission) {
	r.permissions.Append(permissions...)
}

func (r *Role) Permissions() mapset.Set[*Permission] {
	return r.permissions
}

func (r *Role) UnmarshalJSON(data []byte) error {
	tempRole := struct {
		Name        string        `json:"name" yaml:"name"`
		Permissions []*Permission `json:"permissions" yaml:"permissions"`
	}{}

	if err := json.Unmarshal(data, &tempRole); err != nil {
		return err
	}

	r.Name = tempRole.Name
	r.permissions = mapset.NewSet[*Permission](tempRole.Permissions...)

	return nil
}

func GenerateRolesFromFile(rolesPath string) ([]*Role, error) {
	// check if file exists
	if _, err := os.Stat(rolesPath); os.IsNotExist(err) {
		return nil, ErrRolesFileDoesNotExist
	}

	var roles StorageRoles
	if err := cleanenv.ReadConfig(rolesPath, &roles); err != nil {
		return nil, errors.Wrap(err, ErrCanNotRolesFile.Error())
	}

	return roles.Roles, nil
}
