package usecase

import (
	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsByID(userID uint) []domain.Chat
}

func GetChatsForUser(userID uint, chatStorage ChatStore) []domain.Chat {
	chats := chatStorage.GetChatsByID(userID)
	return chats
}
