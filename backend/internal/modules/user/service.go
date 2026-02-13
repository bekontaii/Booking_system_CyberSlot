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
	if input.ClubID != nil {
		updated.ClubID = input.ClubID
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
	if user.ClubID != nil && *user.ClubID <= 0 {
		return ErrInvalidUserInput
	}

	if user.Role != "" && !isValidRole(user.Role) {
		return ErrInvalidUserRole
	}
	if strings.ToUpper(user.Role) == RoleClubAdmin && user.ClubID == nil {
		return ErrInvalidUserInput
	}
	if strings.ToUpper(user.Role) != RoleClubAdmin && user.ClubID != nil {
		return ErrInvalidUserInput
	}

	return nil
}

func (s *Service) UpdateUserRole(id int, role string, clubID *int) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUserInput
	}

	role = strings.ToUpper(strings.TrimSpace(role))
	if !isValidRole(role) {
		return User{}, ErrInvalidUserRole
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return User{}, err
	}

	existing.Role = role
	existing.ClubID = clubID
	if err := validateUserPartial(existing, UpdateUserRequest{Role: &role, ClubID: clubID}); err != nil {
		return User{}, err
	}

	return s.repo.Update(id, existing)
}

func validateUserPartial(user User, input UpdateUserRequest) error {
	if input.Name != nil && user.Name == "" {
		return ErrInvalidUserInput
	}
	if input.Email != nil && user.Email == "" {
		return ErrInvalidUserInput
	}
	if input.ClubID != nil && *input.ClubID <= 0 {
		return ErrInvalidUserInput
	}
	if user.ClubID != nil && *user.ClubID <= 0 {
		return ErrInvalidUserInput
	}
	if user.Role != "" && !isValidRole(user.Role) {
		return ErrInvalidUserRole
	}
	if strings.ToUpper(user.Role) == RoleClubAdmin && user.ClubID == nil {
		return ErrInvalidUserInput
	}
	if strings.ToUpper(user.Role) != RoleClubAdmin && user.ClubID != nil {
		return ErrInvalidUserInput
	}
	return nil
}

func isValidRole(role string) bool {
	switch role {
	case RoleUser, RoleSiteAdmin, RoleClubAdmin:
		return true
	default:
		return false
	}
}
