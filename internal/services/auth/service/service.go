package service

import (
	"dussh/internal/cache/redis"
	"dussh/internal/domain/models"
	"dussh/internal/repository"
	"dussh/internal/services/auth/api/v1"
	"dussh/internal/utils/bytesconv"
	"dussh/pkg/jwt"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternal           = errors.New("internal error")
)

type AuthRepository interface {
	SaveUser(ctx context.Context, user *models.User) (int64, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

func NewAuthService(
	repository AuthRepository,
	manager jwt.TokenManager,
	cache redis.Cache,
	log *zap.Logger,
) v1.AuthService {
	return &authService{
		repo:         repository,
		tokenManager: manager,
		cache:        cache,
		log:          log.Named("auth.service"),
	}
}

type authService struct {
	repo         AuthRepository
	tokenManager jwt.TokenManager
	cache        redis.Cache

	log *zap.Logger
}

// Register registers new user and returns user ID.
// If user with given email already exists, returns error.
func (as *authService) Register(ctx context.Context, user *models.User) (int64, error) {
	//exists, err := as.repo.CheckUserExists(ctx, user.Email)
	//if err != nil {
	//	return 0, err
	//}
	//if exists {
	//	return 0, repository.ErrUserAlreadyExists
	//}

	hashPassword, err := bcrypt.GenerateFromPassword(
		bytesconv.StringToBytes(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	user.Password = bytesconv.BytesToString(hashPassword)
	// TODO: при регистрации человек не может быть с ролью кроме guest
	user.Role = models.Guest

	id, err := as.repo.SaveUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Login checks if user with given credentials exists and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (as *authService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		bytesconv.StringToBytes(user.Password),
		bytesconv.StringToBytes(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := as.createJWTTokenPair(user)
	if err != nil {
		return nil, err
	}

	if err := as.cache.SetRefreshToken(
		ctx,
		strconv.Itoa(int(user.ID)),
		tokenPair.RefreshToken.Token,
		time.Duration(tokenPair.RefreshToken.TTL),
	); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (as *authService) createJWTTokenPair(user *models.User) (*jwt.TokenPair, error) {
	accessToken, err := as.tokenManager.NewAccessToken(user.ID, user.Email, int(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := as.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	if accessToken.Token == "" || refreshToken.Token == "" {
		return nil, ErrInternal
	}

	return &jwt.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authService) RefreshToken(ctx context.Context, token string, userID int64) (*jwt.TokenPair, error) {
	usr := strconv.FormatInt(userID, 10)
	t, err := as.cache.GetRefreshTokenByUserID(ctx, usr)
	if err != nil {
		return nil, err
	}

	if t != token {
		return nil, jwt.ErrInvalidToken
	}

	user, err := as.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokenPair, err := as.createJWTTokenPair(user)
	if err != nil {
		return nil, err
	}

	if err := as.cache.UpdateRefreshToken(
		ctx,
		usr,
		tokenPair.RefreshToken.Token,
		time.Duration(tokenPair.RefreshToken.TTL),
	); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (as *authService) Logout(ctx context.Context, token string, userID int64) error {
	usr := strconv.FormatInt(userID, 10)
	t, err := as.cache.GetRefreshTokenByUserID(ctx, usr)
	if err != nil {
		return err
	}

	if t != token {
		return jwt.ErrInvalidToken
	}

	if err := as.cache.DeleteRefreshToken(
		ctx,
		usr,
		token,
	); err != nil {
		return err
	}

	return nil
}
