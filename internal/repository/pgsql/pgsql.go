package pgsql

import (
	"context"
	"dussh/internal/config"
	"dussh/internal/domain/models"
	"dussh/internal/repository"
	"dussh/internal/repository/pgsql/.gen/dussh/public/table"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type Repository struct {
	db *pgxpool.Pool

	log *zap.Logger
}

func NewPGSQLRepository(ctx context.Context, dbConf *config.DB, log *zap.Logger) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connectionString(dbConf))
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:  pool,
		log: log.Named("repository.pgsql"),
	}, nil
}

func connectionString(dbConf *config.DB) string {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.DBName,
	)
	return connString
}

func (r *Repository) SaveUser(ctx context.Context, user *models.User) (int64, error) {
	r.log.Debug("saving user")
	var (
		credID       int64
		userID       int64
		personalInfo = table.PersonalInfo
		creds        = table.Creds
	)

	if err := withTx(ctx, r.db, func(tx pgx.Tx) error {
		query, args := creds.INSERT(creds.HashedPassword).VALUES(user.Password).RETURNING(creds.CredsID).Sql()
		if err := tx.QueryRow(ctx, query, args...).Scan(&credID); err != nil {
			r.log.Error("failed to create user creds", zap.Error(err))
			return err
		}

		query, args = personalInfo.
			INSERT(personalInfo.Name, personalInfo.MiddleName, personalInfo.Surname,
				personalInfo.Email, personalInfo.Phone, personalInfo.RolesID,
				personalInfo.CredsID,
			).
			VALUES(user.FirstName, user.MiddleName, user.Surname,
				user.Email, user.Phone, user.Role, credID,
			).RETURNING(personalInfo.PersonalInfoID).Sql()

		if err := tx.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
			var pgErr *pgconn.PgError

			r.log.Error("failed to save user", zap.Error(err))
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return repository.ErrUserAlreadyExists
			}

			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}

	r.log.Debug("user saved", zap.Int64("user", userID))
	return userID, nil
}

