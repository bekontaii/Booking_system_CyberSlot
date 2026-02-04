package booking

import "time"

func Conflicts(start, end time.Time, existing Booking) bool {
	return start.Before(existing.EndTime) && end.After(existing.StartTime)
}
