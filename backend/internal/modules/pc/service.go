package pc

import "errors"

var ErrInvalidPCInput = errors.New("invalid pc input")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePC(input PC) (PC, error) {
	if err := validatePC(input); err != nil {
		return PC{}, err
	}

	if input.Status == "" {
		input.Status = StatusActive
	}

	return s.repo.Create(input)
}

func (s *Service) GetAllPCs() ([]PC, error) {
	return s.repo.GetAll()
}

func (s *Service) GetPCsByClub(clubID int) ([]PC, error) {
	if clubID <= 0 {
		return nil, ErrInvalidPCInput
	}

	return s.repo.GetByClubID(clubID)
}

func (s *Service) UpdatePC(id int, input PC) (PC, error) {
	if id <= 0 {
		return PC{}, ErrInvalidPCInput
	}
	if err := validatePC(input); err != nil {
		return PC{}, err
	}
	if input.Status == "" {
		input.Status = StatusActive
	}

	return s.repo.Update(id, input)
}

func (s *Service) DeletePC(id int) error {
	if id <= 0 {
		return ErrInvalidPCInput
	}

	return s.repo.Delete(id)
}

func validatePC(pc PC) error {
	if pc.ClubID <= 0 || pc.Name == "" || pc.RAM <= 0 {
		return ErrInvalidPCInput
	}

	if pc.Status != "" && !isValidStatus(pc.Status) {
		return ErrInvalidPCInput
	}

	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case StatusActive, StatusBroken, StatusMaintenance:
		return true
	default:
		return false
	}
}
