package service

import (
	"context"
	"dussh/internal/cache/redis"
	"dussh/internal/domain/models"
	userv1 "dussh/internal/services/user/api/v1"
	"dussh/internal/utils/bytesconv"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetEmployeePosition(ctx context.Context, id int64) (*models.Position, error)
	GetAllUserPositions(ctx context.Context) ([]*models.Position, error)
	SaveUser(ctx context.Context, user *models.User) (int64, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, id int64, user *models.User) error
	DeleteUser(ctx context.Context, id int64) error
}

func NewUserService(
	repository Repository,
	log *zap.Logger,
) userv1.Service {
	return &userService{
		repo: repository,
		log:  log.Named("user.service"),
	}
}

type userService struct {
	repo  Repository
	cache redis.Cache

	log *zap.Logger
}

func (u *userService) Get(ctx context.Context, userID int64) (*models.User, error) {
	userInfo, err := u.repo.GetUserByID(ctx, userID)
	return userInfo, err
}

func (u *userService) GetEmployeePosition(ctx context.Context, id int64) (*models.Position, error) {
	position, err := u.repo.GetEmployeePosition(ctx, id)
	return position, err
}

// TODO: add checking role access for creating

func (u *userService) Create(ctx context.Context, user *models.User) (int64, error) {
	hashPassword, err := bcrypt.GenerateFromPassword(
		bytesconv.StringToBytes(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	user.Password = bytesconv.BytesToString(hashPassword)

	id, err := u.repo.SaveUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *userService) Update(ctx context.Context, id int64, user *models.User) error {
	return u.repo.UpdateUser(ctx, id, user)
}

func (u *userService) Delete(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}

func (u *userService) GetAllPositions(ctx context.Context) ([]*models.Position, error) {
	return u.repo.GetAllUserPositions(ctx)
}
