package dto

type RegisterRequest struct {
	Email    string `json:"email" validator:"required, email"`
	Password string `json:"password" validator:"required, min=8"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email" validator:"required, email"`
	Password string `json:"password" validator:"required, min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
