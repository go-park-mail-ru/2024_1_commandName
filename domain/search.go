package domain

type ChatSearchRequest struct {
	Word   string `json:"word"`
	UserID uint   `json:"user_id"`
}

type ChatSearchResponse struct {
	Chats []Chat `json:"chats"`
}

type MessagesSearchRequest struct {
	Word   string `json:"word"`
	UserID uint   `json:"user_id"`
}

type MessagesSearchResponse struct {
	Messages []Message `json:"messages"`
}

type ContactsSearchRequest struct {
	Word   string `json:"word"`
	UserID uint   `json:"user_id"`
}

type ContactsSearchResponse struct {
	Contacts []Person `json:"contacts"`
}
