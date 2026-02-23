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
	manager   *auth.JWTManager
	store     repository.UserRepository
}

func NewAuthService(storage repository.UserRepository, manager *auth.JWTManager) *AuthService {
	validator := validator.New()
	return &AuthService{
		validator: validator,
		manager:   manager,
		store:     storage,
	}
}

func (a *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (int64, error) {
	if err := a.validator.Struct(req); err != nil {
		return 0, fmt.Errorf("invalid data: %w", err)
	}

	_, err := a.store.GetByEmail(ctx, req.Email)
	if err == nil {
		return 0, fmt.Errorf("email is already in use")
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

func (a *AuthService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	if err := a.validator.Struct(req); err != nil {
		return "", fmt.Errorf("invalid data: %w", err)
	}

	user, err := a.store.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("user does not exists")
	}

	token, err := a.manager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
