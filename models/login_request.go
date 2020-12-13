package models

// LoginRequest is the user creation request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
