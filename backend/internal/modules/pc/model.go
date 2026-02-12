package pc

import "time"

const (
	StatusActive      = "ACTIVE"
	StatusBroken      = "BROKEN"
	StatusMaintenance = "MAINTENANCE"
)

type PC struct {
	ID       int       `json:"id"`
	ClubID   int       `json:"club_id"`
	PCNumber int       `json:"pc_number"`
	Name     string    `json:"name,omitempty"`
	CPU      string    `json:"cpu,omitempty"`
	GPU      string    `json:"gpu,omitempty"`
	RAM      int       `json:"ram,omitempty"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created_at,omitempty"`
}
