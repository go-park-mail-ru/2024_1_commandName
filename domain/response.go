package domain

// Response[T]
type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Body   T   `json:"body,omitempty"`
}

type Error struct {
	Error string `json:"error" example:"error description"`
}

type Chats struct {
	Chats []Chat `json:"chats"`
}
