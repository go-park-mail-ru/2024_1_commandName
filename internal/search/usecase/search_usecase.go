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
	SearchChats(ctx context.Context, word string, userID uint, chatType string) (foundChatsStructure domain.ChatSearchResponse)
	SendMatchedChatsSearchResponse(response domain.ChatSearchResponse, userID uint)
	SearchMessages(ctx context.Context, word string, userID uint) (foundMessagesStructure domain.MessagesSearchResponse)
	SendMatchedMessagesSearchResponse(response domain.MessagesSearchResponse, userID uint)
	SearchContacts(ctx context.Context, word string, userID uint) (foundContactsStructure domain.ContactsSearchResponse)
	SendMatchedContactsSearchResponse(response domain.ContactsSearchResponse, userID uint)
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, s SearchStore, user domain.Person) error {
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

	for {
		var decodedSearchRequest domain.SearchRequest
		mt, request, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		err = json.Unmarshal(request, &decodedSearchRequest)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "json Unmarshal",
				Message: err.Error(),
				Segment: "method HandleWebSocket, search_usecase.go",
			}
			fmt.Println(customErr.Error())
			continue
		}
		decodedSearchRequest.UserID = user.ID
		logger.Debug("got ws message", "msg", decodedSearchRequest)
		//TODO: валидация
		conn := s.GetConnection(user.ID)
		if conn == nil {
			fmt.Println("conn was closed")
			s.AddConnection(ctx, connection, user.ID)
		}
		if decodedSearchRequest.Type == "chat" {
			SearchChats(ctx, s, user, decodedSearchRequest.Word, decodedSearchRequest.UserID)
		} else if decodedSearchRequest.Type == "message" {
			SearchMessages(ctx, s, user, decodedSearchRequest.Word, decodedSearchRequest.UserID)
		} else if decodedSearchRequest.Type == "contact" {
			SearchContacts(ctx, s, user, decodedSearchRequest.Word, decodedSearchRequest.UserID)
		} else if decodedSearchRequest.Type == "channel" {
			SearchChannels(ctx, s, user, decodedSearchRequest.Word, decodedSearchRequest.UserID)
		} else {
			customErr := &domain.CustomError{
				Type:    "search type",
				Message: err.Error(),
				Segment: "method HandleWebSocket, search.go",
			}
			return customErr
		}
	}
	return nil
}

func SearchChats(ctx context.Context, s SearchStore, user domain.Person, word string, userID uint) {
	s.AddSearchIndexes(ctx)
	matchedChatsStructure := s.SearchChats(ctx, word, userID, "chat")
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return chats:", "chats", matchedChatsStructure)
	s.SendMatchedChatsSearchResponse(matchedChatsStructure, user.ID)
	s.DeleteSearchIndexes(ctx)
}

func SearchChannels(ctx context.Context, s SearchStore, user domain.Person, word string, userID uint) {
	s.AddSearchIndexes(ctx)
	matchedChatsStructure := s.SearchChats(ctx, word, userID, "channel")
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return channels:", "channels", matchedChatsStructure)
	s.SendMatchedChatsSearchResponse(matchedChatsStructure, user.ID)
	s.DeleteSearchIndexes(ctx)
}

func SearchMessages(ctx context.Context, s SearchStore, user domain.Person, word string, userID uint) {
	s.AddSearchIndexes(ctx)
	matchedMessagesStructure := s.SearchMessages(ctx, word, userID)
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return messages:", "messages", matchedMessagesStructure)
	s.SendMatchedMessagesSearchResponse(matchedMessagesStructure, user.ID)
	s.DeleteSearchIndexes(ctx)
}

func SearchContacts(ctx context.Context, s SearchStore, user domain.Person, word string, userID uint) {
	s.AddSearchIndexes(ctx)
	matchedContactsStructure := s.SearchContacts(ctx, word, userID)
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return contacts:", "contacts", matchedContactsStructure)
	s.SendMatchedContactsSearchResponse(matchedContactsStructure, user.ID)
	s.DeleteSearchIndexes(ctx)
}
