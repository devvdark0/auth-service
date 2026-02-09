package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is no longer valid")
)

type JWTValidator struct {
	secret   string
	tokenTTL time.Duration
}

func NewJWTValidator(key string, ttl time.Duration) *JWTValidator {
	return &JWTValidator{
		secret:   key,
		tokenTTL: ttl,
	}
}

func (j *JWTValidator) GenerateToken(id int64, email string) (string, error) {
	exp := time.Now().Add(j.tokenTTL).Unix()

	claims := jwt.MapClaims{
		"sub":   strconv.Itoa(int(id)),
		"email": email,
		"exp":   exp,
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *JWTValidator) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
