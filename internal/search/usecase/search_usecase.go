package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	_ "regexp"

	"ProjectMessenger/domain"
	_ "ProjectMessenger/internal/misc"
	"github.com/gorilla/websocket"
)

type SearchStore interface {
	AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn
	AddSearchIndexes(ctx context.Context)
	DeleteSearchIndexes(ctx context.Context)
	SearchChats(ctx context.Context, word string, userID uint) (foundChatsStructure domain.ChatSearchResponse)
	SendMatchedSearchResponse(response domain.ChatSearchResponse)
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person) {
	fmt.Println("add conn for", user.ID)
	ctx = s.AddConnection(ctx, connection, user.ID)
	defer func() {
		fmt.Println("del conn for", user.ID)
		s.DeleteConnection(user.ID)
		err := connection.Close()
		if err != nil {
			fmt.Println("err:", err)
			//TODO
			return
		}
	}()

	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	s.AddSearchIndexes(ctx)
	for {
		var decodedChatSearchRequest domain.ChatSearchRequest
		mt, request, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(request, &decodedChatSearchRequest)
		if err != nil {
			fmt.Println("err decoding JSON:", err)
			continue
		}
		decodedChatSearchRequest.UserID = user.ID
		logger.Debug("got ws message", "msg", decodedChatSearchRequest)
		//TODO: валидация
		matchedChatsStructure := s.SearchChats(ctx, decodedChatSearchRequest.Word, decodedChatSearchRequest.UserID)
		s.SendMatchedSearchResponse(matchedChatsStructure)
	}
	s.DeleteSearchIndexes(ctx)
}
