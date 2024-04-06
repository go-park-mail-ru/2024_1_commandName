package domain

import "time"

type Chat struct {
	ID                 uint        `json:"id"`
	Type               string      `json:"type"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	AvatarPath         string      `json:"avatar"`
	CreatorID          uint        `json:"creator"`
	Messages           []*Message  `json:"messages"`
	Users              []*ChatUser `json:"users"`
	LastActionDateTime time.Time   `json:"last_action_date_time"`
	LastMessage        Message     `json:"last_message"`
}
