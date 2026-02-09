package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/devvdark0/auth-service/internal/auth"
	"github.com/devvdark0/auth-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailInUse         = errors.New("email is already in use")
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	repo      repository.UserRepository
	validator *auth.JWTValidator
}

func NewAuthService(repository repository.UserRepository, validator *auth.JWTValidator) AuthService {
	return &authService{
		repo:      repository,
		validator: validator,
	}
}

func (us *authService) Register(ctx context.Context, email, password string) (int64, error) {
	_, err := us.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return 0, ErrEmailInUse
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	hashPass, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	userID, err := us.repo.CreateUser(ctx, email, hashPass)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (us *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := auth.VerifyPassword(password, user.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := us.validator.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
