package auth

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	mu             sync.Mutex
	users          map[string]User
	jwtSecret      string
	jwtExpireHours int
}

func NewService(secret string, expireHours int) *Service {
	if expireHours <= 0 {
		expireHours = 24
	}

	return &Service{
		users:          make(map[string]User),
		jwtSecret:      secret,
		jwtExpireHours: expireHours,
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

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[input.Username]; exists {
		return errors.New("username already exists")
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return err
	}

	s.users[input.Username] = User{
		Name:         input.Name,
		Surname:      input.Surname,
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: hash,
	}

	return nil
}

// -------------------- LOGIN --------------------

func (s *Service) Login(ctx context.Context, input LoginRequest) (string, error) {
	if input.Username == "" || input.Password == "" {
		return "", errors.New("username or password is empty")
	}

	s.mu.Lock()
	user, ok := s.users[input.Username]
	s.mu.Unlock()
	if !ok {
		return "", errors.New("invalid credentials")
	}

	if err := CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateToken(user.Username, s.jwtExpireHours, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
