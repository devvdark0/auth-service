package service

import (
	"context"
	"fmt"
	"time"

	"github.com/devvdark0/auth-service/internal/auth"
	"github.com/devvdark0/auth-service/internal/models"
	"github.com/devvdark0/auth-service/internal/repository"
	"github.com/go-playground/validator/v10"
)

type AuthService struct {
	validator *validator.Validate
	store     repository.UserRepository
}

func NewAuthService(storage repository.UserRepository) *AuthService {
	validator := validator.New()
	return &AuthService{
		validator: validator,
		store:     storage,
	}
}

func (a *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (int64, error) {
	if err := a.validator.Struct(req); err != nil {
		return 0, fmt.Errorf("invalid data: %w", err)
	}

	passHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		PassHash:  passHash,
		CreatedAt: time.Now(),
	}

	userId, err := a.store.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	return userId, nil
}
