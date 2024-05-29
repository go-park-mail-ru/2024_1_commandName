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

type User struct {
	User Person `json:"user"`
}

type Contacts struct {
	Contacts []Person `json:"contacts"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

type Alternative struct {
	Message SummarizeMessageRequest `json:"message"`
	Status  string                  `json:"status"`
}

type Usage struct {
	InputTextTokens  string `json:"inputTextTokens"`
	CompletionTokens string `json:"completionTokens"`
	TotalTokens      string `json:"totalTokens"`
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Usage        Usage         `json:"usage"`
	ModelVersion string        `json:"modelVersion"`
}
type SummarizeMessageResponse struct {
	Result Result `json:"result"`
}
