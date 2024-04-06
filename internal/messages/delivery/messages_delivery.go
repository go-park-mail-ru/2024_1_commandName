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

// SendMessage method to send messages
//
// @Summary SendMessage
// @Description Сначала по этому URL надо произвести upgrade до вебсокета, потом слать json сообщений
// @ID sendMessage
// @Accept application/json
// @Produce application/json
// @Param user body  domain.Message true "message that was sent"
// @Success 200 {object}  domain.Response[int]
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error | could not upgrade connection"
// @Router /sendMessage [post]
func (messageHandler *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := messageHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	fmt.Println(userID)

	upgrader := messageRepository.UpgradeConnection()

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}

	usecase.GetMessagesByWebSocket(ctx, connection, userID, messageHandler.Messages)
}

// GetChatMessages returns messages of some chat
//
// @Summary GetChatMessages
// @ID getChatMessages
// @Accept application/json
// @Produce application/json
// @Param user body  RequestChatIDBody true "ID of chat"
// @Success 200 {object}  domain.Response[domain.Messages]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChatMessages [post]
func (messageHandler *MessageHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, _ := messageHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var RequestChatID RequestChatIDBody
	err := decoder.Decode(&RequestChatID)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	limit := 100
	messages := usecase.GetChatMessages(r.Context(), limit, RequestChatID.ChatID, messageHandler.Messages)
	misc.WriteStatusJson(ctx, w, 200, domain.Messages{Messages: messages})
}
