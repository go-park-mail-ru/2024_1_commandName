package usecase

import (
	"ProjectMessenger/internal/chats/usecase"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"ProjectMessenger/domain"

	"github.com/gorilla/websocket"
)

type WebsocketStore interface {
	SendMessageToUser(userID uint, message []byte) error
	AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context
	DeleteConnection(userID uint)
	GetConnection(userID uint) *websocket.Conn
}

type MessageStore interface {
	SetMessage(ctx context.Context, message domain.Message) (messageSaved domain.Message)
	GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, userID uint, wsStorage WebsocketStore, messageStorage MessageStore, chatStorage usecase.ChatStore) {
	ctx = wsStorage.AddConnection(ctx, connection, userID)
	defer func() {
		wsStorage.DeleteConnection(userID)
		connection.Close()
	}()
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}
		var userDecodedMessage domain.Message
		err = json.Unmarshal(message, &userDecodedMessage)
		if err != nil {
			fmt.Println(err)
			return
		}
		logger.Debug("got ws message", "msg", userDecodedMessage)
		userDecodedMessage.UserID = userID
		userDecodedMessage.CreateTimestamp = time.Now()
		messageSaved := messageStorage.SetMessage(ctx, userDecodedMessage)
		SendMessageToOtherUsers(ctx, messageSaved, wsStorage, chatStorage)
	}
}

func SendMessageToOtherUsers(ctx context.Context, message domain.Message, wsStorage WebsocketStore, chatStorage usecase.ChatStore) {
	chatUsers := chatStorage.GetChatUsersByChatID(ctx, message.UserID)
	wg := &sync.WaitGroup{}
	for i := range chatUsers {
		wg.Add(1)
		go func(message domain.Message, i int) {
			conn := wsStorage.GetConnection(chatUsers[i].UserID)
			if conn != nil {
				messageMarshalled, err := json.Marshal(message)
				if err != nil {
					return
				}
				err = wsStorage.SendMessageToUser(chatUsers[i].UserID, messageMarshalled)
				if err != nil {
					return
				}
			}
			wg.Done()
		}(message, i)
	}
	wg.Wait()
}

func GetChatMessages(ctx context.Context, limit int, chatID uint, messageStorage MessageStore) []domain.Message {
	messages := messageStorage.GetChatMessages(ctx, chatID, limit)
	return messages
}
