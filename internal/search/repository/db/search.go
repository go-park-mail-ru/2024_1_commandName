package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/repository/db"
	"ProjectMessenger/internal/chats/usecase"
	ws "ProjectMessenger/internal/messages/repository/db"
	"github.com/gorilla/websocket"
)

type Search struct {
	db          *sql.DB
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
	Chats       usecase.ChatStore
	WebSocket   *ws.Websocket
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
	s.WebSocket.Connections[userID] = connection
	s.mu.Unlock()
	ctx = context.WithValue(ctx, "ws userID", userID)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	logger.Debug("established ws")
	return ctx
}

func (s *Search) DeleteConnection(userID uint) {
	s.mu.Lock()
	delete(s.WebSocket.Connections, userID)
	s.mu.Unlock()
}

func (s *Search) GetConnection(userID uint) *websocket.Conn {
	s.mu.RLock()
	conn := s.WebSocket.Connections[userID]
	s.mu.RUnlock()
	return conn
}

func (s *Search) SendMessageToUser(userID uint, message []byte) error {
	connection := s.GetConnection(userID)
	if connection == nil {
		return errors.New("No connection found for user")
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (s *Search) SearchChats(ctx context.Context, word string, userID uint) (foundChatsStructure domain.ChatSearchResponse) {
	enWord := word

	rows, err := s.db.QueryContext(ctx,
		`SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id 
				FROM chat.chat c
				JOIN chat.chat_user cu ON c.id = cu.chat_id 
				WHERE name ILIKE $1 || '%' AND cu.user_id = $2`, enWord, userID)
	if err != nil {
		//TODO
		fmt.Println("err:", err)
	}
	matchedChats := make([]domain.Chat, 0)
	for rows.Next() {
		var mChat domain.Chat
		err = rows.Scan(&mChat.ID, &mChat.Type, &mChat.Name, &mChat.Description, &mChat.AvatarPath, &mChat.CreatedAt, &mChat.LastActionDateTime, &mChat.CreatorID)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method searchChats, search.go",
			}
			fmt.Println(customErr.Error())
			return foundChatsStructure
		}
		mChat.Messages = s.Chats.GetMessagesByChatID(ctx, mChat.ID)
		if mChat.Messages != nil {
			mChat.Users = s.Chats.GetChatUsersByChatID(ctx, mChat.ID)
		}

		if mChat.Users != nil {
			matchedChats = append(matchedChats, mChat)
		}
	}
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method searchChats, search.go",
		}
		fmt.Println("ERROR: ", customErr.Error())
		return foundChatsStructure
	}

	var chatSearchResponse domain.ChatSearchResponse
	chatSearchResponse.Chats = matchedChats
	chatSearchResponse.UserID = userID

	return chatSearchResponse
}

func ConvertToJSONResponse(chats []domain.Chat, userID uint) (jsonResponse []byte) {
	var chatSearchResponse domain.ChatSearchResponse
	chatSearchResponse.Chats = chats
	chatSearchResponse.UserID = userID
	jsonResponse, err := json.Marshal(chatSearchResponse)
	if err != nil {
		fmt.Println("err encoding JSON:", err)
	}
	return jsonResponse
}

func (s *Search) AddSearchIndexes(ctx context.Context) {
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
}

func (s *Search) DeleteSearchIndexes(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_c CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_user_id CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_cu CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
}

func (s *Search) SendMatchedSearchResponse(response domain.ChatSearchResponse) {

	jsonResp := ConvertToJSONResponse(response.Chats, response.UserID)
	err := s.WebSocket.SendMessageToUser(response.UserID, jsonResp)
	if err != nil {
		//TODO
		fmt.Println("ERROR:", err)
	}
}

func NewSearchStorage(database *sql.DB) *Search {
	slog.Info("created search storage")
	return &Search{
		db:          database,
		Connections: make(map[uint]*websocket.Conn),
		Chats:       db.NewChatsStorage(database),
		WebSocket:   ws.NewWsStorage(database),
	}
}
