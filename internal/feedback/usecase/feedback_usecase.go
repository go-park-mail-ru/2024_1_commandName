package usecase

import (
	"context"

	"ProjectMessenger/domain"
)

type FeedbackStore interface {
	CheckNeedReturnQuestion(ctx context.Context, userID uint, question_id int) (needReturn bool)
	GetQuestions(ctx context.Context, userID uint) []domain.Question
	SetAnswer(ctx context.Context, userID uint, questionID int, grade int) bool
	AddQuestion(ctx context.Context, question domain.Question)
	UpdateQuestion(ctx context.Context, question domain.Question)
	GetAllQuestionStatistic(ctx context.Context) (statistic domain.AllStatistic)
	GetStatisticForOneQuestion(ctx context.Context, questionID int, questionType string) (statistic []int)
}

func IsReturnNeeded(ctx context.Context, fs FeedbackStore, userID uint, typeOfQuestion int) (isNeeded bool) {
	isNeeded = fs.CheckNeedReturnQuestion(ctx, userID, typeOfQuestion)
	return isNeeded
}

func ReturnQuestions(ctx context.Context, fs FeedbackStore, userID uint) []domain.Question {
	questions := fs.GetQuestions(ctx, userID)
	return questions
}

func SetAnswer(ctx context.Context, userID uint, questionID int, grade int, fs FeedbackStore) bool {
	ok := fs.SetAnswer(ctx, userID, questionID, grade)
	return ok
}

func GetAllStatistic(ctx context.Context, fs FeedbackStore) (Statistics domain.AllStatistic) {
	Statistics = fs.GetAllQuestionStatistic(ctx)
	return Statistics
}

func AddQuestion(ctx context.Context, fs FeedbackStore, question domain.Question) {
	fs.AddQuestion(ctx, question)
}

func UpdateQuestion(ctx context.Context, fs FeedbackStore, question domain.Question) {
	fs.UpdateQuestion(ctx, question)
}
