package usecase

import (
	"ProjectMessenger/domain"
	"sort"
)

type ChatStore interface {
	GetChatsByID(userID uint) []domain.Chat
}

func GetChatsForUser(userID uint, chatStorage ChatStore) []domain.Chat {
	chats := chatStorage.GetChatsByID(userID)
	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastMessageSentAt.Before(chats[i].LastMessageSentAt)
	})
	return chats
}
