package user

import (
	"errors"
	"time"
)

var (
	ErrInvalidUserInput = errors.New("invalid user input")
	ErrInvalidUserRole  = errors.New("invalid user role")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(input User) (User, error) {
	if err := validateUser(input); err != nil {
		return User{}, err
	}

	if input.Role == "" {
		input.Role = RoleUser
	}
	if input.CreatedAt.IsZero() {
		input.CreatedAt = time.Now().UTC()
	}

	return s.repo.Create(input)
}

func (s *Service) GetUsers() ([]User, error) {
	return s.repo.GetAll()
}

func (s *Service) GetUserByID(id int) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUserInput
	}

	return s.repo.GetByID(id)
}

func (s *Service) UpdateUser(id int, input User) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUserInput
	}
	if err := validateUser(input); err != nil {
		return User{}, err
	}
	if input.Role == "" {
		input.Role = RoleUser
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return User{}, err
	}
	if input.CreatedAt.IsZero() {
		input.CreatedAt = existing.CreatedAt
	}

	return s.repo.Update(id, input)
}

func (s *Service) DeleteUser(id int) error {
	if id <= 0 {
		return ErrInvalidUserInput
	}

	return s.repo.Delete(id)
}

func validateUser(user User) error {
	if user.Name == "" || user.Email == "" {
		return ErrInvalidUserInput
	}

	if user.Role != "" && !isValidRole(user.Role) {
		return ErrInvalidUserRole
	}

	return nil
}

func isValidRole(role string) bool {
	switch role {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}
