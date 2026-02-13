package booking

import (
	"errors"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrBookingConflict  = errors.New("booking conflict")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrInvalidStatus    = errors.New("invalid booking status")
	ErrClubInactive     = errors.New("Club is currently inactive")
)

type PCClubAvailabilityChecker interface {
	IsPCBookable(pcID int) bool
}

type Service struct {
	repo        Repository
	expireAfter time.Duration
	checker     PCClubAvailabilityChecker
}

func NewService(repo Repository, checker PCClubAvailabilityChecker, expireAfter time.Duration) *Service {
	if expireAfter <= 0 {
		expireAfter = DefaultExpiration
	}

	return &Service{repo: repo, checker: checker, expireAfter: expireAfter}
}

func (s *Service) CreateBooking(request CreateBookingRequest) (Booking, error) {
	if err := validateTimeRange(request.StartTime, request.EndTime); err != nil {
		return Booking{}, err
	}
	if s.checker != nil && !s.checker.IsPCBookable(request.PCID) {
		return Booking{}, ErrClubInactive
	}

	existing := s.repo.GetAll()
	if !IsAvailable(existing, request.PCID, request.StartTime, request.EndTime) {
		return Booking{}, ErrBookingConflict
	}

	booking := Booking{
		PCID:      request.PCID,
		UserID:    request.UserID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		Status:    StatusPending,
	}

	created, err := s.repo.Create(booking)
	if err != nil {
		return Booking{}, err
	}

	s.startAutoExpire(created.ID)
	return created, nil
}

func (s *Service) GetAllBookings() []Booking {
	return s.repo.GetAll()
}

func (s *Service) GetBookingsByClub(clubID int) []Booking {
	if clubID <= 0 {
		return nil
	}
	return s.repo.GetByClubID(clubID)
}

func (s *Service) GetBookingByID(id int) (Booking, bool) {
	if id <= 0 {
		return Booking{}, false
	}
	return s.repo.GetByID(id)
}

func (s *Service) UpdateBooking(id int, request UpdateBookingRequest) (Booking, error) {
	current, ok := s.getByID(id)
	if !ok {
		return Booking{}, ErrBookingNotFound
	}

	updated := current
	if request.PCID != nil {
		updated.PCID = *request.PCID
	}
	if request.UserID != nil {
		updated.UserID = *request.UserID
	}
	if request.StartTime != nil {
		updated.StartTime = *request.StartTime
	}
	if request.EndTime != nil {
		updated.EndTime = *request.EndTime
	}
	if request.Status != nil {
		if !isValidStatus(*request.Status) {
			return Booking{}, ErrInvalidStatus
		}
		updated.Status = *request.Status
	}

	if err := validateTimeRange(updated.StartTime, updated.EndTime); err != nil {
		return Booking{}, err
	}

	if updated.Status == StatusPending || updated.Status == StatusConfirmed {
		existing := s.repo.GetAll()
		if !IsAvailableExcept(existing, updated.ID, updated.PCID, updated.StartTime, updated.EndTime) {
			return Booking{}, ErrBookingConflict
		}
	}

	if err := s.repo.Update(updated); err != nil {
		return Booking{}, err
	}

	return updated, nil
}

func (s *Service) DeleteBooking(id int) error {
	if _, ok := s.getByID(id); !ok {
		return ErrBookingNotFound
	}

	return s.repo.Delete(id)
}

func validateTimeRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() || !start.Before(end) {
		return ErrInvalidTimeRange
	}

	return nil
}

func (s *Service) getByID(id int) (Booking, bool) {
	return s.repo.GetByID(id)
}

func isValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusExpired:
		return true
	default:
		return false
	}
}
