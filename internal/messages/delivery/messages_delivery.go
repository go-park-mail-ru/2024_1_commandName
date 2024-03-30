package delivery

import (
	"database/sql"
	"fmt"
	"net/http"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/messages/usecase"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Messages    usecase.MessageStore
}

func NewMessagesHandler(authHandler *authdelivery.AuthHandler, dataBase *sql.DB) *MessageHandler {
	return &MessageHandler{
		AuthHandler: authHandler,
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

	defer connection.Close()

	messageHandler.readMessages(connection)

}

func (m *MessageHandler) PrintMessage(message []byte) {
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

func (m *MessageHandler) readMessages(connection *websocket.Conn) {
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		if err = connection.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}

		go m.PrintMessage(message)
	}
}
