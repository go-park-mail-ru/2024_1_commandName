package usecase

import (
	"context"

	"github.com/gorilla/websocket"
)

type MessageStore interface {
	PrintMessage(message []byte)
	ReadMessages(ctx context.Context, connection *websocket.Conn, userID uint)
	SendMessageToUser(userID uint, message []byte) error
	AddConnection(connection *websocket.Conn, userID uint)
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn
}

func GetMessagesByWebSocket(ctx context.Context, connection *websocket.Conn, userID uint, messageStorage MessageStore) {
	messageStorage.AddConnection(connection, userID)
	messageStorage.ReadMessages(ctx, connection, userID)
}
