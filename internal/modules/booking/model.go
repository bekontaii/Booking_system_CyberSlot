package booking

import "time"

type Booking struct {
	ID        int       `json:"id"`
	PCID      int       `json:"pc_id"`
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}
