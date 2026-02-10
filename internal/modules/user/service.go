package user

import (
	"errors"
	"strings"
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

func (s *Service) GetUserByUsername(username string) (User, error) {
	if username == "" {
		return User{}, ErrInvalidUserInput
	}

	return s.repo.GetByUsername(username)
}

func (s *Service) UpdateUser(id int, input User) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUserInput
	}
	if err := validateUser(input); err != nil {
		return User{}, err
	}
	if input.Role != "" {
		input.Role = strings.ToUpper(input.Role)
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

func (s *Service) UpdateUserPartial(id int, input UpdateUserRequest) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUserInput
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return User{}, err
	}

	updated := existing
	if input.Name != nil {
		updated.Name = *input.Name
	}
	if input.Surname != nil {
		updated.Surname = *input.Surname
	}
	if input.Username != nil {
		updated.Username = *input.Username
	}
	if input.Email != nil {
		updated.Email = *input.Email
	}
	if input.Role != nil {
		updated.Role = strings.ToUpper(*input.Role)
	}

	if err := validateUserPartial(updated, input); err != nil {
		return User{}, err
	}

	return s.repo.Update(id, updated)
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

func validateUserPartial(user User, input UpdateUserRequest) error {
	if input.Name != nil && user.Name == "" {
		return ErrInvalidUserInput
	}
	if input.Email != nil && user.Email == "" {
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
