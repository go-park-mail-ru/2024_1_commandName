package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/messages/usecase"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Messages    usecase.MessageStore
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func NewMessagesHandler(authHandler *authdelivery.AuthHandler) *MessageHandler {
	return &MessageHandler{
		AuthHandler: authHandler,
		Connections: make(map[uint]*websocket.Conn),
		//Messages:       db.NewChatsStorage(dataBase),
	}
}

func (messageHandler MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	authorized, userID := messageHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	fmt.Println(userID)

	upgrader := UpgradeConnection()

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	messageHandler.AddConnection(connection, userID)
	messageHandler.readMessages(connection, userID)

}

func (m *MessageHandler) PrintMessage(message []byte) {
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

func (m *MessageHandler) readMessages(connection *websocket.Conn, userID uint) {
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

func (m *MessageHandler) SendMessageToUser(userID uint, message []byte) error {
	connection := m.GetConnection(userID)
	if connection == nil {
		return errors.New("No connection found for user")
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (m *MessageHandler) AddConnection(connection *websocket.Conn, userID uint) {
	m.mu.Lock()
	m.Connections[userID] = connection
	m.mu.Unlock()
}

func (m *MessageHandler) DeleteConnection(userID uint) {
	m.mu.Lock()
	delete(m.Connections, userID)
	m.mu.Unlock()
}

func (m *MessageHandler) GetConnection(userID uint) *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Connections[userID]
}
