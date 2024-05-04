package v1

import (
	"context"
	domainerrors "dussh/internal/domain/errors"
	"dussh/internal/domain/models"
	"dussh/internal/domain/response"
	"dussh/internal/services/course"
	"dussh/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Service interface {
	Get(ctx context.Context, id int64) (*models.Course, error)
	Create(ctx context.Context, crs *models.Course) (int64, error)
	Update(ctx context.Context, id int64, crs *models.Course) error
	Delete(ctx context.Context, id int64) error
	DeleteEvent(ctx context.Context, courseID, eventID int64) error
}

func NewCourseAPI(service Service, log *zap.Logger) course.Api {
	return &courseAPI{
		svc: service,
		log: log.Named("course.api"),
	}
}

type courseAPI struct {
	svc Service

	log *zap.Logger
}

func (ca *courseAPI) Get(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	crs, err := ca.svc.Get(c, courseID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"get course successfully",
		response.WithValues(map[string]any{"course": crs}),
	).OK(c)
}

func (ca *courseAPI) Create(c *gin.Context) {
	var crs models.Course

	if err := c.BindJSON(&crs); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(crs); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	courseID, err := ca.svc.Create(c, &crs)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"successfully",
		response.WithValues(map[string]any{"course_id": courseID}),
	).OK(c)
}

type UpdateRequest struct {
	Name                    string          `json:"name,omitempty"`
	MonthlySubscriptionCost *float64        `json:"monthly_subscription_cost,omitempty"`
	Events                  []*models.Event `json:"events,omitempty"`
}

func (ca *courseAPI) Update(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

	crs := &models.Course{
		Name:                    req.Name,
		MonthlySubscriptionCost: req.MonthlySubscriptionCost,
		Events:                  req.Events,
	}

	if err := ca.svc.Update(c, courseID, crs); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"update course successfully",
	).OK(c)
}

func (ca *courseAPI) Delete(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	if err := ca.svc.Delete(c, courseID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"delete course successfully",
	).OK(c)
}

func (ca *courseAPI) DeleteEvent(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	eventID, err := strconv.ParseInt(c.Param("event-id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	if err := ca.svc.DeleteEvent(c, courseID, eventID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"delete course successfully",
	).OK(c)
}