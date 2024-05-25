package domain

import (
	"os"
	"time"
)

type Message struct {
	ID             uint           `json:"id" swaggerignore:"true"`
	ChatID         uint           `json:"chat_id"`
	UserID         uint           `json:"user_id" swaggerignore:"true"`
	Message        string         `json:"message_text"`
	Edited         bool           `json:"edited" swaggerignore:"true"` //TODO
	EditedAt       time.Time      `json:"edited_at" swaggerignore:"true"`
	CreatedAt      time.Time      `json:"sent_at" swaggerignore:"true"`
	SenderUsername string         `json:"username"`
	File           *FileInMessage `json:"file"`
}

type FileInMessage struct {
	OriginalName string `json:"original_name"`
	Path         string `json:"path"`
	Type         string `json:"type"`
}

type FileFromUser struct {
	MessageText    string `json:"message_text"`
	MessageID      uint
	ChatID         uint   `json:"chat_id"`
	AttachmentType string `json:"type"`
}

type FileWithInfo struct {
	FileInfo os.FileInfo
	File     *os.File
}
