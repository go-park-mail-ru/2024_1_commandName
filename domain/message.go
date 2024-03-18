package domain

type Message struct {
	ID      int    `json:"id"`
	ChatID  int    `json:"chat_id"`
	UserID  uint   `json:"user_id"`
	Message string `json:"message_text"`
	Edited  bool   `json:"edited"`
}