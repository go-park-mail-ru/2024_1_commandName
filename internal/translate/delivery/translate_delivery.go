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
	YandexConfig.TranslateKey = "Bearer t1.9euelZqSjI3ImYmRmJGXjZyRmI-YmO3rnpWanMyVzMzLyJuXnJSQzZSQzJnl9PdNQFhO-e86AQfv3fT3DW9VTvnvOgEH783n9euelZrPkZGTxovOlcqOyozGi8mNme_8xeuelZrPkZGTxovOlcqOyozGi8mNmQ.O4GVYAf7g3v5xPfzw4qACh4IOxnioX_fBrPl-8h0uCCGbAi6bbc4TcQ4CsT2lWNsVgAkyhk8zV4w0dyhBh7bCg"
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
