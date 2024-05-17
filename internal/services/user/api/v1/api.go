package v1

import (
	"context"
	domainerrors "dussh/internal/domain/errors"
	"dussh/internal/domain/models"
	"dussh/internal/domain/response"
	"dussh/internal/services/user"
	"dussh/pkg/validator"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Service interface {
	Get(ctx context.Context, id int64) (*models.User, error)
	GetEmployeePosition(ctx context.Context, id int64) (*models.Position, error)
	GetAllPositions(ctx context.Context) ([]*models.Position, error)
	Create(ctx context.Context, user *models.User) (int64, error)
	Update(ctx context.Context, id int64, user *models.User) error
	Delete(ctx context.Context, id int64) error
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

func (u *userAPI) Get(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	usr, err := u.svc.Get(c, userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	userInfo := models.UserInfo{
		ID:         usr.ID,
		FirstName:  usr.FirstName,
		MiddleName: usr.MiddleName,
		Surname:    usr.Surname,
		Email:      usr.Email,
		Phone:      usr.Phone,
		Role:       usr.Role,
	}

	if usr.Role == models.Employee {
		position, err := u.svc.GetEmployeePosition(c, userID)
		if err == nil && position != nil {
			userInfo.PositionName = position.Name
		}
	}

	response.New(
		http.StatusOK,
		"get user successfully",
		response.WithValues(map[string]any{"user": userInfo}),
	).OK(c)
}

func (u *userAPI) GetAllPositions(c *gin.Context) {
	positions, err := u.svc.GetAllPositions(c)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"get all employee positions successfully",
		response.WithValues(map[string]any{"positions": positions}),
	).OK(c)
}

func (u *userAPI) Create(c *gin.Context) {
	var usr models.User

	if err := c.BindJSON(&usr); err != nil {
		response.BadRequest(c, err)
		return
	}

	if usr.Role == models.Employee && usr.PositionID == 0 {
		response.BadRequest(c, errors.New("employee position id is required field"))
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

type UpdateRequest struct {
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	Phone      string `json:"phone,omitempty" validate:"omitempty,e164"`
}

func (u *userAPI) Update(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	var req UpdateRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	usr := &models.User{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Surname:    req.Surname,
		Email:      req.Email,
		Phone:      req.Phone,
	}

	if err := u.svc.Update(c, userID, usr); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"update user successfully",
	).OK(c)
}

func (u *userAPI) Delete(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	if err := u.svc.Delete(c, userID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"delete user successfully",
	).OK(c)
}
