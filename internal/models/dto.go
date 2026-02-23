package models

type RegisterRequest struct {
	Username string `json:"username" validate:"required, min=5, max=100"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=8"`
}

type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
