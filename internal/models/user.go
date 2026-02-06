package models

type User struct {
	ID       int64
	Email    string
	Password string
}

type RegisterRequest struct {
	Email    string
	Password string
}

type RegisterResponse struct {
	ID int64
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}
