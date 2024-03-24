package domain

import "time"

type Chat struct {
	ID                int         `json:"id"`
	Type              string      `json:"type"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	AvatarPath        string      `json:"avatar"`
	CreatorID         string      `json:"creator"`
	Messages          []*Message  `json:"messages"`
	Users             []*ChatUser `json:"users"`
	LastMessageSentAt time.Time   `json:"sent_at"`
}
