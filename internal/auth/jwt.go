package auth

import (
	"fmt"
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

func (j *JWTManager) VerifyToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to parse token")
		}

		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
