package db

import (
	"ProjectMessenger/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type Messages struct {
	db *sql.DB
}

func NewMessageStorage(db *sql.DB) *Messages {
	return &Messages{db: db}
}

func (m *Messages) SetMessage(ctx context.Context, message domain.Message) (messageSaved domain.Message) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	fmt.Println("MESSAGE:", message)
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited_at, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id"
	var messageID uint
	err := m.db.QueryRowContext(ctx, query, message.UserID, message.ChatID, message.Message, message.EditedAt, message.CreatedAt).Scan(&messageID)
	if err != nil {
		fmt.Println("ARGS:", message.UserID, message.ChatID, message.Message, message.EditedAt, message.CreatedAt)
		fmt.Println(err)
		return domain.Message{}
	}
	fmt.Println("made insert", messageID)
	message.ID = messageID
	query = "UPDATE chat.chat SET created_at = $1 WHERE id = $2"
	_, err = m.db.ExecContext(ctx, query, message.CreatedAt, message.ChatID)
	fmt.Println("made update")
	if err != nil {
		fmt.Println("err in SetMessage")
		logger.Error(err.Error())
		return
	}
	//m.SendMessageToOtherUsers(ctx, message)
	logger.Debug("SetMessage: success", "msg", message)
	fmt.Println("return")
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
