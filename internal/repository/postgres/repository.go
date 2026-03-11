package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/devvdark0/auth-service/internal/models"
)

type authRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewAuthRepository(db *sql.DB, log *slog.Logger) *authRepository {
	return &authRepository{
		db:  db,
		log: log,
	}
}

func (a *authRepository) Create(ctx context.Context, u *models.User) (int64, error) {
	a.log.Debug("creating user", "email", u.Email)

	sql := `INSERT INTO users(username, email, password, created_at) VALUES(?, ?, ?, ?)`

	res, err := a.db.ExecContext(ctx, sql, u.Username, u.Email, u.PassHash, u.CreatedAt)
	if err != nil {
		a.log.Error("failed to create user with such credentials", "err", err)
		return 0, err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		a.log.Error("failed to get id of user", "err", err)
		return 0, err
	}

	return userId, nil
}

func (a *authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	a.log.Debug("get user from db", "email", email)

	sql := `SELECT * FROM users WHERE email=?`

	var user models.User

	err := a.db.QueryRowContext(ctx, sql, email).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
		)
	if err != nil {
		a.log.Error("failed to get user", "err", err)
		return nil, err
	}

	return &user, nil
}
