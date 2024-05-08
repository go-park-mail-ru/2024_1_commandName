package usecase

import (
	"ProjectMessenger/domain"
	chats "ProjectMessenger/internal/chats_service/proto"
	"context"
	"fmt"
	"sort"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type ChatManager struct {
	chats.UnimplementedChatServiceServer
	storage ChatStore
}

func NewChatManager(storage ChatStore) *ChatManager {
	return &ChatManager{storage: storage}
}

func convertChat(chat domain.Chat) *chats.Chat {
	messagesGRPC := make([]*chats.Message, 0)
	for i := range chat.Messages {
		messagesGRPC = append(messagesGRPC, &chats.Message{
			Id:          uint64(chat.Messages[i].ID),
			ChatId:      uint64(chat.Messages[i].ChatID),
			UserId:      uint64(chat.Messages[i].UserID),
			MessageText: chat.Messages[i].Message,
			Edited:      chat.Messages[i].Edited,
			EditedAt:    timestamppb.New(chat.Messages[i].EditedAt),
			SentAt:      timestamppb.New(chat.Messages[i].CreatedAt),
			Username:    chat.Messages[i].SenderUsername,
		})
	}
	usersGRPC := make([]*chats.ChatUser, 0)
	for i := range chat.Users {
		usersGRPC = append(usersGRPC, &chats.ChatUser{
			ChatId: uint64(chat.Users[i].ChatID),
			UserId: uint64(chat.Users[i].UserID),
		})
	}
	lastMessage := &chats.Message{
		Id:          uint64(chat.LastMessage.ID),
		ChatId:      uint64(chat.LastMessage.ChatID),
		UserId:      uint64(chat.LastMessage.UserID),
		MessageText: chat.LastMessage.Message,
		Edited:      chat.LastMessage.Edited,
		EditedAt:    timestamppb.New(chat.LastMessage.EditedAt),
		SentAt:      timestamppb.New(chat.LastMessage.CreatedAt),
		Username:    chat.LastMessage.SenderUsername,
	}
	return &chats.Chat{
		Id:                 uint64(chat.ID),
		Type:               chat.Type,
		Name:               chat.Name,
		Description:        chat.Description,
		AvatarPath:         chat.AvatarPath,
		CreatorId:          uint64(chat.CreatorID),
		Messages:           messagesGRPC,
		Users:              usersGRPC,
		CreatedAt:          timestamppb.New(chat.CreatedAt),
		EditedAt:           timestamppb.New(chat.EditedAt),
		LastActionDateTime: timestamppb.New(chat.LastActionDateTime),
		LastMessage:        lastMessage,
		LastSeenMessageId:  0,
	}
}

func (cm *ChatManager) checkUserBelongsToChat(ctx context.Context, chatID uint, userRequestingID uint) bool {
	//logger := slog.With("requestID", ctx.Value("traceID"))
	usersOfChat := cm.storage.GetChatUsersByChatID(ctx, chatID)
	for i := range usersOfChat {
		if usersOfChat[i].UserID == userRequestingID {
			//logger.Debug("CheckUserBelongsToChat: true", "chatID", chatID, "userID", userRequestingID)
			return true
		}
	}
	//logger.Debug("CheckUserBelongsToChat: false", "chatID", chatID, "userID", userRequestingID)
	return false
}

func (cm *ChatManager) leaveChat(ctx context.Context, userID uint, channelID uint) (err error) {
	channel, err := cm.storage.GetChatByChatID(ctx, channelID)
	if err != nil {
		return err
	}
	if channel.Type != "3" && channel.Type != "2" {
		return fmt.Errorf("Неверный id чата")
	}

	belongs := cm.checkUserBelongsToChat(ctx, channelID, userID)
	if !belongs {
		return fmt.Errorf("Пользователь не состоит в этом чате")
	}
	err = cm.storage.RemoveUserFromChat(ctx, userID, channelID)
	if err != nil {
		return err
	}
	return nil
}

func (cm *ChatManager) GetChatByChatID(ctx context.Context, in *chats.UserAndChatID) (*chats.Chat, error) {
	chatID := uint(in.GetChatID())
	userID := uint(in.GetUserID())
	chat, err := cm.storage.GetChatByChatID(ctx, chatID)
	if err != nil {
		return &chats.Chat{}, err
	}
	belongs := cm.checkUserBelongsToChat(ctx, chatID, userID)
	if !belongs {
		//return &chats.Chat{}, status.Error(500, "")
	}
	return convertChat(chat), nil
}

func (cm *ChatManager) GetChatsForUser(ctx context.Context, in *chats.UserID) (*chats.ChatArray, error) {
	userID := uint(in.GetUserID())
	chatsForUser := cm.storage.GetChatsForUser(ctx, userID)

	sort.Slice(chatsForUser, func(i, j int) bool {
		return chatsForUser[j].LastActionDateTime.Before(chatsForUser[i].LastActionDateTime)
	})
	chatsGRPC := &chats.ChatArray{Chats: make([]*chats.Chat, 0)}
	for i := range chatsForUser {
		chatsGRPC.Chats = append(chatsGRPC.Chats, convertChat(chatsForUser[i]))
	}
	return chatsGRPC, nil
}

func (cm *ChatManager) CheckUserBelongsToChat(ctx context.Context, in *chats.UserAndChatID) (*chats.BoolResponse, error) {
	belongs := cm.checkUserBelongsToChat(ctx, uint(in.GetChatID()), uint(in.GetUserID()))
	return &chats.BoolResponse{Res: belongs}, nil
}

func (cm *ChatManager) CreatePrivateChat(ctx context.Context, in *chats.TwoUserIDs) (*chats.CreateChatResponse, error) {
	creatingUserID := uint(in.GetID1())
	companionID := uint(in.GetID2())
	if creatingUserID == companionID {
		return &chats.CreateChatResponse{}, status.Error(400, "Диалог с самим собой пока не поддерживается")
	}

	exists, chatID, err := cm.storage.CheckPrivateChatExists(ctx, creatingUserID, companionID)
	if err != nil {
		return &chats.CreateChatResponse{}, status.Error(500, "")
	}
	if exists {
		return &chats.CreateChatResponse{
			ChatID:    uint64(chatID),
			IsNewChat: false,
		}, nil
	}
	chatID, err = cm.storage.CreateChat(ctx, "", "", creatingUserID, companionID)
	if err != nil {
		return &chats.CreateChatResponse{}, status.Error(500, "")
	}
	return &chats.CreateChatResponse{
		ChatID:    uint64(chatID),
		IsNewChat: true,
	}, nil
}

func (cm *ChatManager) DeleteChat(ctx context.Context, in *chats.UserAndChatID) (*chats.BoolResponse, error) {
	chatID := uint(in.GetChatID())
	deletingUserID := uint(in.GetUserID())
	userBelongsToChat := cm.checkUserBelongsToChat(ctx, chatID, deletingUserID)
	if !userBelongsToChat {
		return &chats.BoolResponse{}, status.Error(400, "Неверный id для удаления")
	}
	chat, err := cm.storage.GetChatByChatID(ctx, chatID)
	if err != nil {
		return &chats.BoolResponse{}, status.Error(500, "")
	}
	if (chat.Type == "3" || chat.Type == "2") && chat.CreatorID != deletingUserID {
		err := cm.leaveChat(ctx, deletingUserID, chatID)
		if err != nil {
			return &chats.BoolResponse{}, status.Error(500, "")
		}
		return &chats.BoolResponse{Res: true}, nil
	}
	wasDeleted, err := cm.storage.DeleteChat(ctx, chatID)
	if err != nil {
		//logger.Error("DeleteChat: error", "error", err.Error(), "wasDeleted", wasDeleted)
		return &chats.BoolResponse{}, status.Error(500, "")
	}
	//logger.Debug("DeleteChat: success", "wasDeleted", wasDeleted)
	return &chats.BoolResponse{Res: wasDeleted}, nil
}

func (cm *ChatManager) CreateGroupChat(ctx context.Context, in *chats.CreateGroupReq) (*chats.CreateChatResponse, error) {
	creatingUserID := uint(in.GetCreatingUserID())
	chatName := in.GetName()
	description := in.GetDescription()
	usersIDs := make([]uint, 0)
	grpcUsers := in.GetUsers().GetUsers()
	for i := range grpcUsers {
		usersIDs = append(usersIDs, uint(grpcUsers[i].GetUserID()))
	}

	if len(usersIDs) < 3 {
		//logger.Info("CreateGroupChat: len < 3", "users", usersIDs)
	}
	userMap := make(map[uint]bool)
	if usersIDs[0] != creatingUserID {

	}
	for i := range usersIDs {
		if userMap[usersIDs[i]] == true {
			//logger.Info("user id is duplicated", "userID", usersIDs[i])
			break
		}
		userMap[usersIDs[i]] = true
	}
	usersIDs = append(usersIDs, creatingUserID)
	usersIDs[0], usersIDs[len(usersIDs)-1] = usersIDs[len(usersIDs)-1], usersIDs[0]

	chatID, err := cm.storage.CreateChat(ctx, chatName, description, usersIDs...)
	if err != nil {
		return &chats.CreateChatResponse{}, status.Error(500, "")
	}
	return &chats.CreateChatResponse{ChatID: uint64(chatID)}, nil
}

func (cm *ChatManager) UpdateGroupChat(ctx context.Context, in *chats.UpdateGroupChatReq) (*chats.Empty, error) {
	chatID := uint(in.GetChatID())
	userID := uint(in.GetUserID())
	name := in.GetName()
	desc := in.GetDescription()
	chat, err := cm.storage.GetChatByChatID(ctx, chatID)
	if chat.Type != "2" && chat.Type != "3" {
		return &chats.Empty{}, status.Error(400, "")
	}
	if err != nil {
		return &chats.Empty{}, status.Error(500, "")
	}
	userWasFound := false
	for i := range chat.Users {
		if chat.Users[i].UserID == userID {
			userWasFound = true
			break
		}
	}
	if !userWasFound {
		return &chats.Empty{}, status.Error(400, "user does not belong to chat")
	}
	if name != "" {
		chat.Name = name
	}
	if desc != "" {
		chat.Description = desc
	}
	ok := cm.storage.UpdateGroupChat(ctx, chat)
	//logger.Info("UpdateGroupChat", "ok", ok)
	if !ok {
		return &chats.Empty{}, status.Error(500, "")
	}
	return &chats.Empty{Dummy: true}, nil
}

func (cm *ChatManager) GetMessagesByChatID(ctx context.Context, in *chats.ChatID) (*chats.MessageArray, error) {
	chatID := uint(in.GetChatID())
	messages := cm.storage.GetMessagesByChatID(ctx, chatID)
	messagesGRPC := make([]*chats.Message, 0)
	for i := range messages {
		messagesGRPC = append(messagesGRPC, &chats.Message{
			Id:          uint64(messages[i].ID),
			ChatId:      uint64(messages[i].ChatID),
			UserId:      uint64(messages[i].UserID),
			MessageText: messages[i].Message,
			Edited:      messages[i].Edited,
			EditedAt:    timestamppb.New(messages[i].EditedAt),
			SentAt:      timestamppb.New(messages[i].CreatedAt),
			Username:    messages[i].SenderUsername,
		})
	}
	return &chats.MessageArray{Messages: messagesGRPC}, nil
}

func (cm *ChatManager) GetPopularChannels(ctx context.Context, in *chats.UserID) (*chats.ChannelsArray, error) {
	userID := uint(in.GetUserID())
	channels, err := cm.storage.GetNPopularChannels(ctx, userID, 10)
	if err != nil {
		return &chats.ChannelsArray{}, status.Error(500, "")
	}
	channelsGRPC := make([]*chats.ChannelWithCounter, 0)
	for i := range channels {
		channelsGRPC = append(channelsGRPC, &chats.ChannelWithCounter{
			Id:          uint64(channels[i].ID),
			Name:        channels[i].Name,
			Description: channels[i].Description,
			CreatorId:   uint32(channels[i].CreatorID),
			Avatar:      channels[i].Avatar,
			IsMember:    channels[i].IsMember,
			NumOfUsers:  int32(channels[i].NumOfUsers),
		})
	}
	return &chats.ChannelsArray{Channels: channelsGRPC}, nil
}

func (cm *ChatManager) JoinChannel(ctx context.Context, in *chats.UserAndChatID) (*chats.Empty, error) {
	channelID := uint(in.GetChatID())
	userID := uint(in.GetUserID())
	channel, err := cm.storage.GetChatByChatID(ctx, channelID)
	if err != nil {
		return &chats.Empty{}, status.Error(500, "")
	}
	if channel.Type != "3" {
		return &chats.Empty{}, status.Error(400, "Неверный id канала")
	}

	belongs := cm.checkUserBelongsToChat(ctx, channelID, userID)
	if belongs {
		return &chats.Empty{}, status.Error(400, "Пользователь уже состоит в этом канале")
	}
	err = cm.storage.AddUserToChat(ctx, userID, channelID)
	if err != nil {
		return &chats.Empty{}, status.Error(500, "")
	}
	return &chats.Empty{Dummy: true}, nil
}

func (cm *ChatManager) LeaveChat(ctx context.Context, in *chats.UserAndChatID) (*chats.Empty, error) {
	userID := uint(in.GetUserID())
	channelID := uint(in.GetChatID())
	channel, err := cm.storage.GetChatByChatID(ctx, channelID)
	if err != nil {
		return &chats.Empty{}, status.Error(500, "")
	}
	if channel.Type != "3" && channel.Type != "2" {
		return &chats.Empty{}, status.Error(400, "Неверный id чата")
	}

	belongs := cm.checkUserBelongsToChat(ctx, channelID, userID)
	if !belongs {
		return &chats.Empty{}, status.Error(400, "Пользователь не состоит в этом чате")
	}
	err = cm.storage.RemoveUserFromChat(ctx, userID, channelID)
	if err != nil {
		return &chats.Empty{}, status.Error(500, "")
	}
	return &chats.Empty{Dummy: true}, nil
}

func (cm *ChatManager) CreateChannel(ctx context.Context, in *chats.CreateChannelReq) (*chats.ChatID, error) {
	creatingUserID := uint(in.GetUserID())
	chatName := in.GetName()
	description := in.GetDescription()
	chatID, err := cm.storage.CreateChat(ctx, chatName, description, creatingUserID)
	if err != nil {
		return &chats.ChatID{}, status.Error(500, "")
	}
	return &chats.ChatID{ChatID: uint64(chatID)}, nil
}
