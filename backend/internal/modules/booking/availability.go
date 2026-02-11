package booking

import "time"

func IsAvailable(bookings []Booking, pcID int, start, end time.Time) bool {
	for _, booking := range bookings {
		if booking.PCID != pcID {
			continue
		}

		if booking.Status != StatusPending && booking.Status != StatusConfirmed {
			continue
		}

		if Conflicts(start, end, booking) {
			return false
		}
	}

	return true
}
