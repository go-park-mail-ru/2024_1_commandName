package usecase

import (
	"context"
	_ "regexp"

	"ProjectMessenger/domain"
	_ "ProjectMessenger/internal/misc"
	"github.com/gorilla/websocket"
)

type SearchStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string)
	AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn
	HandleWebSocket(ctx context.Context, connection *websocket.Conn, user domain.Person)
}
