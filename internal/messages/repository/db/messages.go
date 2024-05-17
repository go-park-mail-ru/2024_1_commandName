package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"ProjectMessenger/domain"
	authusecase "ProjectMessenger/internal/auth/usecase"
	"ProjectMessenger/internal/misc"
)

type Messages struct {
	db *sql.DB
}

func NewMessageStorage(db *sql.DB) *Messages {
	return &Messages{db: db}
}

func (m *Messages) SetMessage(ctx context.Context, message domain.Message) (messageSaved domain.Message) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited_at, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id"
	var messageID uint
	err := m.db.QueryRowContext(ctx, query, message.UserID, message.ChatID, message.Message, message.EditedAt, message.CreatedAt).Scan(&messageID)
	if err != nil {
		// TODO
		fmt.Println(err)
		return domain.Message{}
	}
	message.ID = messageID
	logger.Debug("SetMessage: success", "msg", message)
	return message
}

func (m *Messages) SetFiles(ctx context.Context, multipartFiles []multipart.File, userID uint, messageID uint, userStorage authusecase.UserStore) error {
	_, found := userStorage.GetByUserID(ctx, userID)
	if !found {
		customErr := domain.CustomError{
			Type:    "find user by ID",
			Message: "user not found",
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr)
		return customErr
	}
	for _, multipartFile := range multipartFiles {
		buff := make([]byte, 512)
		if _, err := multipartFile.Read(buff); err != nil {
			customErr := domain.CustomError{
				Type:    "read multipart file",
				Message: err.Error(),
				Segment: "SetFiles, messages.go",
			}
			fmt.Println(customErr)
			return customErr
		}
		seek, err := multipartFile.Seek(0, io.SeekStart)
		if err != nil || seek != 0 {
			customErr := domain.CustomError{
				Type:    "seek multipart file",
				Message: err.Error(),
				Segment: "SetFiles, messages.go",
			}
			fmt.Println(customErr)
			return customErr
		}
		mimeType := http.DetectContentType(buff)
		fmt.Println(mimeType)

		//TODO check type of file
		/*
			if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/pjpeg" && mimeType != "image/webp" {
				return fmt.Errorf("Файл не является изображением")
			}*/
		query := "INSERT INTO chat.file (user_id, message_id, file_path) VALUES($1, $2, $3)"
		dbErr := m.db.QueryRowContext(ctx, query, userID, messageID, "")
		if dbErr != nil {
			// TODO
			fmt.Println(err)
			return domain.Message{}
		}
	}

}

func (m *Messages) StoreFile(ctx context.Context, multipartFile multipart.File, fileHandler *multipart.FileHeader) (name string, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	originalName := fileHandler.Filename
	fileNameSlice := strings.Split(originalName, ".")
	if len(fileNameSlice) < 2 {
		logger.Info("StoreAvatar filename without extension")
		return "", fmt.Errorf("Файл не имеет расширения")
	}
	extension := fileNameSlice[len(fileNameSlice)-1]

	filename := misc.RandStringRunes(20)
	filePath := "files/" + filename + "." + extension

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("StoreAvatar failed to open a file", "path", filePath)
		return "", fmt.Errorf("internal error")
	}
	defer f.Close()

	_, err = io.Copy(f, multipartFile)
	if err != nil {
		logger.Error("StoreFile failed to copy file", "path", filePath)
		return "", fmt.Errorf("internal error")
	}
	logger.Debug("StoreFile success", "path", filePath)
	return filename + "." + extension, nil
}

func (m *Messages) GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)

	rows, err := m.db.QueryContext(ctx, "SELECT message.id, user_id, chat_id, message.message, message.created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = $1", chatID)
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
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.CreatedAt, &mess.EditedAt, &mess.SenderUsername); err != nil {
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
	return chatMessagesArr
}

func (m *Messages) GetMessage(ctx context.Context, messageID uint) (message domain.Message, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	message = domain.Message{}
	err = m.db.QueryRowContext(ctx, "SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), created_at FROM chat.message WHERE id = $1", messageID).Scan(
		&message.ID, &message.UserID, &message.ChatID, &message.Message, &message.Edited, &message.EditedAt, &message.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("EditMessage didn't found message", "messageID", messageID)
			return message, fmt.Errorf("Такого сообщения не существует")
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetChatsForUser, chats.go",
		}
		logger.Error(err.Error(), "EditMessage db error", customErr.Message)
		return message, fmt.Errorf("internal error")
	}
	return message, nil
}

func (m *Messages) UpdateMessageText(ctx context.Context, message domain.Message) (err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	_, err = m.db.ExecContext(ctx, "UPDATE chat.message SET message = $1, edited = $2, edited_at = $3 WHERE id = $4", message.Message, message.Edited, message.EditedAt, message.ID)
	if err != nil {
		logger.Error("UpdateMessageText db error", "messageID", message.ID)
		return fmt.Errorf("internal error")
	}
	return nil
}

func (m *Messages) DeleteMessage(ctx context.Context, messageID uint) error {
	logger := slog.With("requestID", ctx.Value("traceID"))
	_, err := m.db.ExecContext(ctx, "DELETE FROM chat.message WHERE id = $1", messageID)
	if err != nil {
		logger.Error("DeleteMessage db error", "messageID", messageID)
		return fmt.Errorf("internal error")
	}
	return nil
}
