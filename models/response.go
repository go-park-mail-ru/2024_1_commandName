package models

import (
	"encoding/json"
	"net/http"
)

// Response[T]
type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Body   T   `json:"body,omitempty"`
}

type ErrorResponse struct {
	Status int   `json:"status" example:"400"`
	Body   Error `json:"body,omitempty"`
}

type Error struct {
	Error string `json:"error" example:"error description"`
}

type Chats struct {
	Chats []*Chat `json:"chats"`
}

func WriteStatusJson(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
	response := Response[any]{
		Status: status,
		Body:   body,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
