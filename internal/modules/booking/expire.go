package booking

import "time"

const DefaultExpiration = 30 * time.Second

// Auto-expire pending bookings after the configured duration.
func (s *Service) startAutoExpire(bookingID int) {
	expireAfter := s.expireAfter
	if expireAfter <= 0 {
		expireAfter = DefaultExpiration
	}

	go func() {
		time.Sleep(expireAfter)
		s.expireIfPending(bookingID)
	}()
}

func (s *Service) expireIfPending(bookingID int) {
	booking, ok := s.getByID(bookingID)
	if !ok {
		return
	}

	if booking.Status != StatusPending {
		return
	}

	_ = s.repo.UpdateStatus(bookingID, StatusExpired)
}
