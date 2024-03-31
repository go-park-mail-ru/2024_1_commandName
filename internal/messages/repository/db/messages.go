package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"

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

func (m *Messages) ReadMessages(connection *websocket.Conn, userID uint) {
	defer func() {
		m.DeleteConnection(userID)
		connection.Close()
	}()

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}
		err = m.SendMessageToUser(userID, message)
		if err != nil {
			fmt.Println(err)
		}

		go m.PrintMessage(message)
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
