package usecase

import (
	"context"
	"database/sql"
	_ "regexp"

	"ProjectMessenger/domain"
	_ "ProjectMessenger/internal/misc"
)

type TranslateStore interface {
	Translate(request domain.TranslateRequest) (response domain.TranslateResponse)
	GetFolderID() string
	GetUserLanguageByID(ctx context.Context, db *sql.DB, userID uint) (languageCode string)
}

func HandleTranslate(ts TranslateStore, request domain.TranslateRequest) (response domain.TranslateResponse) {
	response = ts.Translate(request)
	return response
}

func GetUserLanguageByID(ctx context.Context, db *sql.DB, ts TranslateStore, userID uint) (languageCode string) {
	languageCode = ts.GetUserLanguageByID(ctx, db, userID)
	return languageCode
}
