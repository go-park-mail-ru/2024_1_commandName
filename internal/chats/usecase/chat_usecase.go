package usecase

import (
	"ProjectMessenger/domain"
	"sort"
)

type ChatStore interface {
	GetChatsPreviewsByID(userID uint) []domain.Chat
}

func GetChatsPreviewsForUser(userID uint, chatStorage ChatStore) []domain.Chat {
	chats := chatStorage.GetChatsPreviewsByID(userID)
	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastMessage.SentAt.Before(chats[i].LastMessage.SentAt)
	})
	return chats
}
