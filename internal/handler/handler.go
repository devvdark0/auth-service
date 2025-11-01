package handler

import "github.com/devvdark0/auth-service/internal/service"

type AuthHandler struct {
	authService service.AuthService
}

func InitAuthHandler(authService service.AuthService) AuthHandler {
	return AuthHandler{authService: authService}
}
