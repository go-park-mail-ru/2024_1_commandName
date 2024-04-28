package usecase

import (
	"context"
	"database/sql"
	"fmt"
	_ "regexp"

	"ProjectMessenger/domain"
	_ "ProjectMessenger/internal/misc"
)

type TranslateStore interface {
	Translate(request domain.TranslateRequest) (response domain.TranslateResponse)
	GetFolderID() string
}

func HandleTranslate(ts TranslateStore, request domain.TranslateRequest) (response domain.TranslateResponse) {
	response = ts.Translate(request)
	return response
}

func GetUserLanguageByID(ctx context.Context, db *sql.DB, userID uint) (languageCode string) {
	err := db.QueryRowContext(ctx, "SELECT language FROM auth.person WHERE id = $1", userID).Scan(&languageCode)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetUserLanguageByID, translate_usecase.go",
		}
		fmt.Println(customErr.Error())
	}
	return languageCode
}
