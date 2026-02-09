package user

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	CreatedAt    time.Time `json:"created"`
	UpdatedAt    time.Time `json:"updated"`
}
