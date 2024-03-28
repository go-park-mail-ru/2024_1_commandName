package domain

import (
	"time"
)

type Person struct {
	ID           uint      `json:"id" swaggerignore:"true"`
	Username     string    `json:"username"`
	Email        string    `json:"email" swaggerignore:"true"`
	Name         string    `json:"name" swaggerignore:"true"`
	Surname      string    `json:"surname" swaggerignore:"true"`
	About        string    `json:"about" swaggerignore:"true"`
	Password     string    `json:"password,omitempty"`
	CreateDate   time.Time `json:"create_date" swaggerignore:"true"`
	LastSeenDate time.Time `json:"last_seen_date" swaggerignore:"true"`
	Avatar       string    `json:"avatar" swaggerignore:"true"`
	PasswordSalt string    `json:"password_salt,omitempty" swaggerignore:"true"`
}
