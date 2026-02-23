package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/devvdark0/auth-service/internal/models"
	"github.com/devvdark0/auth-service/internal/service"
)

type AuthHandler struct {
	log      *slog.Logger
	authServ service.AuthService
}

func NewAuthHandler(log *slog.Logger, service service.AuthService) *AuthHandler {
	return &AuthHandler{
		log:      log,
		authServ: service,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	a.log.Info("start register request...")
	ctx := r.Context()

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userId, err := a.authServ.Register(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	resp := models.RegisterResponse{
		UserID: userId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	a.log.Info("request completed successfully!")
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	a.log.Info("start processing login request...")

	ctx := r.Context()

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := a.authServ.Login(ctx, &req)
	if err != nil {
		http.Error(w, "failed to log in", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	a.log.Info("request completed successfully!")
}
