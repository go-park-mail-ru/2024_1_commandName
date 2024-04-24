package domain

type ChatSearchResponse struct {
	Chats  []Chat `json:"chats"`
	UserID uint   `json:"user_id"`
}
