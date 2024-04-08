package db

import (
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/chats/usecase"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"ProjectMessenger/domain"

	"github.com/gorilla/websocket"
)

type Messages struct {
	db          *sql.DB
	Chats       *chatsdelivery.ChatsHandler
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func UpgradeConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	return upgrader
}

func (m *Messages) ReadMessages(ctx context.Context, connection *websocket.Conn, userID uint) {
	defer func() {
		m.DeleteConnection(userID)
		connection.Close()
	}()
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		userDecodedMessage := decodeJSON(message)
		logger.Debug("got ws message", "msg", userDecodedMessage)
		userDecodedMessage.UserID = userID
		userDecodedMessage.CreateTimestamp = time.Now()
		m.SendMessageToUser(userID, []byte(`{"status": 200}`))
		m.SetMessage(ctx, userDecodedMessage)
	}
}

func (m *Messages) SendMessageToUser(userID uint, message []byte) error {
	connection := m.GetConnection(userID)
	if connection == nil {
		return errors.New("No connection found for user")
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (m *Messages) AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context {
	m.mu.Lock()
	m.Connections[userID] = connection
	m.mu.Unlock()
	ctx = context.WithValue(ctx, "ws userID", userID)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	logger.Debug("established ws")
	return ctx
}

func (m *Messages) DeleteConnection(userID uint) {
	m.mu.Lock()
	delete(m.Connections, userID)
	m.mu.Unlock()
}

func (m *Messages) GetConnection(userID uint) *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Connections[userID]
}

func NewMessageStorage(db *sql.DB) *Messages {
	return &Messages{
		db:          db,
		Connections: make(map[uint]*websocket.Conn),
	}
}

func decodeJSON(message []byte) domain.Message {
	var mess domain.Message
	err := json.Unmarshal(message, &mess)
	if err != nil {
		// TODO
		fmt.Println(err)
		return domain.Message{}
	}
	return mess
}

func (m *Messages) SendMessageToOtherUsers(ctx context.Context, message domain.Message, chatStorage usecase.ChatStore) {
	users := chatStorage.GetChatUsersByChatID(ctx, message.ChatID)
	for i := range users {
		conn := m.GetConnection(users[i].UserID)
		if conn != nil {
			messageMarshalled, err := json.Marshal(message)
			if err != nil {
				return
			}
			err = m.SendMessageToUser(users[i].UserID, messageMarshalled)
			if err != nil {
				return
			}
		}
	}

}

func (m *Messages) SetMessage(ctx context.Context, message domain.Message) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited, create_datetime) VALUES($1, $2, $3, $4, $5) "
	_, err := m.db.ExecContext(ctx, query, message.UserID, message.ChatID, message.Message, message.Edited, message.CreateTimestamp)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	query = `UPDATE chat.chat SET last_action_datetime = $1 WHERE id = $2`
	_, err = m.db.ExecContext(ctx, query, message.CreateTimestamp, message.ChatID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Debug("SetMessage: success", "msg", message)
}

func (m *Messages) GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)

	rows, err := m.db.QueryContext(ctx, "SELECT id, user_id, chat_id, message.message, edited, create_datetime FROM chat.message WHERE chat_id = $1 ORDER BY create_datetime DESC LIMIT $2", chatID, limit)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, profile.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var mess domain.Message
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.Edited, &mess.CreateTimestamp); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetMessagesByChatID, profile.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		chatMessagesArr = append(chatMessagesArr, mess)
	}
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, profile.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	return chatMessagesArr
}
