package user

import "time"

const (
	RoleUser      = "USER"
	RoleClubAdmin = "CLUB_ADMIN"
	RoleSiteAdmin = "SITE_ADMIN"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	ClubID       *int      `json:"club_id,omitempty"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
}
