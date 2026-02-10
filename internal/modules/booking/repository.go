package booking

import (
	"errors"
	"sync"
)

type Repository interface {
	Create(booking Booking) (Booking, error)
	GetAll() []Booking
	UpdateStatus(id int, status string) error
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
