package models

import (
	"time"
)

type Person struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	About        string    `json:"about"`
	PasswordHash string    `json:"password_hash"`
	CreateDate   time.Time `json:"create_date"`
	LastSeenDate time.Time `json:"last_seen_date"`
	Avatar       string    `json:"avatar"`
	Password     string    `json:"password"`
}
