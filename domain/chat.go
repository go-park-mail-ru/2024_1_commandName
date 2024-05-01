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
	CreatedAt          time.Time   `json:"created_at"`
	EditedAt           time.Time   `json:"edited_at"`
	LastActionDateTime time.Time   `json:"last_action_date_time"`
	LastMessage        Message     `json:"last_message"`
	LastSeenMessageID  int         `json:"last_seen_message_id"`
}

type ChannelWithCounter struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NumOfUsers  int    `json:"numOfUsers"`
}
