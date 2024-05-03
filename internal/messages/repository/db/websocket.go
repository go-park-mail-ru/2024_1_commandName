package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	chatsdelivery "ProjectMessenger/internal/chats/delivery"
)

type Websocket struct {
	db          *sql.DB
	Chats       *chatsdelivery.ChatsHandler
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func NewWsStorage(db *sql.DB) *Websocket {
	return &Websocket{
		db:          db,
		Connections: make(map[uint]*websocket.Conn),
	}
}

func UpgradeConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	return upgrader
}

func (m *Websocket) SendMessageToUser(userID uint, message []byte) error {
	connection := m.GetConnection(userID)
	if connection == nil {
		return errors.New("No connection found for user")
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (m *Websocket) AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context {
	fmt.Println("add con  for ", userID)
	m.mu.Lock()
	m.Connections[userID] = connection
	m.mu.Unlock()
	ctx = context.WithValue(ctx, "ws userID", userID)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	logger.Debug("established ws")
	return ctx
}

func (m *Websocket) DeleteConnection(userID uint) {
	m.mu.Lock()
	delete(m.Connections, userID)
	m.mu.Unlock()
}

func (m *Websocket) GetConnection(userID uint) *websocket.Conn {
	m.mu.RLock()
	conn := m.Connections[userID]
	m.mu.RUnlock()
	return conn
}
