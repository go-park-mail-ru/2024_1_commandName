package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/usecase"
	chats "ProjectMessenger/internal/chats_service/proto"

	"google.golang.org/grpc/status"
)

func convertChat(chat *chats.Chat) domain.Chat {
	messages := make([]*domain.Message, 0)
	users := make([]*domain.ChatUser, 0)
	if chat != nil {

		for i := range chat.Messages {
			messages = append(messages, &domain.Message{
				ID:             uint(chat.Messages[i].GetId()),
				ChatID:         uint(chat.Messages[i].GetChatId()),
				UserID:         uint(chat.Messages[i].GetUserId()),
				Message:        chat.Messages[i].GetMessageText(),
				Edited:         chat.Messages[i].GetEdited(),
				EditedAt:       chat.Messages[i].EditedAt.AsTime(),
				CreatedAt:      chat.Messages[i].SentAt.AsTime(),
				SenderUsername: chat.Messages[i].Username,
			})
		}
		for i := range chat.Users {
			users = append(users, &domain.ChatUser{
				ChatID: int(chat.Users[i].ChatId),
				UserID: uint(chat.Users[i].UserId),
			})
		}
	}
	lastMessage := domain.Message{}
	if chat != nil {
		lastMessage = domain.Message{
			ID:             uint(chat.LastMessage.GetId()),
			ChatID:         uint(chat.LastMessage.GetChatId()),
			UserID:         uint(chat.LastMessage.GetUserId()),
			Message:        chat.LastMessage.GetMessageText(),
			Edited:         chat.LastMessage.GetEdited(),
			EditedAt:       chat.LastMessage.EditedAt.AsTime(),
			CreatedAt:      chat.LastMessage.SentAt.AsTime(),
			SenderUsername: chat.LastMessage.Username,
		}
	}
	return domain.Chat{
		ID:                 uint(chat.GetId()),
		Type:               chat.GetType(),
		Name:               chat.GetName(),
		Description:        chat.GetDescription(),
		AvatarPath:         chat.GetAvatarPath(),
		CreatorID:          uint(chat.GetCreatorId()),
		Messages:           messages,
		Users:              users,
		CreatedAt:          chat.CreatedAt.AsTime(),
		EditedAt:           chat.EditedAt.AsTime(),
		LastActionDateTime: chat.LastActionDateTime.AsTime(),
		LastMessage:        lastMessage,
		LastSeenMessageID:  int(chat.GetLastSeenMessageId()),
	}
}

func GetChatByChatID(ctx context.Context, userID, chatID uint, userStorage usecase.UserStore, chatsGRPC chats.ChatServiceClient) (domain.Chat, error) {
	chatGRPC, err := chatsGRPC.GetChatByChatID(ctx, &chats.UserAndChatID{UserID: uint64(userID), ChatID: uint64(chatID)})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case 500:
				return domain.Chat{}, fmt.Errorf("internal error")
			case 400:
				return domain.Chat{}, fmt.Errorf(e.Message())
			}
		}
		fmt.Println(err)
		return domain.Chat{}, err
	}
	chat := convertChat(chatGRPC)

	if chat.Type == "1" {
		name, _ := GetCompanionNameForPrivateChat(ctx, chat, userID, userStorage)
		chat.Name = name
	}
	return chat, nil
}

func GetChatsForUser(ctx context.Context, userID uint, chatsGRPC chats.ChatServiceClient, userStorage usecase.UserStore) []domain.Chat {
	//logger := slog.With("requestID", ctx.Value("traceID"))
	chatsResp, err := chatsGRPC.GetChatsForUser(ctx, &chats.UserID{UserID: uint64(userID)})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	chatsRes := make([]domain.Chat, 0)
	for i := range chatsResp.Chats {
		current := convertChat(chatsResp.Chats[i])
		if chatsResp.Chats[i].Type == "1" {
			name, ok := GetCompanionNameForPrivateChat(ctx, current, userID, userStorage)
			if !ok {
				//logger.Debug("GetChatsForUser: getchatname failed", "userID", userID, "chatID", chats[i].ID)
				continue
			}
			current.Name = name
		}
		fmt.Println(current)
		chatsRes = append(chatsRes, current)
	}

	sort.Slice(chatsRes, func(i, j int) bool {
		return chatsRes[j].LastActionDateTime.Before(chatsRes[i].LastActionDateTime)
	})
	return chatsRes
}

func CheckUserBelongsToChat(ctx context.Context, chatID uint, userRequestingID uint, chatsGRPC chats.ChatServiceClient) bool {
	logger := slog.With("requestID", ctx.Value("traceID"))
	belongsGRPC, err := chatsGRPC.CheckUserBelongsToChat(ctx, &chats.UserAndChatID{
		UserID: uint64(userRequestingID),
		ChatID: uint64(chatID),
	})
	if err != nil {
		return false
	}
	logger.Debug("CheckUserBelongsToChat: false", "chatID", chatID, "userID", userRequestingID)
	return belongsGRPC.Res
}

