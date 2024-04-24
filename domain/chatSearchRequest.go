package domain

type ChatSearchRequest struct {
	Word   string `json:"word"`
	UserID uint   `json:"user_id"`
}
