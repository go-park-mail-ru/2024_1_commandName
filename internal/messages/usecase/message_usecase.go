package usecase

import (
	"context"

	"ProjectMessenger/domain"
	"github.com/gorilla/websocket"
)

type MessageStore interface {
	PrintMessage(message []byte)
	ReadMessages(ctx context.Context, connection *websocket.Conn, userID uint)
	SendMessageToUser(userID uint, message []byte) error
	AddConnection(connection *websocket.Conn, userID uint)
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn

	GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message
}

func GetMessagesByWebSocket(ctx context.Context, connection *websocket.Conn, userID uint, messageStorage MessageStore) {
	messageStorage.AddConnection(connection, userID)
	messageStorage.ReadMessages(ctx, connection, userID)
}

func GetChatMessages(ctx context.Context, limit int, chatID uint, messageStorage MessageStore) []domain.Message {
	messages := messageStorage.GetChatMessages(ctx, chatID, limit)
	return messages
}
