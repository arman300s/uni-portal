package contracts

type SignupInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID uint   `json:"id"`
}
