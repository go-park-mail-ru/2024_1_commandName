package domain

type Chat struct {
	ID          uint        `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	AvatarPath  string      `json:"avatar"`
	CreatorID   string      `json:"creator"`
	Messages    []*Message  `json:"messages,omitempty"`
	Users       []*ChatUser `json:"users"`
	LastMessage Message     `json:"last_message"`
}
