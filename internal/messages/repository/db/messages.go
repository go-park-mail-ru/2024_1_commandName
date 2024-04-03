package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"ProjectMessenger/domain"

	"github.com/gorilla/websocket"
)

type Messages struct {
	db          *sql.DB
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func (m *Messages) PrintMessage(message []byte) {
	fmt.Print("От пользователя пришло сообщение: ")
	fmt.Println(string(message))
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

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		userDecodedMessage := DecodeJSON(message)
		userDecodedMessage.UserID = userID
		userDecodedMessage.CreateTimestamp = time.Now()
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

func (m *Messages) AddConnection(connection *websocket.Conn, userID uint) {
	m.mu.Lock()
	m.Connections[userID] = connection
	m.mu.Unlock()
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

func DecodeJSON(message []byte) domain.Message {
	var mess domain.Message
	err := json.Unmarshal(message, &mess)
	if err != nil {
		// TODO
		fmt.Println(err)
		return domain.Message{}
	}
	return mess
}

func (m *Messages) SetMessage(ctx context.Context, message domain.Message) {
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited, create_datetime) VALUES($1, $2, $3, $4, $5) "
	_, err := m.db.ExecContext(ctx, query, message.UserID, message.ChatID, message.Message, message.Edited, message.CreateTimestamp)
	if err != nil {
		// TODO
		fmt.Println(err)
	}
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
