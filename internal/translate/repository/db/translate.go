package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
)

type Translate struct {
	db     *sql.DB
	Config domain.YandexConfig
}

func (t *Translate) Translate(request domain.TranslateRequest) (response domain.TranslateResponse) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		//TODO
		fmt.Println(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(t.Config.Method, t.Config.Url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "http new request",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	req.Header.Add("Content-Type", t.Config.Header)
	req.Header.Add("Authorization", t.Config.TranslateKey)
	resp, err := client.Do(req)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "http do request",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "read response",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	response = ParseTranslateResponse(body)
	return response
}

func ParseTranslateResponse(jsonResponse []byte) (response domain.TranslateResponse) {
	err := json.Unmarshal(jsonResponse, &response)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "json Unmarshal",
			Message: err.Error(),
			Segment: "method ParseTranslateResponse, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	return response
}

func (t *Translate) GetFolderID() string {
	return t.Config.FolderID
}

func NewTranslateStorage(database *sql.DB, YandexConfig domain.YandexConfig) *Translate {
	slog.Info("created translate storage")
	return &Translate{
		db:     database,
		Config: YandexConfig,
	}
}
