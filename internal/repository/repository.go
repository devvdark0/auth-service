package repository

import (
	"context"

	"github.com/devvdark0/auth-service/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
