package domain

type TranslateRequest struct {
	Messages           []string `json:"texts"`
	FolderID           string   `json:"folderId"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
}
