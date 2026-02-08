package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/devvdark0/auth-service/internal/dto"
	"github.com/devvdark0/auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	log         *slog.Logger
}

func NewAuthHandler(auth service.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		log:         log,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ah.log.Info(
		"Handling register request",
		"method", r.Method,
		"path", r.URL.Path,
		"ip", r.RemoteAddr,
	)
	ctx := r.Context()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//TODO: validation

	userId, err := ah.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	resp := &dto.RegisterResponse{
		ID: userId,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	ah.log.Info(
		"Request completed",
		"status", http.StatusCreated,
		"userId", userId,
	)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ah.log.Info(
		"Handling Login request",
		"method", r.Method,
		"path", r.URL.Path,
		"ip", r.RemoteAddr,
	)

	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//TODO: validation

	token, err := ah.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Inernal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.LoginResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
