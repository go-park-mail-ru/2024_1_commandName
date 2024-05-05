package misc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
)

const INTERNALERRORJSON = "{\"status\": 500,\"body\":{\"error\": \"Internal server error\"}}"

func WriteStatusJson(ctx context.Context, w http.ResponseWriter, status int, body any) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	w.Header().Set("Content-Type", "application/json")
	jsonByte, err := MarshalStatusJson(status, body)
	if err != nil {
		WriteInternalErrorJson(ctx, w)
		return
	}
	_, err = w.Write(jsonByte)
	if err != nil {
		WriteInternalErrorJson(ctx, w)
		return
	}
	logger.Info("response", "status", status, "body", body)
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

func WriteInternalErrorJson(ctx context.Context, w http.ResponseWriter) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	// implementation similar to http.Error, only difference is the Content-type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	_, _ = fmt.Fprintln(w, INTERNALERRORJSON)
	logger.Info("response internal error", "body", INTERNALERRORJSON)
}

func WriteErrorMessageJson(ctx context.Context, w http.ResponseWriter, statusCode int, errorMessage string) {
	errorResponse := domain.Error{Error: errorMessage}
	WriteStatusJson(ctx, w, statusCode, errorResponse)
}
