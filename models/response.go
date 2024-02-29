package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int `json:"status" example:"200"`
	Body   any `json:"body,omitempty"`
}

type ErrorResponse struct {
	Status int   `json:"status" example:"400"`
	Body   Error `json:"body,omitempty"`
}

type Error struct {
	Error string `json:"error" example:"user not found"`
}

func WriteStatusJson(w http.ResponseWriter, status int, body any) error {
	w.WriteHeader(status)
	_, err := w.Write(MarshalStatusJson(status, body))
	if err != nil {
		return err
	}
	return nil
}

func MarshalStatusJson(status int, body any) []byte {
	response := Response{
		Status: status,
		Body:   body,
	}
	marshal, _ := json.Marshal(response)
	return marshal
}
