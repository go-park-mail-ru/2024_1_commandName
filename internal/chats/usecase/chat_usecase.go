package usecase

import (
	"ProjectMessenger/internal/auth/usecase"
	"context"
	"log/slog"
	"sort"

	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsByID(ctx context.Context, userID uint) []domain.Chat
	GetChatUsersByChatID(ctx context.Context, chatID int) []*domain.ChatUser
}

func GetChatsForUser(ctx context.Context, userID uint, chatStorage ChatStore, userStorage usecase.UserStore) []domain.Chat {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chats := chatStorage.GetChatsByID(ctx, userID)
	for i := range chats {
		if chats[i].Type == "1" {
			usersOfChat := chatStorage.GetChatUsersByChatID(ctx, chats[i].ID)
			if len(usersOfChat) != 2 {
				logger.Error("GetChatsForUser: number of users in private chat doesn't equal two")
				return nil
			}
			if usersOfChat[0].UserID == userID {
				user, found := userStorage.GetByUserID(ctx, usersOfChat[1].UserID)
				if !found {
					logger.Error("GetChatsForUser: user0 not found")
				}
				chats[i].Name = user.Username
			} else if usersOfChat[1].UserID == userID {
				user, found := userStorage.GetByUserID(ctx, usersOfChat[0].UserID)
				if !found {
					logger.Error("GetChatsForUser: user1 not found")
				}
				chats[i].Name = user.Username
			} else {
				logger.Error("GetChatsForUser: private chat does not contain it's users")
			}
		}
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastMessage.CreateTimestamp.Before(chats[i].LastMessage.CreateTimestamp)
	})
	return chats
}
