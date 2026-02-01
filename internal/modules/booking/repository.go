package booking

import "sync"

var (
	bookings = make([]Booking, 0)
	mu       sync.Mutex
)

func Save(b Booking) {
	mu.Lock()
	defer mu.Unlock()
	bookings = append(bookings, b)
}

func GetAll() []Booking {
	return bookings
}
