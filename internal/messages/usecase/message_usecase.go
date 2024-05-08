package usecase

import (
	chats "ProjectMessenger/internal/chats_service/proto"
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
	GetMessage(ctx context.Context, messageID uint) (message domain.Message, err error)
	UpdateMessageText(ctx context.Context, message domain.Message) (err error)
	DeleteMessage(ctx context.Context, messageID uint) error
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, user domain.Person, wsStorage WebsocketStore, messageStorage MessageStore, chatStorage chats.ChatServiceClient) {
	ctx = wsStorage.AddConnection(ctx, connection, user.ID)
	defer func() {
		wsStorage.DeleteConnection(user.ID)
		connection.Close()
	}()
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		var userDecodedMessage domain.Message
		err = json.Unmarshal(message, &userDecodedMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}
		logger.Debug("got ws message", "msg", userDecodedMessage)
		userDecodedMessage.UserID = user.ID
		userDecodedMessage.CreatedAt = time.Now().UTC()
		userDecodedMessage.SenderUsername = user.Username
		messageSaved := messageStorage.SetMessage(ctx, userDecodedMessage)

		SendMessageToOtherUsers(ctx, messageSaved, wsStorage, chatStorage)
	}
}

func SendMessageToOtherUsers(ctx context.Context, message domain.Message, wsStorage WebsocketStore, chatStorage chats.ChatServiceClient) {
	//chatUsers := chatStorage.GetChatUsersByChatID(ctx, message.ChatID)
	resp, _ := chatStorage.GetChatByChatID(ctx, &chats.UserAndChatID{UserID: 0, ChatID: uint64(message.ChatID)})

	chatUsers := make([]domain.ChatUser, 0)
	for i := range resp.Messages {
		chatUsers = append(chatUsers, domain.ChatUser{
			ChatID: int(resp.Messages[i].ChatId),
			UserID: uint(resp.Messages[i].UserId),
		})
	}

	wg := &sync.WaitGroup{}
	for i := range chatUsers {
		wg.Add(1)
		go func(userID uint, i int, message domain.Message) {
			defer wg.Done()
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
		}(chatUsers[i].UserID, i, message)
	}
	wg.Wait()
}

func GetChatMessages(ctx context.Context, limit int, chatID uint, messageStorage MessageStore) []domain.Message {
	messages := messageStorage.GetChatMessages(ctx, chatID, limit)
	return messages
}

func EditMessage(ctx context.Context, userID uint, messageID uint, newMessageText string, messageStorage MessageStore) (err error) {
	message, err := messageStorage.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}
	if message.UserID != userID {
		return fmt.Errorf("Пользователь не является отправителем")
	}
	message.Message = newMessageText
	message.EditedAt = time.Now().UTC()
	message.Edited = true
	err = messageStorage.UpdateMessageText(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMessage(ctx context.Context, userID uint, messageID uint, messageStorage MessageStore) error {
	message, err := messageStorage.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}
	if message.UserID != userID {
		return fmt.Errorf("Пользователь не является отправителем")
	}
	err = messageStorage.DeleteMessage(ctx, messageID)
	if err != nil {
		return err
	}
	return nil
}
