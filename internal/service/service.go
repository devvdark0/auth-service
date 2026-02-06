package service

import (
	"context"
	"fmt"

	"github.com/devvdark0/auth-service/internal/auth"
	"github.com/devvdark0/auth-service/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repo: repository,
	}
}

func (us *userService) Register(ctx context.Context, email, password string) (int64, error) {
	if _, err := us.repo.GetUserByEmail(ctx, email); err == nil {
		return 0, fmt.Errorf("user is already exists with such email %w", err)
	}

	hashPass, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password %w", err)
	}

	userID, err := us.repo.CreateUser(ctx, email, hashPass)
	if err != nil {
		return 0, fmt.Errorf("failed to create user %w", err)
	}

	return userID, nil
}

func (us *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("user not found %w", err)
	}

	if err := auth.VerifyPassword(password, user.Password); err != nil {
		return "", fmt.Errorf("incorect password %w", err)
	}

	//TODO: return token
	return "", nil
}
