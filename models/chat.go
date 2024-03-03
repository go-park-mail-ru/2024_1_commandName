package models

type Chat struct {
	ID          int        `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AvatarPath  string     `json:"avatar"`
	CreatorID   string     `json:"creator"`
	Messeges    *[]Message `json:"messeges"`
	Users       []ChatUser `json:"useres"`
}
