package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/devvdark0/auth-service/internal/auth"
	"github.com/devvdark0/auth-service/internal/models"
)

type AuthRepository interface {
	Create(ctx context.Context, u *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type authService struct {
	jwt  *auth.JWTManager
	repo AuthRepository
}

func NewAuthService(repo AuthRepository, jwt *auth.JWTManager) *authService {
	return &authService{
		jwt:  jwt,
		repo: repo,
	}
}

func (a *authService) Register(ctx context.Context, username, email, password string) (int64, error) {

	_, err := a.repo.GetByEmail(ctx, email)
	if err == nil {
		return 0, fmt.Errorf("email already in use")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	passHash, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("hashing password err: %w", err)
	}

	user := models.User{
		Username:  username,
		Email:     email,
		PassHash:  passHash,
		CreatedAt: time.Now(),
	}

	userId, err := a.repo.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userId, nil
}

func (a *authService) Login(ctx context.Context, email, password string) (string, error) {

	u, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := auth.VerifyPasssword(password, u.PassHash); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := a.jwt.GenerateToken(u.ID, u.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
