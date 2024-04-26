package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"
	repo "ProjectMessenger/internal/translate/repository/db"
	"ProjectMessenger/internal/translate/usecase"
)

type TranslateHandler struct {
	Translate    usecase.TranslateStore
	ChatsHandler *delivery.ChatsHandler
	Database     *sql.DB
	Config       domain.YandexConfig
}

func NewTranslateHandler(database *sql.DB, chatsHandler *delivery.ChatsHandler) *TranslateHandler {
	var YandexConfig domain.YandexConfig
	YandexConfig.TranslateKey = "Bearer t1.9euelZqelYrMyciLnJDHj5PKzpyclO3rnpWanMyVzMzLyJuXnJSQzZSQzJnl8_dlHlBO-e80ShNo_d3z9yVNTU757zRKE2j9zef1656VmozOzZPGlMidmZTHjcjNk86e7_zF656VmozOzZPGlMidmZTHjcjNk86e.dbhRbkheLJfmVeunG45CqjxpeIosd9qEl3g0HlRvQSQBnn3QvPOBklVEm5GxoOUKTBWvWJIxBTsOXvRb9fOIDA"
	YandexConfig.Url = "https://translate.api.cloud.yandex.net/translate/v2/translate"
	YandexConfig.FolderID = "b1gq4i9e5unl47m0kj5f"
	YandexConfig.Header = "application/json"
	YandexConfig.Method = "POST"
	return &TranslateHandler{
		Database:     database,
		Config:       YandexConfig,
		ChatsHandler: chatsHandler,
		Translate:    repo.NewTranslateStorage(database, YandexConfig),
	}
}

func (ts *TranslateHandler) TranslateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := ts.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		fmt.Println("not auth")
	}
	var request domain.TranslateRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		//TODO
	}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		fmt.Println(err)
		//TODO
	}
	request.FolderID = ts.Config.FolderID
	request.TargetLanguageCode = usecase.GetUserLanguageByID(ctx, ts.Database, userID)
	fmt.Println(request)
	response := usecase.HandleTranslate(ts.Translate, request)
	fmt.Println("translated: ", response)
	misc.WriteStatusJson(ctx, w, 200, response)
}
