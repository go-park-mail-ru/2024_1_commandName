package domain

type ChatSearchRequest struct {
	Word   string `json:"word"`
	UserID uint   `json:"user_id"`
}

type ChatSearchResponse struct {
	Chats  []Chat `json:"chats"`
	UserID uint   `json:"user_id"`
}
