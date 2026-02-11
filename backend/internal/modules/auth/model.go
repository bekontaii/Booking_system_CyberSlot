package auth

type User struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}
