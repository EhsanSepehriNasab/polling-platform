package models

import "time"

// User defines the structure for a user
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSignUp struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
