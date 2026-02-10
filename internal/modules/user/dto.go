package user

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Role     *string `json:"role"`
}
