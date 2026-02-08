package booking

import "time"

type CreateBookingRequest struct {
	PCID      int       `json:"pc_id"`
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
