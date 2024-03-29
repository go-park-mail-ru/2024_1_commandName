package usecase

import (
	"context"
	"sort"

	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsByID(ctx context.Context, userID uint) []domain.Chat
}

func GetChatsForUser(ctx context.Context, userID uint, chatStorage ChatStore) []domain.Chat {
	chats := chatStorage.GetChatsByID(ctx, userID)
	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastMessage.CreateTimestamp.Before(chats[i].LastMessage.CreateTimestamp)
	})
	return chats
}
