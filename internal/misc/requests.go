package misc

import (
	"ProjectMessenger/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

const INTERNALERRORJSON = "{\"status\": 500,\"body\":{\"error\": \"Internal server error\"}}"

func WriteStatusJson(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	jsonByte, err := MarshalStatusJson(status, body)
	if err != nil {
		WriteInternalErrorJson(w)
		return
	}
	_, err = w.Write(jsonByte)
	if err != nil {
		WriteInternalErrorJson(w)
		return
	}
}

func MarshalStatusJson(status int, body any) ([]byte, error) {
	response := domain.Response[any]{
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
