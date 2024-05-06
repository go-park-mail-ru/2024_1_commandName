package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"ProjectMessenger/internal/auth/usecase"

	"ProjectMessenger/domain"
)

type ChatStore interface {
	GetChatsForUser(ctx context.Context, userID uint) []domain.Chat
	GetChatUsersByChatID(ctx context.Context, chatID uint) []*domain.ChatUser
	CheckPrivateChatExists(ctx context.Context, userID1, userID2 uint) (exists bool, chatID uint, err error)
	GetChatByChatID(ctx context.Context, chatID uint) (domain.Chat, error)
	CreateChat(ctx context.Context, name, description string, userIDs ...uint) (chatID uint, err error)
	DeleteChat(ctx context.Context, chatID uint) (wasDeleted bool, err error)
	UpdateGroupChat(ctx context.Context, updatedChat domain.Chat) (ok bool)
	GetLastSeenMessageId(ctx context.Context, chatID uint, userID uint) (lastSeenMessageID int)
	GetFirstChatMessageID(ctx context.Context, chatID uint) (firstMessageID int)

	GetNPopularChannels(ctx context.Context, userID uint, n int) ([]domain.ChannelWithCounter, error)
	AddUserToChat(ctx context.Context, userID uint, chatID uint) (err error)
	RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) (err error)
	GetMessagesByChatID(ctx context.Context, chatID uint) []domain.Message
}

func GetChatByChatID(ctx context.Context, userID, chatID uint, chatStorage ChatStore, userStorage usecase.UserStore) (domain.Chat, error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chat, err := chatStorage.GetChatByChatID(ctx, chatID)
	if err != nil {
		return domain.Chat{}, err
	}
	belongs := CheckUserBelongsToChat(ctx, chatID, userID, chatStorage)
	if !belongs {
		logger.Info("GetChatByChatID: user does not belong", "userID", userID, "chatID", chatID)

		customErr := &domain.CustomError{
			Type:    "internal",
			Message: "user does not belong to chat",
			Segment: "method CheckUserBelongsToChat, chat_usecase.go",
		}
		return domain.Chat{}, customErr
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
		return chats[j].LastActionDateTime.Before(chats[i].LastActionDateTime)
	})
	return chats
}

