package domain

import "time"

type Message struct {
	ID              uint      `json:"id" swaggerignore:"true"`
	ChatID          uint      `json:"chat_id"`
	UserID          uint      `json:"user_id" swaggerignore:"true"`
	Message         string    `json:"message_text"`
	Edited          bool      `json:"edited" swaggerignore:"true"`
	CreateTimestamp time.Time `json:"sent_at" swaggerignore:"true"`
	SenderUsername  string    `json:"username"`
}
