package usecase

import (
	"context"
)

type FeedbackStore interface {
	CheckNeedCSAT(ctx context.Context, userID uint, typeOfQuestion int) (isUserAnswered bool)
}

func IsUserAnsweredOnGlobalCSAT(ctx context.Context, fs FeedbackStore, userID uint) (isAnswered bool) {
	typeOfCSAT := 1
	isAnswered = fs.CheckNeedCSAT(ctx, userID, typeOfCSAT)
	return isAnswered
}

func IsUserAnsweredOnFitchCSAT(ctx context.Context, fs FeedbackStore, userID uint) (isAnswered bool) {
	typeOfCSAT := 2
	isAnswered = fs.CheckNeedCSAT(ctx, userID, typeOfCSAT)
	return isAnswered
}