func CheckUserBelongsToChat(ctx context.Context, chatID uint, userRequestingID uint, chatStorage ChatStore) bool {
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
func CreatePrivateChat(ctx context.Context, creatingUserID uint, companionID uint, chatStorage ChatStore, userStorage usecase.UserStore) (chatID uint, isNewChat bool, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	if creatingUserID == companionID {
		return 0, false, fmt.Errorf("Диалог с самим собой пока не поддерживается")
	}

	companion, found := userStorage.GetByUserID(ctx, companionID)
	fmt.Println("Comp: ", companion)
	if !found {
		logger.Error("CreatePrivateChat: user wasn't found", "companionID", companionID)
		return 0, false, fmt.Errorf("Пользователь, с которым вы хотите создать диалог, не найден")
	}
	exists, chatID, err := chatStorage.CheckPrivateChatExists(ctx, creatingUserID, companion.ID)
	if err != nil {
		return 0, false, err
	}
	if exists {
		return chatID, false, nil
	}
	chatID, err = chatStorage.CreateChat(ctx, companion.Username, "", creatingUserID, companion.ID)
	if err != nil {
		return 0, false, err
	}
	return chatID, true, nil
}

func DeleteChat(ctx context.Context, deletingUserID, chatID uint, chatStorage ChatStore) (wasDeleted bool, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("DeleteChat: enter", "userID", deletingUserID, "chatID", chatID)

	userBelongsToChat := CheckUserBelongsToChat(ctx, chatID, deletingUserID, chatStorage)
	if !userBelongsToChat {
		return false, fmt.Errorf("Неверный id для удаления")
	}
	chat, err := chatStorage.GetChatByChatID(ctx, chatID)
	if err != nil {
		return false, err
	}
	if (chat.Type == "3" || chat.Type == "2") && chat.CreatorID != deletingUserID {
		err := LeaveChat(ctx, deletingUserID, chatID, chatStorage)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	wasDeleted, err = chatStorage.DeleteChat(ctx, chatID)
	if err != nil {
		logger.Error("DeleteChat: error", "error", err.Error(), "wasDeleted", wasDeleted)
		return false, err
	}
	logger.Debug("DeleteChat: success", "wasDeleted", wasDeleted)
	return wasDeleted, err
}

func CreateGroupChat(ctx context.Context, creatingUserID uint, usersIDs []uint, chatName, description string, chatStorage ChatStore) (chatID uint, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	if len(usersIDs) < 3 {
		logger.Info("CreateGroupChat: len < 3", "users", usersIDs)
	}
	userMap := make(map[uint]bool)
	if usersIDs[0] != creatingUserID {

	}
	for i := range usersIDs {
		if userMap[usersIDs[i]] == true {
			logger.Info("user id is duplicated", "userID", usersIDs[i])
			break
		}
		userMap[usersIDs[i]] = true
	}
	usersIDs = append(usersIDs, creatingUserID)
	usersIDs[0], usersIDs[len(usersIDs)-1] = usersIDs[len(usersIDs)-1], usersIDs[0]

	chatID, err = chatStorage.CreateChat(ctx, chatName, description, usersIDs...)
	if err != nil {
		return 0, nil
	}
	return chatID, nil
}

func UpdateGroupChat(ctx context.Context, userID, chatID uint, name, desc *string, chatStorage ChatStore) (err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chat, err := chatStorage.GetChatByChatID(ctx, chatID)
	if chat.Type != "2" && chat.Type != "3" {
		return fmt.Errorf("internal error")
	}
	if err != nil {
		return fmt.Errorf("internal error")
	}
	userWasFound := false
	for i := range chat.Users {
		if chat.Users[i].UserID == userID {
			userWasFound = true
			break
		}
	}
	if !userWasFound {
		return fmt.Errorf("user does not belong to chat")
	}
	if name != nil {
		chat.Name = *name
	}
	if desc != nil {
		chat.Description = *desc
	}
	ok := chatStorage.UpdateGroupChat(ctx, chat)
	logger.Info("UpdateGroupChat", "ok", ok)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}

func GetMessagesByChatID(ctx context.Context, chatStorage ChatStore, chatID uint) []domain.Message {
	messages := chatStorage.GetMessagesByChatID(ctx, chatID)
	return messages
}

func GetPopularChannels(ctx context.Context, userID uint, chatStorage ChatStore) ([]domain.ChannelWithCounter, error) {
	channels, err := chatStorage.GetNPopularChannels(ctx, userID, 10)
	return channels, err
}

func JoinChannel(ctx context.Context, userID uint, channelID uint, chatStorage ChatStore) (err error) {
	channel, err := chatStorage.GetChatByChatID(ctx, channelID)
	if err != nil {
		return err
	}
	if channel.Type != "3" {
		return fmt.Errorf("Неверный id канала")
	}

	belongs := CheckUserBelongsToChat(ctx, channelID, userID, chatStorage)
	if belongs {
		return fmt.Errorf("Пользователь уже состоит в этом канале")
	}
	err = chatStorage.AddUserToChat(ctx, userID, channelID)
	if err != nil {
		return err
	}
	return nil
}

func LeaveChat(ctx context.Context, userID uint, channelID uint, chatStorage ChatStore) (err error) {
	channel, err := chatStorage.GetChatByChatID(ctx, channelID)
	if err != nil {
		return err
	}
	if channel.Type != "3" && channel.Type != "2" {
		return fmt.Errorf("Неверный id чата")
	}

	belongs := CheckUserBelongsToChat(ctx, channelID, userID, chatStorage)
	if !belongs {
		return fmt.Errorf("Пользователь не состоит в этом чате")
	}
	err = chatStorage.RemoveUserFromChat(ctx, userID, channelID)
	if err != nil {
		return err
	}
	return nil
}

func CreateChannel(ctx context.Context, creatingUserID uint, chatName, description string, chatStorage ChatStore) (chatID uint, err error) {
	chatID, err = chatStorage.CreateChat(ctx, chatName, description, creatingUserID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
