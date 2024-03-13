package httpapp

import (
	httpserver "dussh/internal/http"
	"dussh/internal/services/auth"
	authapi "dussh/internal/services/auth/api"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type App struct {
	httpServer *http.Server
	addr       string
}

// New creates new http app.
func New(
	addr string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
) *App {
	router := httpserver.NewRouter()

	baseRouteGroup, err := httpserver.SetAPIPath(router)
	if err != nil {
		panic(err)
	}

	authService :=
	authAPI := authapi.NewAuthAPI()
	auth.Routes(baseRouteGroup, authAPI)

	server := httpserver.NewServer(
		addr,
		router,
		readTimeout,
		writeTimeout,
		idleTimeout,
	)

	return &App{
		httpServer: server,
		addr:       addr,
	}
}

func (a *App) MustRun() {
	if err := a.httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
