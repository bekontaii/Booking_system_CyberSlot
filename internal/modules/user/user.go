package user

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created"`
	UpdatedAt    time.Time `json:"updated"`
}
