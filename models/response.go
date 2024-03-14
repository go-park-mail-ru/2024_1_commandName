package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const INTERNALERRORJSON = "{\"status\": 500,\"body\":{\"error\": \"Internal server error\"}}"

// Response[T]
type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Body   T   `json:"body,omitempty"`
}

type Error struct {
	Error string `json:"error" example:"error description"`
}

type Chats struct {
	Chats []*Chat `json:"chats"`
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

func WriteInternalErrorJson(w http.ResponseWriter) {
	// implementation similar to http.Error, only difference is the Content-type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	_, _ = fmt.Fprintln(w, INTERNALERRORJSON)
}
