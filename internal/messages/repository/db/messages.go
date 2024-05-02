package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"ProjectMessenger/domain"
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

func (m *Messages) GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)

	rows, err := m.db.QueryContext(ctx, "SELECT message.id, user_id, chat_id, message.message, created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = $1", chatID)
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
