package auth

import (
	"strconv"
	"time"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:   cfg.Secret,
		tokenTTL: cfg.TokenTTL,
	}
}

func (j *JWTManager) GenerateToken(id int64, email string) (string, error) {
	exp := time.Now().Add(j.tokenTTL).Unix()

	claims := jwt.MapClaims{
		"sub":   strconv.Itoa(int(id)),
		"email": email,
		"exp":   exp,
		"iat":   time.Now(),
	}

	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenStr.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
