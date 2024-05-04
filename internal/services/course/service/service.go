package service

import (
	"context"
	"dussh/internal/cache/redis"
	"dussh/internal/domain/models"
	coursev1 "dussh/internal/services/course/api/v1"
	"go.uber.org/zap"
)

type Repository interface {
	GetCourse(ctx context.Context, courseID int64) (*models.Course, error)
	SaveCourse(ctx context.Context, crs *models.Course) (int64, error)
	UpdateCourse(ctx context.Context, id int64, crs *models.Course) error
	DeleteCourse(ctx context.Context, id int64) error
	DeleteEvent(ctx context.Context, courseID, eventID int64) error
}

func NewCourseService(
	repository Repository,
	log *zap.Logger,
) coursev1.Service {
	return &courseService{
		repo: repository,
		log:  log.Named("course.service"),
	}
}

type courseService struct {
	repo  Repository
	cache redis.Cache

	log *zap.Logger
}

func (c *courseService) Get(ctx context.Context, courseID int64) (*models.Course, error) {
	return c.repo.GetCourse(ctx, courseID)
}

func (c *courseService) Create(ctx context.Context, crs *models.Course) (int64, error) {
	return c.repo.SaveCourse(ctx, crs)
}

func (c *courseService) Update(ctx context.Context, id int64, crs *models.Course) error {
	return c.repo.UpdateCourse(ctx, id, crs)
}

func (c *courseService) Delete(ctx context.Context, id int64) error {
	return c.repo.DeleteCourse(ctx, id)
}

func (c *courseService) DeleteEvent(ctx context.Context, courseID, eventID int64) error {
	return c.repo.DeleteEvent(ctx, courseID, eventID)
}
