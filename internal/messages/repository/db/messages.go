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
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited, create_datetime) VALUES($1, $2, $3, $4, $5) returning id"
	var messageID uint
	m.db.QueryRowContext(ctx, query, message.UserID, message.ChatID, message.Message, message.Edited, message.CreateTimestamp).Scan(&messageID)
	message.ID = messageID
	query = `UPDATE chat.chat SET last_action_datetime = $1 WHERE id = $2`
	_, err := m.db.ExecContext(ctx, query, message.CreateTimestamp, message.ChatID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//m.SendMessageToOtherUsers(ctx, message)
	logger.Debug("SetMessage: success", "msg", message)
	return message
}

func (m *Messages) GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)

	rows, err := m.db.QueryContext(ctx, "SELECT id, user_id, chat_id, message.message, edited, create_datetime FROM chat.message WHERE chat_id = $1 ORDER BY create_datetime DESC LIMIT $2", chatID, limit)
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
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.Edited, &mess.CreateTimestamp); err != nil {
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
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, profile.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	return chatMessagesArr
}