func GetCompanionNameForPrivateChat(ctx context.Context, chat domain.Chat, userRequestingID uint, userStorage usecase.UserStore) (companionUsername string, belongs bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))

	usersOfChat := chat.Users
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
func CreatePrivateChat(ctx context.Context, creatingUserID uint, companionID uint, chatsGRPC chats.ChatServiceClient, userStorage usecase.UserStore) (chatID uint, isNewChat bool, err error) {
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
	resp, err := chatsGRPC.CreatePrivateChat(ctx, &chats.TwoUserIDs{
		ID1: uint64(creatingUserID),
		ID2: uint64(companionID),
	})
	if err != nil {
		return 0, false, err
	}
	return uint(resp.GetChatID()), resp.GetIsNewChat(), nil
}

func DeleteChat(ctx context.Context, deletingUserID, chatID uint, chatsGRPC chats.ChatServiceClient) (wasDeleted bool, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("DeleteChat: enter", "userID", deletingUserID, "chatID", chatID)
	success, err := chatsGRPC.DeleteChat(ctx, &chats.UserAndChatID{
		UserID: uint64(deletingUserID),
		ChatID: uint64(chatID),
	})
	if err != nil {
		return false, err
	}
	return success.Res, err
}

func CreateGroupChat(ctx context.Context, creatingUserID uint, usersIDs []uint, chatName, description string, chatsGRPC chats.ChatServiceClient) (chatID uint, err error) {
	//logger := slog.With("requestID", ctx.Value("traceID"))
	usersGRPC := make([]*chats.UserID, 0)
	for i := range usersIDs {
		usersGRPC = append(usersGRPC, &chats.UserID{UserID: uint64(usersIDs[i])})
	}

	resp, err := chatsGRPC.CreateGroupChat(ctx, &chats.CreateGroupReq{
		CreatingUserID: uint64(creatingUserID),
		Users:          &chats.CreateGroupReq_UserArray{Users: usersGRPC},
		Name:           chatName,
		Description:    description,
	})
	if err != nil {
		return 0, err
	}
	return uint(resp.GetChatID()), nil
}

func UpdateGroupChat(ctx context.Context, userID, chatID uint, name, desc *string, chatsGRPC chats.ChatServiceClient) (err error) {
	//logger := slog.With("requestID", ctx.Value("traceID"))
	_, err = chatsGRPC.UpdateGroupChat(ctx, &chats.UpdateGroupChatReq{
		UserID:      uint64(userID),
		ChatID:      uint64(chatID),
		Name:        *name,
		Description: *desc,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetMessagesByChatID(ctx context.Context, chatsGRPC chats.ChatServiceClient, chatID uint) []domain.Message {
	resp, err := chatsGRPC.GetMessagesByChatID(ctx, &chats.ChatID{ChatID: uint64(chatID)})
	if err != nil {
		return nil
	}
	messages := make([]domain.Message, 0)
	for i := range resp.Messages {
		messages = append(messages, domain.Message{
			ID:             uint(resp.Messages[i].GetId()),
			ChatID:         uint(resp.Messages[i].GetChatId()),
			UserID:         uint(resp.Messages[i].GetUserId()),
			Message:        resp.Messages[i].GetMessageText(),
			Edited:         resp.Messages[i].GetEdited(),
			EditedAt:       resp.Messages[i].GetEditedAt().AsTime(),
			CreatedAt:      resp.Messages[i].GetSentAt().AsTime(),
			SenderUsername: resp.Messages[i].GetUsername(),
		})
	}
	return messages
}

func GetPopularChannels(ctx context.Context, userID uint, chatsGRPC chats.ChatServiceClient) ([]domain.ChannelWithCounter, error) {
	resp, err := chatsGRPC.GetPopularChannels(ctx, &chats.UserID{UserID: uint64(userID)})
	if err != nil {
		return nil, err
	}
	channels := make([]domain.ChannelWithCounter, 0)
	for i := range resp.Channels {
		channels = append(channels, domain.ChannelWithCounter{
			ID:          uint(resp.Channels[i].GetId()),
			Name:        resp.Channels[i].GetName(),
			Description: resp.Channels[i].GetDescription(),
			CreatorID:   uint(resp.Channels[i].GetCreatorId()),
			Avatar:      resp.Channels[i].GetAvatar(),
			IsMember:    resp.Channels[i].GetIsMember(),
			NumOfUsers:  int(resp.Channels[i].GetNumOfUsers()),
		})
	}
	return channels, nil
}

func JoinChannel(ctx context.Context, userID uint, channelID uint, chatsGRPC chats.ChatServiceClient) (err error) {
	_, err = chatsGRPC.JoinChannel(ctx, &chats.UserAndChatID{
		UserID: uint64(userID),
		ChatID: uint64(channelID),
	})
	if err != nil {
		return err
	}
	return nil
}

func LeaveChat(ctx context.Context, userID uint, channelID uint, chatsGRPC chats.ChatServiceClient) (err error) {
	_, err = chatsGRPC.LeaveChat(ctx, &chats.UserAndChatID{
		UserID: uint64(userID),
		ChatID: uint64(channelID),
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateChannel(ctx context.Context, creatingUserID uint, chatName, description string, chatsGRPC chats.ChatServiceClient) (chatID uint, err error) {
	channel, err := chatsGRPC.CreateChannel(ctx, &chats.CreateChannelReq{
		UserID:      uint64(creatingUserID),
		Name:        chatName,
		Description: description,
	})
	if err != nil {
		return 0, err
	}
	return uint(channel.GetChatID()), nil
}
