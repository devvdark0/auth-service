package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(origPass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(origPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(origPass, hashPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(origPass))
}
