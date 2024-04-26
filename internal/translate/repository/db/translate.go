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
	req.Header.Add("Content-Type", t.Config.Header)
	req.Header.Add("Authorization", t.Config.TranslateKey)
	if err != nil {
		//TODO
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		//TODO
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		//TODO
	}
	response = ParseTranslateResponse(body)
	fmt.Println(response)
	return response
}

func ParseTranslateResponse(jsonResponse []byte) (response domain.TranslateResponse) {
	err := json.Unmarshal(jsonResponse, &response)
	if err != nil {
		//TODO
		fmt.Println(err)
	}
	return response
}

func (t *Translate) GetFolderID() string {
	return t.Config.FolderID
}

func NewTranslateStorage(database *sql.DB, YandexConfig domain.YandexConfig) *Translate {
	slog.Info("created search storage")
	return &Translate{
		db:     database,
		Config: YandexConfig,
	}
}
