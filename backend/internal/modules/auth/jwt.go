package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
)

func GenerateToken(u user.User, expireHours int, secret string) (string, error) {
	if u.Username == "" {
		return "", errors.New("username cannot be empty")
	}
	if u.ID <= 0 {
		return "", errors.New("user id must be positive")
	}
	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}

	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Duration(expireHours) * time.Hour)
	claims := jwt.MapClaims{
		"sub":     u.Username,
		"user_id": u.ID,
		"role":    u.Role,
		"club_id": u.ClubID,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
