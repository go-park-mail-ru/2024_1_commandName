package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"ProjectMessenger/domain"
	"github.com/gorilla/websocket"
)

type Search struct {
	db          *sql.DB
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func (s *Search) GetUserIDbySessionID(ctx context.Context, sessionID string) {

}

func UpgradeConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	return upgrader
}

func (s *Search) AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context {
	s.mu.Lock()
	s.Connections[userID] = connection
	s.mu.Unlock()
	ctx = context.WithValue(ctx, "ws userID", userID)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	logger.Debug("established ws")
	return ctx
}

func (s *Search) DeleteConnection(userID uint) {
	s.mu.Lock()
	delete(s.Connections, userID)
	s.mu.Unlock()
}

func (s *Search) GetConnection(userID uint) *websocket.Conn {
	s.mu.RLock()
	conn := s.Connections[userID]
	s.mu.RUnlock()
	return conn
}

func (s *Search) HandleWebSocket(ctx context.Context, connection *websocket.Conn, user domain.Person) {
	ctx = s.AddConnection(ctx, connection, user.ID)
	defer func() {
		s.DeleteConnection(user.ID)
		err := connection.Close()
		if err != nil {
			fmt.Println("err:", err)
			//TODO
			return
		}
	}()

	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	_, err := s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_c ON chat.chat (id)")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_user_id ON chat.chat_user (user_id)")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_cu ON chat.chat_user (chat_id);")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	for {
		mt, wordToSearch, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		var decodedWordToSearch string
		err = json.Unmarshal(wordToSearch, &decodedWordToSearch)
		if err != nil {
			fmt.Println("err decoding JSON:", err)
			continue
		}
		logger.Debug("got ws message", "msg", decodedWordToSearch)
		//TODO: валидация
		messageSaved := messageStorage.SetMessage(ctx, decodedWordToSearch)
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_c RESTRICT")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_user_id RESTRICT")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_cu RESTRICT")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
}

func (s *Search) searchChats(ctx context.Context, word string, userID uint) (matchedChats []*domain.Chat) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id 
				FROM chat.chat c
				JOIN chat.chat_user cu ON c.id = cu.chat_id 
				WHERE name LIKE '%' || $1 AND cu.user_id = $2`, word, userID)
	if err != nil {
		//TODO
		fmt.Println("err:", err)
	}
	for rows.Next() {

	}
}

func NewSearchStorage(db *sql.DB) *Search {
	slog.Info("created search storage")
	return &Search{
		db:          db,
		Connections: make(map[uint]*websocket.Conn),
	}
}
