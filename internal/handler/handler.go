package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/devvdark0/auth-service/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
}

func NewAuthHandler(log *slog.Logger, authService AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	a.log.Info("start processing register request")

	w.Header().Set("Contetn-Type", "application/json")
	ctx := r.Context()

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.Error("failed to parse body", "err", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userId, err := a.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		a.log.Error("failed to register user", "err", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	resp := models.RegisterResponse{
		UserID: userId,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.log.Error("failed to encode response", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	a.log.Info("user was sucessfully registered", "user_id", userId)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	a.log.Info("start processing login request")

	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.Error("failed to decode request", "err", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	token, err := a.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		a.log.Error("failde to login user", "err", err)
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.log.Error("failed to encode response", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	a.log.Info("user was successfully logged in", "token", token)
}
