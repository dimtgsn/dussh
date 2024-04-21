package pgsql

import (
	"context"
	"dussh/internal/config"
	"dussh/internal/domain/models"
	"dussh/internal/repository"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
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
		credID int
		userID int64
		q1     = `INSERT INTO creds(hashed_password) VALUES ($1) RETURNING creds_id;`
		q2     = `INSERT INTO personal_info (creds_id, name, middle_name, surname, email, roles_id, phone)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING personal_info_id;`
	)

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			r.log.Debug("user saved", zap.Int64("user", userID))
			tx.Commit(ctx)
		}
	}()

	if err := tx.QueryRow(ctx, q1, user.Password).Scan(&credID); err != nil {
		r.log.Error("failed to create user creds", zap.Error(err))
		return 0, err
	}

	if err := tx.QueryRow(
		ctx,
		q2,
		credID,
		user.FirstName,
		user.MiddleName,
		user.Surname,
		user.Email,
		int(user.Role),
		user.Phone,
	).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		r.log.Error("failed to save user", zap.Error(err))
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, repository.ErrUserAlreadyExists
		}
		return 0, err
	}

	return userID, nil
}

func (r *Repository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	r.log.Debug("checking user exists")

	var (
		q = `SELECT personal_info_id FROM personal_info WHERE email = $1;`
	)

	if err := r.db.QueryRow(ctx, q, email).Scan(nil); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.log.Debug("getting user by email")

	var (
		user models.User
		q    = `SELECT personal_info_id,name, middle_name, surname, email, roles_id, phone, c.hashed_password as password
				FROM personal_info p 
				    JOIN creds c ON c.creds_id=p.creds_id 
				WHERE email = $1;`
	)

	if err := pgxscan.Get(ctx, r.db, &user, q, email); err != nil {
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

	var (
		user models.User
		q    = `SELECT personal_info_id,name, middle_name, surname, email, roles_id, phone, c.hashed_password as password
				FROM personal_info p 
				    JOIN creds c ON c.creds_id=p.creds_id 
				WHERE personal_info_id = $1;`
	)

	if err := pgxscan.Get(ctx, r.db, &user, q, id); err != nil {
		r.log.Debug("failed to get user by id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
