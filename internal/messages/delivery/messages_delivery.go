package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/misc"

	//chatsInMemoryRepository "ProjectMessenger/internal/chats/repository/inMemory"
	messageRepository "ProjectMessenger/internal/messages/repository/db"
	"ProjectMessenger/internal/messages/usecase"
	"github.com/gorilla/websocket"
)

type RequestChatIDBody struct {
	ChatID uint `json:"chatID"`
}

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
	}
}

func NewMessagesHandlerMemory(authHandler *authdelivery.AuthHandler) *MessageHandler {
	return &MessageHandler{
		AuthHandler: authHandler,
		//Messages:    chatsInMemoryRepository.NewChatsStorage(),
	}
}

func (messageHandler MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
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

	usecase.GetMessagesByWebSocket(r.Context(), connection, userID, messageHandler.Messages)
}

func (messageHandler MessageHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	authorized, _ := messageHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var RequestChatID RequestChatIDBody
	err := decoder.Decode(&RequestChatID)
	if err != nil {
		http.Error(w, "wrong json structure", 400)
		fmt.Println(err)
		return
	}
	limit := 100
	messages := usecase.GetChatMessages(r.Context(), limit, RequestChatID.ChatID, messageHandler.Messages)
	misc.WriteStatusJson(w, 200, domain.Messages{Messages: messages})
}