func (r *Repository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	r.log.Debug("checking user exists")

	personalInfo := table.PersonalInfo
	query, args := personalInfo.SELECT(
		personalInfo.PersonalInfoID,
	).
		WHERE(personalInfo.Email.EQ(postgres.String(email))).
		Sql()

	if err := r.db.QueryRow(ctx, query, args...).Scan(nil); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.log.Debug("getting user by email")

	var user models.User
	personalInfo := table.PersonalInfo

	query, args := postgres.SELECT(
		personalInfo.AllColumns.Except(personalInfo.CredsID),
		table.Creds.HashedPassword.AS("personal_info.password"),
	).
		FROM(personalInfo.INNER_JOIN(table.Creds, personalInfo.CredsID.EQ(table.Creds.CredsID))).
		WHERE(personalInfo.Email.EQ(postgres.String(email))).
		Sql()

	if err := pgxscan.Get(ctx, r.db, &user, query, args...); err != nil {
		r.log.Debug("failed to get user by email", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	r.log.Debug("getting user by id")

	var user models.User
	personalInfo := table.PersonalInfo

	query, args := postgres.SELECT(
		personalInfo.AllColumns.Except(personalInfo.CredsID),
		table.Creds.HashedPassword.AS("personal_info.password"),
	).
		FROM(personalInfo.INNER_JOIN(table.Creds, personalInfo.CredsID.EQ(table.Creds.CredsID))).
		WHERE(personalInfo.PersonalInfoID.EQ(postgres.Int(id))).
		Sql()

	if err := pgxscan.Get(ctx, r.db, &user, query, args...); err != nil {
		r.log.Debug("failed to get user by id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int64, user *models.User) error {
	r.log.Debug("updating user by id")

	updatePersonalInfo := table.PersonalInfo.UPDATE()
	var columns []any
	if user.FirstName != "" {
		columns = append(columns, table.PersonalInfo.Name.SET(postgres.String(user.FirstName)))
	}
	if user.Surname != "" {
		columns = append(columns, table.PersonalInfo.Surname.SET(postgres.String(user.Surname)))
	}
	if user.MiddleName != "" {
		columns = append(columns, table.PersonalInfo.MiddleName.SET(postgres.String(user.MiddleName)))
	}
	if user.Email != "" {
		columns = append(columns, table.PersonalInfo.Email.SET(postgres.String(user.Email)))
	}
	if user.Phone != "" {
		columns = append(columns, table.PersonalInfo.Phone.SET(postgres.String(user.Phone)))
	}

	if len(columns) < 1 {
		r.log.Debug("nothing to updated")
		return nil
	}

	query, args := updatePersonalInfo.SET(columns[0], columns[1:]...).
		WHERE(table.PersonalInfo.PersonalInfoID.EQ(postgres.Int(id))).Sql()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		r.log.Debug("failed to update user", zap.Error(err))
		return err
	}

	r.log.Debug("user updated successfully")
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	r.log.Debug("deleting user")

	query, args := table.PersonalInfo.DELETE().
		WHERE(table.PersonalInfo.PersonalInfoID.EQ(postgres.Int(id))).
		Sql()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		r.log.Debug("failed to delete user", zap.Error(err))
		return err
	}

	r.log.Debug("user deleted successfully")
	return nil
}

func (r *Repository) GetCourse(ctx context.Context, courseID int64) (*models.Course, error) {
	r.log.Debug("getting course")

	var (
		csr     models.Course
		courses = table.Courses
		events  []*models.Event

		rowsProcessed int
	)

	query, args := postgres.SELECT(
		courses.AllColumns,
		table.Events.AllColumns,
	).
		FROM(courses.INNER_JOIN(table.Events, table.Events.CourseID.EQ(courses.CourseID))).
		WHERE(courses.CourseID.EQ(postgres.Int(courseID))).
		Sql()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.log.Debug("failed to get course", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rowsProcessed++
		event := models.Event{}
		var (
			startDate  time.Time
			periodType string
		)
		if err := rows.Scan(
			&csr.ID, &csr.Name, &csr.MonthlySubscriptionCost,
			&event.ID, &event.Description, &startDate,
			&event.RecurrentCount, &event.PeriodFreq, &periodType, &event.CourseID,
		); err != nil {
			return nil, err
		}
		startDateTime := models.MyTime(startDate)
		event.StartDate = &startDateTime
		periodTypeModel, err := models.PeriodTypeString(periodType)
		event.PeriodType = &periodTypeModel
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if rowsProcessed == 0 {
		return nil, repository.ErrCourseNotFound
	}

	csr.Events = events

	return &csr, nil
}

func (r *Repository) SaveCourse(ctx context.Context, crs *models.Course) (int64, error) {
	r.log.Debug("creating course")

	var (
		courseID int64
		courses  = table.Courses
	)

	if len(crs.Events) < 1 {
		return 0, repository.ErrEventsRequired
	}

	if err := withTx(ctx, r.db, func(tx pgx.Tx) error {
		query, args := courses.INSERT(courses.AllColumns.Except(courses.CourseID)).
			VALUES(crs.Name, crs.MonthlySubscriptionCost).RETURNING(courses.CourseID).Sql()

		if err := tx.QueryRow(ctx, query, args...).Scan(&courseID); err != nil {
			r.log.Error("failed to create course", zap.Error(err))
			return err
		}

		if err := r.courseEventsCreate(ctx, tx, courseID, crs.Events); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}

	r.log.Debug("course created successfully", zap.Int64("course_id", courseID))
	return courseID, nil
}

func (r *Repository) SaveEvents(ctx context.Context, courseID int64, events []*models.Event) error {
	r.log.Debug("creating course events")

	if err := withTx(ctx, r.db, func(tx pgx.Tx) error {
		if err := r.courseEventsCreate(ctx, tx, courseID, events); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	r.log.Debug("events course created successfully", zap.Int64("course_id", courseID))
	return nil
}

func withTx(ctx context.Context, db *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *Repository) courseEventsCreate(
	ctx context.Context,
	tx pgx.Tx,
	courseID int64,
	events []*models.Event,
) error {
	for _, e := range events {
		if e != nil {
			query, args := table.Events.
				INSERT(table.Events.AllColumns.Except(table.Events.EventID)).
				VALUES(e.Description, time.Time(*e.StartDate), e.RecurrentCount, e.PeriodFreq, e.PeriodType, courseID).
				RETURNING(table.Events.EventID).Sql()

			if _, err := tx.Exec(ctx, query, args...); err != nil {
				r.log.Error("failed to create course event", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// todo add unique constraint to courses table

func (r *Repository) UpdateCourse(ctx context.Context, id int64, crs *models.Course) error {
	r.log.Debug("updating course")

	updateCourse := table.Courses.UPDATE()
	var columns []any
	if crs.Name != "" {
		columns = append(columns, table.Courses.CourseName.SET(postgres.String(crs.Name)))
	}
	if crs.MonthlySubscriptionCost != nil {
		columns = append(columns, table.Courses.MonthlySubscriptionCost.SET(postgres.Float(*crs.MonthlySubscriptionCost)))
	}

	if len(crs.Events) > 0 {
		if err := r.courseEventsUpdate(ctx, crs.Events); err != nil {
			return err
		}
	}

	if len(columns) < 1 {
		r.log.Debug("course has nothing to update")
		return nil
	}

	query, args := updateCourse.SET(columns[0], columns[1:]...).
		WHERE(table.Courses.CourseID.EQ(postgres.Int(id))).Sql()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		r.log.Debug("failed to update course", zap.Error(err))
		return err
	}

	r.log.Debug("course updated successfully")
	return nil
}

func (r *Repository) courseEventsUpdate(ctx context.Context, events []*models.Event) error {
	for _, e := range events {
		eventUpdate := table.Events.UPDATE()
		var columns []any
		if e.Description != "" {
			columns = append(columns, table.Events.EventDescription.SET(postgres.String(e.Description)))
		}
		if e.RecurrentCount != nil {
			columns = append(columns, table.Events.RecurrentCount.SET(postgres.Int(*e.RecurrentCount)))
		}
		if e.PeriodType != nil {
			columns = append(columns, table.Events.PeriodType.SET(postgres.String(e.PeriodType.String())))
		}
		if e.PeriodFreq != nil {
			columns = append(columns, table.Events.PeriodFreq.SET(postgres.Int(*e.PeriodFreq)))
		}
		if e.StartDate != nil {
			columns = append(columns, table.Events.StartDate.SET(postgres.TimestampT(time.Time(*e.StartDate))))
		}

		if len(columns) < 1 {
			continue
		}

		query, args := eventUpdate.SET(columns[0], columns[1:]...).
			WHERE(table.Events.EventID.EQ(postgres.Int(e.ID))).Sql()
		if _, err := r.db.Exec(ctx, query, args...); err != nil {
			r.log.Debug("failed to update event", zap.Error(err))
			return err
		}
		r.log.Debug("event updated successfully", zap.Int64("event_id", e.ID))
	}

	return nil
}

func (r *Repository) DeleteCourse(ctx context.Context, id int64) error {
	r.log.Debug("deleting course")

	query, args := table.Courses.DELETE().
		WHERE(table.Courses.CourseID.EQ(postgres.Int(id))).
		Sql()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		r.log.Debug("failed to delete course", zap.Error(err))
		return err
	}

	r.log.Debug("course deleted successfully")
	return nil
}

func (r *Repository) DeleteEvent(ctx context.Context, courseID, eventID int64) error {
	r.log.Debug("deleting course event")

	query, args := table.Events.DELETE().
		WHERE(
			postgres.AND(
				table.Events.CourseID.EQ(postgres.Int(courseID)),
				table.Events.EventID.EQ(postgres.Int(eventID)),
			),
		).
		Sql()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		r.log.Debug("failed to delete course event", zap.Error(err))
		return err
	}

	r.log.Debug("course event deleted successfully")
	return nil
}

func (r *Repository) CheckCountEvents(ctx context.Context, courseID int64) (int, error) {
	r.log.Debug("check count of course events")

	events := table.Events
	var count int

	query, args := events.SELECT(postgres.COUNT(events.EventID)).
		WHERE(events.CourseID.EQ(postgres.Int(courseID))).
		Sql()

	if err := r.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		r.log.Debug("failed to check count of course events", zap.Error(err))
		return 0, err
	}

	r.log.Debug("check count of course events successfully")
	return count, nil
}
