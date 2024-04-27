package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"ProjectMessenger/domain"
)

type Feedback struct {
	db *sql.DB
}

func (f *Feedback) CheckNeedReturnQuestion(ctx context.Context, userID uint, question_id int) (needReturn bool) {
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
			Segment: "method CheckNeedReturnQuestion, feedback.go",
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
	logger.Debug("GetQuestions", "userID", userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetQuestions, feedback.go",
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

	logger.Debug("GetQuestions: success", "questions to return", questions)
	return questions
}

func NewFeedbackStorage(db *sql.DB) *Feedback {
	fillFake(db)
	return &Feedback{db: db}
}

func (f *Feedback) SetAnswer(ctx context.Context, answer domain.FeedbackAnswer) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "INSERT INTO feedback.survey_answers (user_id, question_id, grade, answered_at) VALUES ($1, $2, $3, $4)"
	_, err := f.db.ExecContext(ctx, query, answer.UserID, answer.QuestionID, answer.Grade, time.Now())
	logger.Debug("SetAnswer", "answer", answer)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method SetAnswer, feedback.go",
		}
		logger.Error(customErr.Error())
	}
	logger.Debug("SetAnswer: success", "adding answer", true)
}

func (f *Feedback) GetStatisticForCSAT(ctx context.Context) (statistic domain.FeedbackStatisticCSAT) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "SELECT grade FROM feedback.survey_answers sa JOIN feedback.survey_questions sq ON sa.question_id = sq.id WHERE sq.questiontype = $1"
	rows, err := f.db.QueryContext(ctx, query, "CSAT")
	logger.Debug("getStatisticForCSAT", "type", "CSAT")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method getStatisticForCSAT, feedback.go",
		}
		logger.Error(customErr.Error())
		return statistic
	}

	for rows.Next() {
		currGrade := 0
		err = rows.Scan(&currGrade)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		switch currGrade {
		case 1:
			statistic.OneStarsCount++
		case 2:
			statistic.TwoStarsCount++
		case 3:
			statistic.ThreeStarsCount++
		case 4:
			statistic.FourStarsCount++
		case 5:
			statistic.FiveStarsCount++
		}
	}
	return statistic
}

func (f *Feedback) GetStatisticForNSP(ctx context.Context) (statistic domain.FeedbackStatisticNSP) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "SELECT grade FROM feedback.survey_answers sa JOIN feedback.survey_questions sq ON sa.question_id = sq.id WHERE sq.questiontype = $1"
	rows, err := f.db.QueryContext(ctx, query, "NSP")
	logger.Debug("getStatisticForNSP", "type", "NSP")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method getStatisticForNSP, feedback.go",
		}
		logger.Error(customErr.Error())
		return statistic
	}

	for rows.Next() {
		currGrade := 0
		err = rows.Scan(&currGrade)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
		switch currGrade {
		case 1:
			statistic.OneStarsCount++
		case 2:
			statistic.TwoStarsCount++
		case 3:
			statistic.ThreeStarsCount++
		case 4:
			statistic.FourStarsCount++
		case 5:
			statistic.FiveStarsCount++
		case 6:
			statistic.SixStarsCount++
		case 7:
			statistic.SevenStarsCount++
		case 8:
			statistic.EightStarsCount++
		case 9:
			statistic.NineStarsCount++
		case 10:
			statistic.TenStarsCount++
		}
	}
	return statistic
}

func fillFake(db *sql.DB) {
	counter := 0
	_ = db.QueryRow("SELECT count(id) FROM feedback.survey_questions").Scan(&counter)
	if counter != 0 {
		return
	}
	query := `INSERT INTO feedback.survey_questions (question_text, questiontype) VALUES ($1, $2)`
	db.Exec(query, "Как вам наш сервис?", "CSAT")
}
