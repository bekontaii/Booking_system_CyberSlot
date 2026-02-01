package auth

import (
	"context"
	"errors"

	"github.com/bekontaii/Booking_system_CyberSlot/internal/modules/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// -------------------- PASSWORD --------------------

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

// -------------------- REGISTER --------------------

func (s *Service) Register(ctx context.Context, input RegisterRequest) error {
	if input.Email == "" || input.Password == "" || input.Username == "" {
		return errors.New("email, username or password is empty")
	}

	// check if user already exists
	_, err := s.repo.FindByUsername(ctx, input.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return err
	}

	u := &user.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
		Name:         input.Name,
		Surname:      input.Surname,
	}

	return s.repo.Create(ctx, u)
}

// -------------------- LOGIN --------------------

func (s *Service) Login(ctx context.Context, input LoginRequest) (*user.User, error) {
	if input.Username == "" || input.Password == "" {
		return nil, errors.New("username or password is empty")
	}

	u, err := s.repo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := CheckPasswordHash(input.Password, u.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}
