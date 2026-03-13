package auth

import (
	"fmt"
	"time"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewJWTManager(cfg *config.AuthConfig) *JWTManager {
	return &JWTManager{
		secret:   cfg.Secret,
		tokenTTL: cfg.TokenTTL,
	}
}

func (j *JWTManager) GenerateToken(id int64, email string) (string, error) {
	exp := time.Now().Add(j.tokenTTL)

	claims := jwt.MapClaims{
		"sub":   string(id),
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *JWTManager) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("invalid token")
		}

		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
