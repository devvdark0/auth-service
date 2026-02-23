package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func ValidatePassword(passHash, originPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(passHash), []byte(originPass))
}
