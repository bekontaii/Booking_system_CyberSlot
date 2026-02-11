package pc

const (
	StatusActive      = "ACTIVE"
	StatusBroken      = "BROKEN"
	StatusMaintenance = "MAINTENANCE"
)

type PC struct {
	ID     int    `json:"id"`
	ClubID int    `json:"club_id"`
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	GPU    string `json:"gpu"`
	RAM    int    `json:"ram"`
	Status string `json:"status"`
}
