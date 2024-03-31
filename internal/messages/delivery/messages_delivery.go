package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	//chatsInMemoryRepository "ProjectMessenger/internal/chats/repository/inMemory"
	messageRepository "ProjectMessenger/internal/messages/repository/db"
	"ProjectMessenger/internal/messages/usecase"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Messages    usecase.MessageStore
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
}

func NewMessagesHandler(authHandler *authdelivery.AuthHandler, database *sql.DB) *MessageHandler {
	return &MessageHandler{
		AuthHandler: authHandler,
		Connections: make(map[uint]*websocket.Conn),
		Messages:    messageRepository.NewMessageStorage(database),
		//Messages:       db.NewChatsStorage(dataBase),
	}
}

func NewMessagesHandlerMemory(authHandler *authdelivery.AuthHandler) *MessageHandler {
	return &MessageHandler{
		AuthHandler: authHandler,
		//Messages:    chatsInMemoryRepository.NewChatsStorage(),
	}
}

func (messageHandler MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	authorized, userID := messageHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	fmt.Println(userID)

	upgrader := messageRepository.UpgradeConnection()

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	usecase.GetMessagesByWebSocket(connection, userID, messageHandler.Messages)
	//messageHandler.AddConnection(connection, userID)
	//messageHandler.ReadMessages(connection, userID)

}
