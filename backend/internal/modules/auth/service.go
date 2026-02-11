package auth

import (
	"context"
	"errors"
	user2 "github.com/bekontaii/Booking_system_CyberSlot/backend/internal/modules/user"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo           user2.Repository
	jwtSecret      string
	jwtExpireHours int
}

func NewService(repo user2.Repository, secret string, expireHours int) *Service {
	if expireHours <= 0 {
		expireHours = 24
	}

	return &Service{
		repo:           repo,
		jwtSecret:      secret,
		jwtExpireHours: expireHours,
	}
}

func HashPassword(plainPassword string) (string, error) {
	if plainPassword == "" {
		return "", errors.New("password is empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(plainPassword string, hashedPassword string) error {
	if plainPassword == "" || hashedPassword == "" {
		return errors.New("password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func (s *Service) Register(ctx context.Context, input RegisterRequest) error {
	if input.Email == "" || input.Password == "" || input.Username == "" {
		return errors.New("email, username or password is empty")
	}
	if input.Name == "" || input.Surname == "" {
		return errors.New("name or surname is empty")
	}

	if _, err := s.repo.GetByUsername(input.Username); err == nil {
		return errors.New("username already exists")
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return err
	}

	_, err = s.repo.Create(user2.User{
		Name:         input.Name,
		Surname:      input.Surname,
		Email:        input.Email,
		Username:     input.Username,
		Role:         user2.RoleUser,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, input LoginRequest) (string, error) {
	if input.Username == "" || input.Password == "" {
		return "", errors.New("username or password is empty")
	}

	u, err := s.repo.GetByUsername(input.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := CheckPasswordHash(input.Password, u.PasswordHash); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateToken(u.Username, s.jwtExpireHours, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
