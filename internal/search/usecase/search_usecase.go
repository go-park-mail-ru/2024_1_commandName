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
	SendMatchedChatsSearchResponse(response domain.ChatSearchResponse, userID uint)
	SearchMessages(ctx context.Context, word string, userID uint) (foundMessagesStructure domain.MessagesSearchResponse)
	SendMatchedMessagesSearchResponse(response domain.MessagesSearchResponse, userID uint)
	SearchContacts(ctx context.Context, word string, userID uint) (foundContactsStructure domain.ContactsSearchResponse)
	SendMatchedContactsSearchResponse(response domain.ContactsSearchResponse, userID uint)
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person, typeToSearch string) {
	fmt.Println("add conn for", user.ID)
	ctx = s.AddConnection(ctx, connection, user.ID)
	defer func() {
		fmt.Println("del conn for", user.ID)
		s.DeleteConnection(user.ID)
		err := connection.Close()
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "websocket close connection",
				Message: err.Error(),
				Segment: "method HandleWebSocket, search_usecase.go",
			}
			fmt.Println(customErr.Error())
			return
		}
	}()

	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	if typeToSearch == "chat" {
		SearchChats(ctx, connection, s, user, logger)
	}
	if typeToSearch == "message" {
		SearchMessages(ctx, connection, s, user, logger)
	}
	if typeToSearch == "contact" {
		SearchContacts(ctx, connection, s, user, logger)
	}
}

func SearchChats(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person, logger *slog.Logger) {
	s.AddSearchIndexes(ctx)
	for {
		var decodedChatSearchRequest domain.ChatSearchRequest
		mt, request, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(request, &decodedChatSearchRequest)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "json Unmarshal",
				Message: err.Error(),
				Segment: "method HandleWebSocket, search_usecase.go",
			}
			fmt.Println(customErr.Error())
			continue
		}
		decodedChatSearchRequest.UserID = user.ID
		logger.Debug("got ws message", "msg", decodedChatSearchRequest)
		//TODO: валидация
		matchedChatsStructure := s.SearchChats(ctx, decodedChatSearchRequest.Word, decodedChatSearchRequest.UserID)
		s.SendMatchedChatsSearchResponse(matchedChatsStructure, user.ID)
	}
	s.DeleteSearchIndexes(ctx)
}

func SearchMessages(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person, logger *slog.Logger) {
	s.AddSearchIndexes(ctx)
	for {
		var decodedMessageSearchRequest domain.MessagesSearchRequest
		mt, request, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(request, &decodedMessageSearchRequest)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "json Unmarshal",
				Message: err.Error(),
				Segment: "method HandleWebSocket, search_usecase.go",
			}
			fmt.Println(customErr.Error())
			continue
		}
		decodedMessageSearchRequest.UserID = user.ID
		logger.Debug("got ws message", "msg", decodedMessageSearchRequest)
		//TODO: валидация
		matchedMessagesStructure := s.SearchMessages(ctx, decodedMessageSearchRequest.Word, decodedMessageSearchRequest.UserID)
		s.SendMatchedMessagesSearchResponse(matchedMessagesStructure, user.ID)
	}
	s.DeleteSearchIndexes(ctx)
}

func SearchContacts(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person, logger *slog.Logger) {
	s.AddSearchIndexes(ctx)
	for {
		var decodedContactSearchRequest domain.ContactsSearchRequest
		mt, request, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(request, &decodedContactSearchRequest)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "json Unmarshal",
				Message: err.Error(),
				Segment: "method SearchContacts, search_usecase.go",
			}
			fmt.Println(customErr.Error())
			continue
		}
		decodedContactSearchRequest.UserID = user.ID
		logger.Debug("got ws message", "msg", decodedContactSearchRequest)
		//TODO: валидация
		matchedContactsStructure := s.SearchContacts(ctx, decodedContactSearchRequest.Word, decodedContactSearchRequest.UserID)
		s.SendMatchedContactsSearchResponse(matchedContactsStructure, user.ID)
	}
	s.DeleteSearchIndexes(ctx)
}
