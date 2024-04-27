package domain

type Question struct {
	Id           int    `json:"question_id"`
	QuestionText string `json:"question_text"`
	QuestionType string `json:"question_type"`
}
