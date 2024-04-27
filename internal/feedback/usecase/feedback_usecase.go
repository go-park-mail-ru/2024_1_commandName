package usecase

import (
	"context"

	"ProjectMessenger/domain"
)

type FeedbackStore interface {
	CheckNeedReturnQuestion(ctx context.Context, userID uint, question_id int) (needReturn bool)
	GetQuestions(ctx context.Context, userID uint) []domain.Question
	SetAnswer(ctx context.Context, userID uint, questionID int, grade int) bool
	GetStatisticForCSAT(ctx context.Context) (statistic []int)
	GetStatisticForNSP(ctx context.Context) (statistic []int)
	AddQuestion(ctx context.Context, question domain.Question)
	UpdateQuestion(ctx context.Context, question domain.Question)
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
	Statistics.CSAT = fs.GetStatisticForCSAT(ctx)
	Statistics.NSP = fs.GetStatisticForNSP(ctx)
	return Statistics
}

func AddQuestion(ctx context.Context, fs FeedbackStore, question domain.Question) {
	fs.AddQuestion(ctx, question)
}

func UpdateQuestion(ctx context.Context, fs FeedbackStore, question domain.Question) {
	fs.UpdateQuestion(ctx, question)
}
