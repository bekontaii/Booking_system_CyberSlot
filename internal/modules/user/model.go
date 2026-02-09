package user

import "time"

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User1 struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
