package service

import (
	"context"
	"dussh/internal/broker/rabbit/publisher"
	"dussh/internal/cache/redis"
	"dussh/internal/domain/models"
	coursev1 "dussh/internal/services/course/api/v1"
	"errors"
	"go.uber.org/zap"
)

type Repository interface {
	GetCourse(ctx context.Context, courseID int64) (*models.Course, error)
	SaveCourse(ctx context.Context, crs *models.Course) (int64, error)
	SaveEvents(ctx context.Context, courseID int64, events []*models.Event) error
	SaveEmployees(ctx context.Context, courseID int64, employees []int64) error
	SaveEnrollment(ctx context.Context, courseID, userID int64) (int64, error)
	UpdateCourse(ctx context.Context, id int64, crs *models.Course) error
	DeleteCourse(ctx context.Context, id int64) error
	DeleteEvent(ctx context.Context, courseID, eventID int64) error
	DeleteEmployee(ctx context.Context, courseID, employeeID int64) error
	DeleteEnrollment(ctx context.Context, enrollmentID int64) error
	CheckCountEvents(ctx context.Context, courseID int64) (int, error)
	CheckCountEmployees(ctx context.Context, courseID int64) (int, error)
}

func NewCourseService(
	repository Repository,
	enrollmentBroker publisher.Publisher[models.EnrollmentEvent],
	log *zap.Logger,
) coursev1.Service {
	return &courseService{
		repo:             repository,
		enrollmentBroker: enrollmentBroker,
		log:              log.Named("course.service"),
	}
}

var (
	ErrMustBeAtLeastOneEvent        = errors.New("there must be at least one event")
	ErrMustBeAtLeastOneBindEmployee = errors.New("there must be at least one bind employee")
)

type courseService struct {
	repo  Repository
	cache redis.Cache
	// TODO добавить отправку в очередь событий
	enrollmentBroker publisher.Publisher[models.EnrollmentEvent]

	log *zap.Logger
}

func (c *courseService) Get(ctx context.Context, courseID int64) (*models.Course, error) {
	return c.repo.GetCourse(ctx, courseID)
}

func (c *courseService) Create(ctx context.Context, crs *models.Course) (int64, error) {
	return c.repo.SaveCourse(ctx, crs)
}

func (c *courseService) CreateEnrollment(ctx context.Context, courseID, userID int64) (int64, error) {
	return c.repo.SaveEnrollment(ctx, courseID, userID)
}

func (c *courseService) AddEvents(ctx context.Context, courseID int64, events []*models.Event) error {
	return c.repo.SaveEvents(ctx, courseID, events)
}

func (c *courseService) AddEmployees(ctx context.Context, courseID int64, employees []int64) error {
	return c.repo.SaveEmployees(ctx, courseID, employees)
}

func (c *courseService) Update(ctx context.Context, id int64, crs *models.Course) error {
	return c.repo.UpdateCourse(ctx, id, crs)
}

func (c *courseService) Delete(ctx context.Context, id int64) error {
	return c.repo.DeleteCourse(ctx, id)
}

func (c *courseService) DeleteEvent(ctx context.Context, courseID, eventID int64) error {
	count, err := c.repo.CheckCountEvents(ctx, courseID)
	if err != nil {
		return err
	}

	if count <= 1 {
		return ErrMustBeAtLeastOneEvent
	}

	return c.repo.DeleteEvent(ctx, courseID, eventID)
}

func (c *courseService) DeleteEmployee(ctx context.Context, courseID, employeeID int64) error {
	count, err := c.repo.CheckCountEmployees(ctx, courseID)
	if err != nil {
		return err
	}

	if count <= 1 {
		return ErrMustBeAtLeastOneBindEmployee
	}

	return c.repo.DeleteEmployee(ctx, courseID, employeeID)
}

func (c *courseService) DeleteEnrollment(ctx context.Context, enrollmentID int64) error {
	return c.repo.DeleteEnrollment(ctx, enrollmentID)
}
