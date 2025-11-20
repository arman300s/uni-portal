package contracts

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role"`
}

type UpdateUserInput struct {
	Email    string `json:"email"`
	RoleName string `json:"role"`
}
