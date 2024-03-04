package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int `json:"status" example:"200"`
	Body   any `validate:"optional" json:"body,omitempty"`
}

type ErrorResponse struct {
	Status int   `json:"status" example:"400"`
	Body   Error `json:"body,omitempty"`
}

type Error struct {
	Error string `json:"error" example:"error description"`
}

func WriteStatusJson(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	jsonByte, err := MarshalStatusJson(status, body)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonByte)
	if err != nil {
		return err
	}
	return nil
}

func MarshalStatusJson(status int, body any) ([]byte, error) {
	response := Response{
		Status: status,
		Body:   body,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
