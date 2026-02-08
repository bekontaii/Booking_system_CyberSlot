package booking

import (
	"errors"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrBookingConflict  = errors.New("booking conflict")
)

type Service struct {
	repo        Repository
	expireAfter time.Duration
}

func NewService(repo Repository, expireAfter time.Duration) *Service {
	if expireAfter <= 0 {
		expireAfter = DefaultExpiration
	}

	return &Service{repo: repo, expireAfter: expireAfter}
}

func (s *Service) CreateBooking(request CreateBookingRequest) (Booking, error) {
	if err := validateTimeRange(request.StartTime, request.EndTime); err != nil {
		return Booking{}, err
	}

	existing := s.repo.GetAll()
	if !IsAvailable(existing, request.PCID, request.StartTime, request.EndTime) {
		return Booking{}, ErrBookingConflict
	}

	booking := Booking{
		ID:        nextID(existing),
		PCID:      request.PCID,
		UserID:    request.UserID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		Status:    StatusPending,
	}

	if err := s.repo.Create(booking); err != nil {
		return Booking{}, err
	}

	s.startAutoExpire(booking.ID)
	return booking, nil
}

func (s *Service) GetAllBookings() []Booking {
	return s.repo.GetAll()
}

func validateTimeRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() || !start.Before(end) {
		return ErrInvalidTimeRange
	}

	return nil
}

func nextID(bookings []Booking) int {
	maxID := 0
	for _, booking := range bookings {
		if booking.ID > maxID {
			maxID = booking.ID
		}
	}

	return maxID + 1
}

func (s *Service) getByID(id int) (Booking, bool) {
	bookings := s.repo.GetAll()
	for _, booking := range bookings {
		if booking.ID == id {
			return booking, true
		}
	}

	return Booking{}, false
}
