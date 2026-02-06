package postgres

import (
	"context"
	"log/slog"

	"github.com/devvdark0/auth-service/internal/models"
	"github.com/devvdark0/auth-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewPOSTGRESQLRepository(db *pgxpool.Pool, log *slog.Logger) repository.UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, email, password string) (int64, error) {
	sql := `INSERT INTO users(email, password) VALUES($1, $2) RETURNING id`

	var userId int64

	err := ur.db.QueryRow(ctx, sql, email, password).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, email, password FROM users WHERE email=$1`

	var user models.User

	err := ur.db.QueryRow(ctx, sql, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
