package usecase

import (
	"context"

	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsByID(ctx context.Context, userID uint) []domain.Chat
}

func GetChatsForUser(ctx context.Context, userID uint, chatStorage ChatStore) []domain.Chat {
	chats := chatStorage.GetChatsByID(ctx, userID)
	return chats
}
