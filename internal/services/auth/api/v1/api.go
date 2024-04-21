package v1

import (
	"context"
	"dussh/internal/domain/models"
	"dussh/internal/domain/response"
	"dussh/internal/repository"
	"dussh/internal/services/auth"
	"dussh/pkg/jwt"
	"dussh/pkg/validator"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AuthService interface {
	Register(ctx context.Context, user *models.User) (int64, error)
	Login(ctx context.Context, email, password string) (*jwt.TokenPair, error)
	RefreshToken(ctx context.Context, token string, userID int64) (*jwt.TokenPair, error)
	Logout(ctx context.Context, token string, userID int64) error
}

func NewAuthAPI(authService AuthService, log *zap.Logger) auth.Api {
	return &authAPI{
		srv: authService,
		log: log.Named("auth.api"),
	}
}

type authAPI struct {
	srv AuthService

	log *zap.Logger
}

func (a *authAPI) Register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(user); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	userID, err := a.srv.Register(c, &user)
	if err != nil {
		r := response.New(http.StatusInternalServerError, err.Error())
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			r.Code = http.StatusBadRequest
		}

		r.Error(c)
		return
	}

	response.New(
		http.StatusOK,
		"user registered successful",
		response.WithValues(map[string]any{"user_id": userID}),
	).OK(c)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (a *authAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	tokenPair, err := a.srv.Login(c, req.Email, req.Password)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	setRefreshTokenToCookie(c, tokenPair.RefreshToken)
	response.New(
		http.StatusOK,
		"user logged successful",
		response.WithValues(map[string]any{"token": tokenPair.AccessToken.Token}),
	).OK(c)
}

func setRefreshTokenToCookie(c *gin.Context, token *jwt.Token) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"RefreshToken",
		token.Token,
		token.TTL, "", "", false, true) // TODO: change path and domain
}

type RefreshTokenRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

func (a *authAPI) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	token, err := jwt.GetRefreshToken(c)
	if err != nil {
		response.New(http.StatusUnauthorized, err.Error()).Error(c)
		return
	}

	tokenPair, err := a.srv.RefreshToken(c, token, req.UserID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	setRefreshTokenToCookie(c, tokenPair.RefreshToken)
	response.New(
		http.StatusOK,
		"refresh token successful",
		response.WithValues(map[string]any{"token": tokenPair.AccessToken.Token}),
	).OK(c)
}

type LogoutRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

func (a *authAPI) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	token, err := jwt.GetRefreshToken(c)
	if err != nil {
		response.New(http.StatusUnauthorized, err.Error()).Error(c)
		return
	}

	if err := a.srv.Logout(c, token, req.UserID); err != nil {
		response.InternalError(c, err)
		return
	}

	setRefreshTokenToCookie(c, &jwt.Token{
		Token: "",
		TTL:   -1,
	})
	response.New(
		http.StatusOK,
		"user logout successful",
	).OK(c)
}
