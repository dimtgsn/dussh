package rbac

import (
	"dussh/internal/role"
	"dussh/pkg/rbac"
	"go.uber.org/zap"
)

type App struct {
	roleManager rbac.RoleManager
}

func New(log *zap.Logger) *App {
	log.Info("rbac app creating")

	roleManager := role.MustNew()

	log.Info("rbac app created")
	return &App{
		roleManager: roleManager,
	}
}

func (a *App) RoleManager() rbac.RoleManager {
	return a.roleManager
}
