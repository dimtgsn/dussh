package v1

import (
	"context"
	"dussh/internal/domain/models"
	"dussh/internal/domain/response"
	"dussh/internal/services/user"
	"dussh/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	Get(ctx context.Context, userID int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) (int64, error)
}

func NewUserAPI(service Service, log *zap.Logger) user.Api {
	return &userAPI{
		svc: service,
		log: log.Named("user.api"),
	}
}

type userAPI struct {
	svc Service

	log *zap.Logger
}

type GetRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

func (u *userAPI) Get(c *gin.Context) {
	var usr models.User

	if err := c.ShouldBind(&usr); err != nil {
		response.BadRequest(c, err)
		return
	}

	userInfo, err := u.svc.Get(c, usr.ID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"successfully",
		response.WithValues(map[string]any{"user": userInfo}),
	).OK(c)
}

func (u *userAPI) Create(c *gin.Context) {
	var usr models.User

	if err := c.BindJSON(&usr); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(usr); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	userID, err := u.svc.Create(c, &usr)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"successfully",
		response.WithValues(map[string]any{"user_id": userID}),
	).OK(c)
}

func (u *userAPI) Update(c *gin.Context) {}

func (u *userAPI) Delete(c *gin.Context) {}
