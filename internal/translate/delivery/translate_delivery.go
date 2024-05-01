package delivery

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"
	repo "ProjectMessenger/internal/translate/repository/db"
	"ProjectMessenger/internal/translate/usecase"
	"gopkg.in/yaml.v3"
)

type TranslateHandler struct {
	Translate    usecase.TranslateStore
	ChatsHandler *delivery.ChatsHandler
	Database     *sql.DB
	Config       domain.YandexConfig
}

func LoadConfig() domain.Config {
	envPath := os.Getenv("GOCHATME_HOME")
	slog.Debug("env home =" + envPath)
	f, err := os.Open(envPath + "config.yml")
	slog.Debug("trying to open " + envPath + "config.yml")
	if err != nil {
		slog.Error("load config failed", "err", err)
		panic(err)
	}
	defer f.Close()

	var cfg domain.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewTranslateHandler(database *sql.DB, chatsHandler *delivery.ChatsHandler) *TranslateHandler {
	var YandexConfig domain.YandexConfig
	cfg := LoadConfig()
	fmt.Println()
	fmt.Println("loaded config:", cfg)
	YandexConfig.TranslateKey = cfg.Yandex.TranslateKey
	YandexConfig.Url = cfg.Yandex.Url
	YandexConfig.FolderID = cfg.Yandex.FolderID
	YandexConfig.Header = cfg.Yandex.Header
	YandexConfig.Method = cfg.Yandex.Method
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
		err := errors.New("person not authorized")
		customErr := &domain.CustomError{
			Type:    "http new request",
			Message: err.Error(),
			Segment: "method TranslateMessage, translate_delivery.go",
		}
		fmt.Println(customErr.Error())
		return
	}
	var request domain.TranslateRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "http read response",
			Message: err.Error(),
			Segment: "method TranslateMessage, translate_delivery.go",
		}
		fmt.Println(customErr.Error())
	}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "json Unmarshal",
			Message: err.Error(),
			Segment: "method TranslateMessage, translate_delivery.go",
		}
		fmt.Println(customErr.Error())
	}
	request.FolderID = ts.Config.FolderID
	request.TargetLanguageCode = usecase.GetUserLanguageByID(ctx, ts.Database, ts.Translate, userID)
	response := usecase.HandleTranslate(ts.Translate, request)
	misc.WriteStatusJson(ctx, w, 200, response)
}
