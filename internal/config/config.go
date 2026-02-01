package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret      string
	JWTExpireHours int
}

func Load() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	expireStr := os.Getenv("JWT_EXPIRE_HOURS")
	expireHours := 24 // default

	if expireStr != "" {
		value, err := strconv.Atoi(expireStr)
		if err != nil {
			return nil, errors.New("JWT_EXPIRE_HOURS must be a number")
		}
		expireHours = value
	}

	return &Config{
		JWTSecret:      jwtSecret,
		JWTExpireHours: expireHours,
	}, nil
}
