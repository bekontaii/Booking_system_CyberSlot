package club

import "errors"

var ErrInvalidClubInput = errors.New("invalid club input")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateClub(input Club) (Club, error) {
	if err := validateClub(input); err != nil {
		return Club{}, err
	}

	return s.repo.Create(input)
}

func (s *Service) GetClubs() ([]Club, error) {
	return s.repo.GetAll()
}

func (s *Service) GetClubByID(id int) (Club, error) {
	if id <= 0 {
		return Club{}, ErrInvalidClubInput
	}

	return s.repo.GetByID(id)
}

func (s *Service) UpdateClub(id int, input Club) (Club, error) {
	if id <= 0 {
		return Club{}, ErrInvalidClubInput
	}
	if err := validateClub(input); err != nil {
		return Club{}, err
	}

	return s.repo.Update(id, input)
}

func (s *Service) DeleteClub(id int) error {
	if id <= 0 {
		return ErrInvalidClubInput
	}

	return s.repo.Delete(id)
}

func validateClub(club Club) error {
	if club.Name == "" || club.City == "" || club.Address == "" {
		return ErrInvalidClubInput
	}

	return nil
}
