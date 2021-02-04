package entity

// Credentials user's credentials with login&password.
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
