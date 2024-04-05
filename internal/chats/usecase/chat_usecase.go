package usecase

import (
	"ProjectMessenger/internal/auth/usecase"
	"context"
	"fmt"
	"log/slog"
	"sort"

	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsForUser(ctx context.Context, userID uint) []domain.Chat
	GetChatUsersByChatID(ctx context.Context, chatID uint) []*domain.ChatUser
	CheckPrivateChatExists(ctx context.Context, userID1, userID2 uint) (exists bool, chatID uint, err error)
	GetChatByChatID(ctx context.Context, chatID uint) (domain.Chat, error)
	CreateChat(ctx context.Context, userIDs ...uint) (chatID uint, err error)
}

func GetChatByChatID(ctx context.Context, userID, chatID uint, chatStorage ChatStore, userStorage usecase.UserStore) (domain.Chat, error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chat, err := chatStorage.GetChatByChatID(ctx, chatID)
	if err != nil {
		return domain.Chat{}, err
	}
	belongs := CheckUserBelongsToChat(ctx, chatID, userID, chatStorage, userStorage)
	if !belongs {
		logger.Info("GetChatByChatID: user does not belong", "userID", userID, "chatID", chatID)
		return domain.Chat{}, fmt.Errorf("user does not belong to chat")
	}

	if chat.Type == "1" {
		name, _ := GetCompanionNameForPrivateChat(ctx, chatID, userID, chatStorage, userStorage)
		chat.Name = name
	}
	return chat, nil
}

func GetChatsForUser(ctx context.Context, userID uint, chatStorage ChatStore, userStorage usecase.UserStore) []domain.Chat {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chats := chatStorage.GetChatsForUser(ctx, userID)
	for i := range chats {
		if chats[i].Type == "1" {
			name, ok := GetCompanionNameForPrivateChat(ctx, chats[i].ID, userID, chatStorage, userStorage)
			if !ok {
				logger.Debug("GetChatsForUser: getchatname failed", "userID", userID, "chatID", chats[i].ID)
				continue
			}
			chats[i].Name = name
		}
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[j].LastMessage.CreateTimestamp.Before(chats[i].LastMessage.CreateTimestamp)
	})
	return chats
}

func CheckUserBelongsToChat(ctx context.Context, chatID uint, userRequestingID uint, chatStorage ChatStore, userStorage usecase.UserStore) bool {
	logger := slog.With("requestID", ctx.Value("traceID"))
	usersOfChat := chatStorage.GetChatUsersByChatID(ctx, chatID)
	for i := range usersOfChat {
		if usersOfChat[i].UserID == userRequestingID {
			logger.Debug("CheckUserBelongsToChat: true", "chatID", chatID, "userID", userRequestingID)
			return true
		}
	}
	logger.Debug("CheckUserBelongsToChat: false", "chatID", chatID, "userID", userRequestingID)
	return false
}

func GetCompanionNameForPrivateChat(ctx context.Context, chatID uint, userRequestingID uint, chatStorage ChatStore, userStorage usecase.UserStore) (companionUsername string, belongs bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	usersOfChat := chatStorage.GetChatUsersByChatID(ctx, chatID)
	if len(usersOfChat) != 2 {
		logger.Error("GetCompanionNameForPrivateChat: number of users in private chat doesn't equal two")
		return "", false
	}
	if usersOfChat[0].UserID == userRequestingID {
		user, found := userStorage.GetByUserID(ctx, usersOfChat[1].UserID)
		if !found {
			logger.Error("GetCompanionNameForPrivateChat: user0 not found")
		}
		return user.Username, true
	} else if usersOfChat[1].UserID == userRequestingID {
		user, found := userStorage.GetByUserID(ctx, usersOfChat[0].UserID)
		if !found {
			logger.Error("GetCompanionNameForPrivateChat: user1 not found")
		}
		return user.Username, true
	} else {
		logger.Error("GetCompanionNameForPrivateChat: private chat does not contain it's users")
		return "", false
	}
}

// CreatePrivateChat created chat, or returns existing
func CreatePrivateChat(ctx context.Context, creatingUserID uint, companionID uint, chatStorage ChatStore, userStorage usecase.UserStore) (chatID uint, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	if creatingUserID == companionID {
		return 0, fmt.Errorf("Диалог с самим собой пока не поддерживается")
	}

	companion, found := userStorage.GetByUserID(ctx, companionID)
	if !found {
		logger.Error("CreatePrivateChat: user wasn't found", "companionID", companionID)
		return 0, fmt.Errorf("Пользователь, с которым вы хотите создать диалог, не найден")
	}
	exists, chatID, err := chatStorage.CheckPrivateChatExists(ctx, creatingUserID, companion.ID)
	if err != nil {
		return 0, err
	}
	if exists {
		return chatID, nil
	}
	chatID, err = chatStorage.CreateChat(ctx, creatingUserID, companion.ID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
