package booking

import (
	"errors"
	"sync"
)

type Repository interface {
	Create(booking Booking) (Booking, error)
	GetAll() []Booking
	GetByID(id int) (Booking, bool)
	GetByClubID(clubID int) []Booking
	UpdateStatus(id int, status string) error
	Update(booking Booking) error
	Delete(id int) error
}

type InMemoryRepository struct {
	mu       sync.Mutex
	bookings []Booking
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{bookings: make([]Booking, 0)}
}

func (r *InMemoryRepository) Create(booking Booking) (Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	maxID := 0
	for _, existing := range r.bookings {
		if existing.ID > maxID {
			maxID = existing.ID
		}
	}
	booking.ID = maxID + 1

	r.bookings = append(r.bookings, booking)
	return booking, nil
}

func (r *InMemoryRepository) GetAll() []Booking {
	r.mu.Lock()
	defer r.mu.Unlock()

	copyBookings := make([]Booking, len(r.bookings))
	copy(copyBookings, r.bookings)
	return copyBookings
}

func (r *InMemoryRepository) GetByID(id int) (Booking, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, b := range r.bookings {
		if b.ID == id {
			return b, true
		}
	}
	return Booking{}, false
}

func (r *InMemoryRepository) GetByClubID(clubID int) []Booking {
	r.mu.Lock()
	defer r.mu.Unlock()

	// In-memory booking repository has no relation to clubs.
	_ = clubID
	out := make([]Booking, len(r.bookings))
	copy(out, r.bookings)
	return out
}

func (r *InMemoryRepository) UpdateStatus(id int, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, booking := range r.bookings {
		if booking.ID == id {
			r.bookings[i].Status = status
			return nil
		}
	}

	return errors.New("booking not found")
}

func (r *InMemoryRepository) Update(booking Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.bookings {
		if r.bookings[i].ID == booking.ID {
			r.bookings[i] = booking
			return nil
		}
	}

	return errors.New("booking not found")
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, booking := range r.bookings {
		if booking.ID == id {
			r.bookings = append(r.bookings[:i], r.bookings[i+1:]...)
			return nil
		}
	}

	return errors.New("booking not found")
}
