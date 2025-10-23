package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}