package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(subject string, expireHours int, secret string) (string, error) {
	if subject == "" {
		return "", errors.New("subject cannot be empty")
	}
	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(expireHours) * time.Hour)
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
