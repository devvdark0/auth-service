package service

import "github.com/devvdark0/auth-service/internal/repository"

type AuthService struct {
	repository repository.Storage
}
