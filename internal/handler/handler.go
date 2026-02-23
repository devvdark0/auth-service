package handler

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userId); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
