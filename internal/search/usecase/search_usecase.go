package usecase

import (
	"context"
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
	SearchMessages(ctx context.Context, word string, userID uint, chatID uint) (foundMessagesStructure domain.MessagesSearchResponse)
	SendMatchedMessagesSearchResponse(response domain.MessagesSearchResponse, userID uint)
	SearchContacts(ctx context.Context, word string, userID uint) (foundContactsStructure domain.ContactsSearchResponse)
	SendMatchedContactsSearchResponse(response domain.ContactsSearchResponse, userID uint)
}

func SearchChats(ctx context.Context, s SearchStore, word string, userID uint) (matchedChatsStructure domain.ChatSearchResponse) {
	s.AddSearchIndexes(ctx)
	matchedChatsStructure = s.SearchChats(ctx, word, userID, "chat")
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return chats:", "chats", matchedChatsStructure)
	return matchedChatsStructure
}

func SearchChannels(ctx context.Context, s SearchStore, word string, userID uint) (matchedChatsStructure domain.ChatSearchResponse) {
	s.AddSearchIndexes(ctx)
	matchedChatsStructure = s.SearchChats(ctx, word, userID, "channel")
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return channels:", "channels", matchedChatsStructure)
	return matchedChatsStructure
}

func SearchMessages(ctx context.Context, s SearchStore, word string, userID uint, chatID uint) (matchedMessagesStructure domain.MessagesSearchResponse) {
	s.AddSearchIndexes(ctx)
	matchedMessagesStructure = s.SearchMessages(ctx, word, userID, chatID)
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return messages:", "messages", matchedMessagesStructure)
	return matchedMessagesStructure
}

func SearchContacts(ctx context.Context, s SearchStore, word string, userID uint) (matchedContactsStructure domain.ContactsSearchResponse) {
	s.AddSearchIndexes(ctx)
	matchedContactsStructure = s.SearchContacts(ctx, word, userID)
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("return contacts:", "contacts", matchedContactsStructure)
	return matchedContactsStructure
}
