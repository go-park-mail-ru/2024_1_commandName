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
	StickerPath    string         `json:"sticker_path"`
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
	FileID         uint   `json:"file_id"`
}

type StickerFromUser struct {
	StickerID uint `json:"sticker_id"`
}

type FileWithInfo struct {
	FileInfo os.FileInfo
	File     *os.File
}

type Sticker struct {
	StickerID   int    `json:"sticker_id"`
	StickerDesc string `json:"sticker_desc"`
	StickerType string `json:"sticker_type"`
	StickerPath string `json:"sticker_path"`
}

type SummarizeMessageRequest struct {
	Role     string `json:"role,omitempty"`
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
}

type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   string  `json:"maxTokens"`
}

type APIRequest struct {
	ModelURI          string                    `json:"modelUri"`
	CompletionOptions CompletionOptions         `json:"completionOptions"`
	Messages          []SummarizeMessageRequest `json:"messages"`
}
