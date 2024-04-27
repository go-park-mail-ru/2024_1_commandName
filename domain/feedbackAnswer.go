package domain

import "time"

type FeedbackAnswer struct {
	UserID     uint      `json:"user_id"`
	QuestionID int       `json:"question_id"`
	Grade      int       `json:"grade"`
	AnsweredAt time.Time `json:"answered_at"`
}
