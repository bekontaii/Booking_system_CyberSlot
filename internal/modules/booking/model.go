package booking

import "time"

type Booking struct {
	ID        int
	PCID      int
	UserID    int
	StartTime time.Time
	EndTime   time.Time
	Status    string
}
