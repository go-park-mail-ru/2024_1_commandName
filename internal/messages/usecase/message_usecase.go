package usecase

import (
	"context"

	"ProjectMessenger/domain"

	"github.com/gorilla/websocket"
)

type MessageStore interface {
	ReadMessages(ctx context.Context, connection *websocket.Conn, userID uint)
	SendMessageToUser(userID uint, message []byte) error
	AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn

	GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message
}

func GetMessagesByWebSocket(ctx context.Context, connection *websocket.Conn, userID uint, messageStorage MessageStore) {
	ctx = messageStorage.AddConnection(ctx, connection, userID)
	messageStorage.ReadMessages(ctx, connection, userID)
}

func GetChatMessages(ctx context.Context, limit int, chatID uint, messageStorage MessageStore) []domain.Message {
	messages := messageStorage.GetChatMessages(ctx, chatID, limit)
	return messages
}
