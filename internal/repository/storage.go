package repository

import "github.com/devvdark0/auth-service/internal/model"

type Storage interface {
	Create(model.User) error
}