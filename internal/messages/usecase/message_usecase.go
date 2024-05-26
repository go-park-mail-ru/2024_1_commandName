package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"mime/multipart"
	"os"
	"time"

	users "ProjectMessenger/internal/auth/usecase"
	chats2 "ProjectMessenger/microservices/chats_service/proto"

	"ProjectMessenger/domain"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	SetFile(ctx context.Context, multipartFile multipart.File, userID uint, messageID uint, request domain.FileFromUser, userStorage users.UserStore, fileHandler *multipart.FileHeader) error
	GetFileByPath(filePath string) (file *os.File, fileInfo os.FileInfo)
	GetFilePathByMessageID(ctx context.Context, messageID uint) (filePath []string)
	GetAllStickers(ctx context.Context) (stickers []domain.Sticker)
	GetStickerPathByID(ctx context.Context, stickerID uint) (filePah string)
}

type FileWithInfo struct {
	fileInfo os.FileInfo
	file     *os.File
}

func HandleWebSocket(ctx context.Context, connection *websocket.Conn, user domain.Person, wsStorage WebsocketStore, messageStorage MessageStore, chatStorage chats2.ChatServiceClient, userStorage users.UserStore, firebase *firebase.App) {
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
			customErr := domain.CustomError{
				Type:    "json.Unmarshal",
				Message: err.Error(),
				Segment: "HandleWebSocket, messages_usecase.go",
			}
			fmt.Println(customErr.Error())
			continue
		}
		logger.Debug("got ws message", "msg", userDecodedMessage)
		userDecodedMessage.UserID = user.ID
		userDecodedMessage.CreatedAt = time.Now().UTC()
		userDecodedMessage.SenderUsername = user.Username
		userDecodedMessage.StickerPath = ""
		messageSaved := messageStorage.SetMessage(ctx, userDecodedMessage)
		chatStorage.UpdateLastActionTime(ctx, &chats2.LastAction{
			ChatID: uint64(userDecodedMessage.ChatID),
			Time:   timestamppb.New(userDecodedMessage.CreatedAt),
		})

		SendMessageToOtherUsers(ctx, messageSaved, user.ID, wsStorage, chatStorage, userStorage, firebase)
	}
}

func SendNotification(app *firebase.App, token string, messageText string, senderUsername string) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	registrationToken := token

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: senderUsername,
			Body:  messageText,
		},
		Token: registrationToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully sent message:", response)
}

func SendMessageToOtherUsers(ctx context.Context, message domain.Message, userID uint, wsStorage WebsocketStore, chatStorage chats2.ChatServiceClient, userStorage users.UserStore, firebase *firebase.App) {
	resp, _ := chatStorage.GetChatByChatID(ctx, &chats2.UserAndChatID{UserID: uint64(userID), ChatID: uint64(message.ChatID)})

	chatUsers := make([]domain.ChatUser, 0)
	for i := range resp.Users {
		chatUsers = append(chatUsers, domain.ChatUser{
			ChatID: int(resp.Users[i].ChatId),
			UserID: uint(resp.Users[i].UserId),
		})
	}

	for i := range chatUsers {
		func(userID uint, i int, senderID uint, message domain.Message) {
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
			if userID != senderID {
				tokensForUser, _ := userStorage.GetTokensForUser(ctx, userID)
				for j := range tokensForUser {
					slog.Info("sending notificatoin")
					SendNotification(firebase, tokensForUser[j], message.Message, message.SenderUsername)
				}
			}
		}(chatUsers[i].UserID, i, userID, message)
	}
}

func SetFile(ctx context.Context, file multipart.File, userID uint, fileHeader *multipart.FileHeader, request domain.FileFromUser, messageStorage MessageStore, userStorage users.UserStore, wsStorage WebsocketStore, chatStorage chats2.ChatServiceClient, firebase *firebase.App) {
	user, found := userStorage.GetByUserID(ctx, userID)
	if !found {
		return
	}

	dummyMessage := domain.Message{
		ID:             0,
		ChatID:         request.ChatID,
		UserID:         user.ID,
		Message:        request.MessageText,
		Edited:         false,
		EditedAt:       time.Time{},
		CreatedAt:      time.Now().UTC(),
		SenderUsername: user.Username,
		File:           nil,
		StickerPath:    "",
	}
	messageSaved := messageStorage.SetMessage(ctx, dummyMessage)

	messageStorage.SetFile(ctx, file, user.ID, messageSaved.ID, request, userStorage, fileHeader)
	SendMessageToOtherUsers(ctx, messageSaved, user.ID, wsStorage, chatStorage, userStorage, firebase)
}

func GetAllStickers(ctx context.Context, messageStorage MessageStore) (stickers []domain.Sticker) {
	stickers = messageStorage.GetAllStickers(ctx)
	return stickers
}

func SendSticker(ctx context.Context, messageStore MessageStore, wsStorage WebsocketStore, chatStorage chats2.ChatServiceClient, request domain.FileFromUser, user domain.Person, userStorage users.UserStore, firebase *firebase.App) {
	stickerPath := messageStore.GetStickerPathByID(ctx, request.FileID)
	sticker := &domain.FileInMessage{Path: stickerPath, Type: "sticker"}
	stickerMessage := domain.Message{
		ID:             0,
		ChatID:         request.ChatID,
		UserID:         user.ID,
		Message:        "",
		Edited:         false,
		EditedAt:       time.Time{},
		CreatedAt:      time.Now().UTC(),
		SenderUsername: user.Username,
		File:           sticker,
		StickerPath:    stickerPath,
	}

	fmt.Println("end of fun")
	messageSaved := messageStore.SetMessage(ctx, stickerMessage)
	SendMessageToOtherUsers(ctx, messageSaved, user.ID, wsStorage, chatStorage, userStorage, firebase)
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
