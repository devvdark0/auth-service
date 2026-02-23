package postgresql

import (
	"context"
	"log/slog"
	"time"

	"github.com/devvdark0/auth-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewPOSTGRESQLRepository(db *pgxpool.Pool, log *slog.Logger) *userRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sql := `INSERT INTO users(email, username, password, created_at) VALUES($1, $2, $3, $4) RETURNING id`

	var id int64

	err := ur.db.QueryRow(ctx, sql, user.Email, user.Username, user.PassHash, user.CreatedAt).Scan(&id)
	if err != nil {
		ur.log.Error("failed to save user", "err", err)
		return 0, err
	}

	ur.log.Debug("user succesfully saved", "id", id)
	return id, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sql := `SELECT id, email, username, password, created_at FROM users WHERE email=$1`
	var user models.User

	err := ur.db.QueryRow(ctx, sql, email).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PassHash,
			&user.CreatedAt,
		)
	if err != nil {
		ur.log.Error("failed to get user from db", "email", email)
		return nil, err
	}

	ur.log.Debug("successfully get user", "email", email)
	return &user, nil
}
