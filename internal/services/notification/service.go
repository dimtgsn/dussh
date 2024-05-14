package notification

import (
	"bytes"
	"context"
	"dussh/internal/domain/models"
	coursev1 "dussh/internal/services/course/api/v1"
	userv1 "dussh/internal/services/user/api/v1"
	"dussh/pkg/notify"
	"dussh/pkg/notify/notification"
	"errors"
	"html/template"
	"log"
	"strings"
)

var ErrNotificationConfigIsInvalid = errors.New("notification config is invalid")

type Service interface {
	Notify(context.Context, *notification.Notification) error
	CreateNotificationByEnrollmentEvent(context.Context, models.EnrollmentEvent) (*notification.Notification, error)
}

func NewService(
	cfg notify.Config,
	courseSvc coursev1.Service,
	userSvc userv1.Service,
) Service {
	return &service{
		cfg:       cfg,
		courseSvc: courseSvc,
		userSvc:   userSvc,
	}
}

type service struct {
	cfg       notify.Config
	courseSvc coursev1.Service
	userSvc   userv1.Service
}

const enrollmentSubject = "Новая запсь на курс"

type enrollmentInfo struct {
	Subject    string
	CourseName string
	User       string
}

func (s *service) CreateNotificationByEnrollmentEvent(
	ctx context.Context,
	e models.EnrollmentEvent,
) (*notification.Notification, error) {
	course, err := s.courseSvc.Get(ctx, e.CourseID)
	if err != nil {
		return nil, err
	}

	user, err := s.userSvc.Get(ctx, e.UserID)
	if err != nil {
		return nil, err
	}

	info := enrollmentInfo{
		Subject:    enrollmentSubject,
		CourseName: course.Name,
		User:       strings.Join([]string{user.Surname, user.FirstName, user.MiddleName}, " "),
	}

	t := template.New("email.html")

	t, err = t.ParseFiles("../../domain/template/email.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, info); err != nil {
		log.Println(err)
	}

	n := &notification.Notification{
		Type:        notification.TypeEmail,
		ContentType: notification.ContentTypeHTML,
		To:          []string{user.Email},
		Subject:     enrollmentSubject,
		Body:        tpl.String(),
	}

	return n, nil
}

func (s *service) Notify(ctx context.Context, n *notification.Notification) error {
	provider := s.cfg.GetNotificationProviderByType(n.Type)
	if !provider.IsValid() {
		return ErrNotificationConfigIsInvalid
	}

	return provider.Send(ctx, n)
}
