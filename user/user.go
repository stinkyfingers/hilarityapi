package user

import "time"

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Auth struct {
	User      User      `json:"user"`
	OIDCToken string    `json:"token"`
	Expires   time.Time `json:"expires"`
}
