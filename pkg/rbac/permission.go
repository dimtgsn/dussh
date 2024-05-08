package rbac

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

type Permission struct {
	Method string `json:"method" yaml:"method"`
	routes mapset.Set[string]
}

func NewPermission(method string, routes ...string) *Permission {
	permRoutes := mapset.NewSet[string]()
	for _, r := range routes {
		permRoutes.Add(NormalizeRoute(r))
	}

	return &Permission{
		Method: method,
		routes: permRoutes,
	}
}

func (p *Permission) Routes() mapset.Set[string] {
	return p.routes
}

func (p *Permission) UnmarshalJSON(data []byte) error {
	tempPermission := struct {
		Method string   `json:"method" yaml:"method"`
		Routes []string `json:"routes" yaml:"routes"`
	}{}

	if err := json.Unmarshal(data, &tempPermission); err != nil {
		return err
	}

	p.Method = tempPermission.Method
	p.routes = mapset.NewSet[string]()
	for _, r := range tempPermission.Routes {
		p.routes.Add(NormalizeRoute(r))
	}

	return nil
}

func NormalizeRoute(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "/", "")
}
