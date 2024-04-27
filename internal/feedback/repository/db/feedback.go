package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"ProjectMessenger/domain"
)

type Feedback struct {
	db *sql.DB
}

func (f *Feedback) CheckNeedCSAT(ctx context.Context, userID uint, typeOfQuestion int) (isUserAnswered bool) {
	// 1 - global
	// 2 - фича
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := ""
	if typeOfQuestion == 1 {
		query = "SELECT is_answered FROM feedback.survey_questions WHERE user_id = $1 and question_id = $2"
	} else {
		query = "SELECT is_answered FROM feedback.survey_questions WHERE user_id = $1 and question_id = $2"
	}
	err := f.db.QueryRowContext(ctx, query, userID, typeOfQuestion).Scan(&isUserAnswered)
	logger.Debug("CheckNeedCSAT", "userID", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("didn't found user by userID", "userID", userID)
			return false
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CheckNeedCSAT, feedback.go",
		}
		logger.Error(customErr.Error())
		return false
	}
	logger.Debug("CheckNeedCSAT: success", "isUserAnswered", isUserAnswered)
	return isUserAnswered
}

func NewFeedbackStorage(db *sql.DB) *Feedback {
	return &Feedback{db: db}
}
