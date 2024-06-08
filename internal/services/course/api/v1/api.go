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
	CreateEnrollment(ctx context.Context, courseID, userID int64) (int64, error)
	AddEvents(ctx context.Context, courseID int64, events []*models.Event) error
	AddEmployees(ctx context.Context, courseID int64, employees []int64) error
	Update(ctx context.Context, id int64, crs *models.Course) error
	Delete(ctx context.Context, id int64) error
	DeleteEvent(ctx context.Context, courseID, eventID int64) error
	DeleteEmployee(ctx context.Context, courseID, employeeID int64) error
	DeleteEnrollment(ctx context.Context, enrollmentID int64) error
	List(ctx context.Context) ([]*models.Course, error)
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

type AddEventsRequest struct {
	Events []*models.Event `json:"events" validate:"required,dive"`
}

func (ca *courseAPI) AddEvents(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	var req AddEventsRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	if err := ca.svc.AddEvents(c, courseID, req.Events); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"add events to course successfully",
	).OK(c)
}

type AddEmployeesRequest struct {
	Employees []int64 `json:"employees" validate:"required"`
}

func (ca *courseAPI) AddEmployees(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	var req AddEmployeesRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	if err := ca.svc.AddEmployees(c, courseID, req.Employees); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"add employees to course successfully",
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
		"delete event course successfully",
	).OK(c)
}

func (ca *courseAPI) DeleteEmployee(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employee-id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	if err := ca.svc.DeleteEmployee(c, courseID, employeeID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"delete employee course bind successfully",
	).OK(c)
}

type CreateEnrollmentRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

func (ca *courseAPI) CreateEnrollment(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	var req CreateEnrollmentRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if validateErrors := validator.StructValidate(req); validateErrors != nil {
		response.BadRequest(c, validateErrors)
		return
	}

	enrollmentID, err := ca.svc.CreateEnrollment(c, courseID, req.UserID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"create new enrollment to course successfully",
		response.WithValues(map[string]any{"enrollment_id": enrollmentID}),
	).OK(c)
}

func (ca *courseAPI) DeleteEnrollment(c *gin.Context) {
	enrollmentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, domainerrors.ErrInvalidURLPattern)
		return
	}

	if err := ca.svc.DeleteEnrollment(c, enrollmentID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"delete course enrollment successfully",
	).OK(c)
}

func (ca *courseAPI) List(c *gin.Context) {
	courses, err := ca.svc.List(c)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.New(
		http.StatusOK,
		"course list received successfully",
		response.WithValues(map[string]any{
			"courses": courses,
		}),
	).OK(c)
}
