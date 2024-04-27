package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"ProjectMessenger/domain"
)

type Feedback struct {
	db *sql.DB
}

func (f *Feedback) CheckNeedReturnQuestion(ctx context.Context, userID uint, question_id int) (needReturn bool) {
	// 1 - global
	// 2 - <какая то фича>
	counter := 0
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "SELECT COUNT(*) FROM feedback.survey_answers WHERE user_id = $1 and question_id = $2"
	err := f.db.QueryRowContext(ctx, query, userID, question_id).Scan(&counter)
	logger.Debug("CheckNeedReturnQuestion", "userID", userID)
	if counter == 0 {
		return true
	}

	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CheckNeedCSAT, feedback.go",
		}
		logger.Error(customErr.Error())
		return false
	}
	logger.Debug("CheckNeedReturnQuestion: success", "needReturn", false)
	return false
}

func (f *Feedback) GetQuestions(ctx context.Context, userID uint) []domain.Question {
	questions := make([]domain.Question, 0)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "SELECT sq.id, sq.questiontype, sq.question_text FROM feedback.survey_questions sq LEFT JOIN feedback.survey_answers sa ON sq.id = sa.question_id AND sa.user_id = $1 WHERE sa.user_id IS NULL;"
	rows, err := f.db.QueryContext(ctx, query, userID)
	logger.Debug("CheckNeedReturnQuestion", "userID", userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CheckNeedCSAT, feedback.go",
		}
		logger.Error(customErr.Error())
		return questions
	}

	for rows.Next() {
		question := domain.Question{}
		err = rows.Scan(&question.Id, &question.QuestionType, &question.QuestionText)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		questions = append(questions, question)
	}

	logger.Debug("CheckNeedReturnQuestion: success", "needReturn", false)
	return questions
}

func NewFeedbackStorage(db *sql.DB) *Feedback {
	return &Feedback{db: db}
}
