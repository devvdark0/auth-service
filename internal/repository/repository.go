package repository

import (
	"context"

	"github.com/devvdark0/auth-service/internal/models"
)

type AuthRepository interface {
	Create(ctx context.Context, u *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
