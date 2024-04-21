package rbac

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

type Permission struct {
	Name   string `json:"name" yaml:"name"`
	routes mapset.Set[string]
}

func NewPermission(name string, routes ...string) *Permission {
	permRoutes := mapset.NewSet[string]()
	for _, r := range routes {
		permRoutes.Add(NormalizeRoute(r))
	}

	return &Permission{
		Name:   name,
		routes: permRoutes,
	}
}

func (p *Permission) UnmarshalJSON(data []byte) error {
	tempPermission := struct {
		Name   string   `json:"name" yaml:"name"`
		Routes []string `json:"routes" yaml:"routes"`
	}{}

	if err := json.Unmarshal(data, &tempPermission); err != nil {
		return err
	}

	p.Name = tempPermission.Name
	p.routes = mapset.NewSet[string]()
	for _, r := range tempPermission.Routes {
		p.routes.Add(NormalizeRoute(r))
	}

	return nil
}

func NormalizeRoute(s string) string {
	return strings.ReplaceAll(s, "/", "")
}
