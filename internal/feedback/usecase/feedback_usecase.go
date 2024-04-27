package usecase

import (
	"context"

	"ProjectMessenger/domain"
)

type FeedbackStore interface {
	CheckNeedReturnQuestion(ctx context.Context, userID uint, question_id int) (needReturn bool)
	GetQuestions(ctx context.Context, userID uint) []domain.Question
	SetAnswer(ctx context.Context, answer domain.FeedbackAnswer)
}

func IsReturnNeeded(ctx context.Context, fs FeedbackStore, userID uint, typeOfQuestion int) (isNeeded bool) {
	isNeeded = fs.CheckNeedReturnQuestion(ctx, userID, typeOfQuestion)
	return isNeeded
}

func ReturnQuestions(ctx context.Context, fs FeedbackStore, userID uint) []domain.Question {
	questions := fs.GetQuestions(ctx, userID)
	return questions
}

func SetQuestion(ctx context.Context, answer domain.FeedbackAnswer, fs FeedbackStore) {
	fs.SetAnswer(ctx, answer)
}
