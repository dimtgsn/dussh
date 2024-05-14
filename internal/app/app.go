package app

import (
	brokerapp "dussh/internal/app/broker"
	cacheapp "dussh/internal/app/cache"
	httpapp "dussh/internal/app/http"
	rbacapp "dussh/internal/app/rbac"
	repoapp "dussh/internal/app/repo"
	"dussh/internal/broker/rabbit/publisher"
	"dussh/internal/config"
	"dussh/internal/domain/models"
	authapi "dussh/internal/services/auth/api/v1"
	authservice "dussh/internal/services/auth/service"
	courseapi "dussh/internal/services/course/api/v1"
	courseservice "dussh/internal/services/course/service"
	"dussh/internal/services/notification"
	userapi "dussh/internal/services/user/api/v1"
	userservice "dussh/internal/services/user/service"
	"dussh/pkg/jwt"
	"dussh/pkg/notify"
	"dussh/pkg/notify/provider/email"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type App struct {
	httpServer *httpapp.App
	broker     *brokerapp.App
	cache      *cacheapp.App
	repo       *repoapp.App
	rbac       *rbacapp.App
}

func New(ctx context.Context, log *zap.Logger, cfg config.Config) *App {
	repoApp := repoapp.New(ctx, &cfg.DB, log)
	cacheApp := cacheapp.New(ctx, &cfg.Redis, log)
	rbacApp := rbacapp.New(log)

	jwtManager, err := jwt.NewManager(
		cfg.Auth.SecretKey,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
	)
	if err != nil {
		panic(err)
	}

	authService := authservice.NewAuthService(
		repoApp.PGSQL(),
		jwtManager,
		cacheApp.Redis(),
		log,
	)
	authAPI := authapi.NewAuthAPI(authService, log)

	userSvc := userservice.NewUserService(repoApp.PGSQL(), log)
	userAPI := userapi.NewUserAPI(userSvc, log)

	ePublisher := publisher.NewEventPublisher[models.EnrollmentEvent](cfg.RabbitMQ)
	courseSvc := courseservice.NewCourseService(repoApp.PGSQL(), ePublisher, log)
	courseAPI := courseapi.NewCourseAPI(courseSvc, log)

	emailCfg := notify.Config{Email: &email.NotificationProvider{
		From:      cfg.Notify.EmailProvider.From,
		Username:  cfg.Notify.EmailProvider.Username,
		Password:  cfg.Notify.EmailProvider.Password,
		Host:      cfg.Notify.EmailProvider.Host,
		Port:      cfg.Notify.EmailProvider.Port,
		TLSEnable: cfg.Notify.EmailProvider.TLSEnable,
	}}
	notificationSvc := notification.NewService(emailCfg, courseSvc, userSvc)

	brokerApp := brokerapp.New(ctx, cfg.RabbitMQ, notificationSvc, log)
	httpApp := httpapp.New(ctx, &cfg, authAPI, userAPI, courseAPI, rbacApp, log)

	return &App{
		httpServer: httpApp,
		broker:     brokerApp,
		cache:      cacheApp,
		repo:       repoApp,
		rbac:       rbacApp,
	}
}

func (a *App) Run() {
	a.broker.MustRun()
	a.httpServer.MustRun()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.cache.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.broker.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
