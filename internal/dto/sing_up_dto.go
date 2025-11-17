package dto

type SignUpDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Guest    bool   `json:"guest"`
}
