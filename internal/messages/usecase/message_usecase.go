package usecase

import "github.com/gorilla/websocket"

type MessageStore interface {
	PrintMessage(message []byte)
	ReadMessages(connection *websocket.Conn, userID uint)
	SendMessageToUser(userID uint, message []byte) error
	AddConnection(connection *websocket.Conn, userID uint)
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn
}

func GetMessagesByWebSocket(connection *websocket.Conn, userID uint, messageStorage MessageStore) {
	messageStorage.AddConnection(connection, userID)
	messageStorage.ReadMessages(connection, userID)
}
