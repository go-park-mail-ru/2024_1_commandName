package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/chats/repository/db"
	"ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
)

type ChatsHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Chats       usecase.ChatStore
}

type chatIDIsNewJsonResponse struct {
	ChatID    uint `json:"chat_id"`
	IsNewChat bool `json:"is_new_chat"`
}

type chatIDStruct struct {
	ChatID uint `json:"chat_id"`
}

type chatJsonResponse struct {
	Chat domain.Chat `json:"chat"`
}

type userIDJson struct {
	ID uint `json:"user_id"`
}

type createGroupJson struct {
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
	Users       []uint `json:"user_ids"`
}

type createChannelJson struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type deleteChatJsonResponse struct {
	SuccessfullyDeleted bool `json:"successfully_deleted"`
}

type updateChatJson struct {
	ChatID         uint    `json:"chat_id"`
	NewName        *string `json:"new_name"`
	NewDescription *string `json:"new_description"`
}

type getPopularChannelsResponse struct {
	Channels []domain.ChannelWithCounter `json:"channels"`
}

func NewChatsHandler(authHandler *authdelivery.AuthHandler, dataBase *sql.DB) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler: authHandler,
		Chats:       db.NewChatsStorage(dataBase),
	}
}

// GetChats gets Chats previews for user
//
// @Summary gets Chats previews for user
// @ID GetChats
// @Produce json
// @Success 200 {object}  domain.Response[domain.Chats]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChats [get]
func (chatsHandler ChatsHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	chats := usecase.GetChatsForUser(ctx, userID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	misc.WriteStatusJson(ctx, w, 200, domain.Chats{Chats: chats})
}

// GetChat gets one chat
//
// @Summary gets one chat
// @ID GetChat
// @Accept json
// @Produce json
// @Param user body  chatIDStruct true "id of chat to get"
// @Success 200 {object}  domain.Response[chatJsonResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChat [post]
func (chatsHandler ChatsHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := chatIDStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chat, err := usecase.GetChatByChatID(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, chatJsonResponse{Chat: chat})
}

// CreatePrivateChat creates dialogue
//
// @Summary creates dialogue
// @ID CreatePrivateChat
// @Accept json
// @Produce json
// @Param user body userIDJson true "ID of person to create private chat with"
// @Success 200 {object}  domain.Response[chatIDIsNewJsonResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized | Пользователь, с которым вы хотите создать дилаог, не найден | Чат с этим пользователем уже существует"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createPrivateChat [post]
func (chatsHandler ChatsHandler) CreatePrivateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	userIDFromRequest := userIDJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userIDFromRequest)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, isNewChat, err := usecase.CreatePrivateChat(ctx, userID, userIDFromRequest.ID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, chatIDIsNewJsonResponse{ChatID: chatID, IsNewChat: isNewChat})
}

// DeleteChat deletes chat
//
// @Summary deletes chat
// @ID DeleteChat
// @Accept json
// @Produce json
// @Param user body chatIDStruct true "ID of chat to delete"
// @Success 200 {object}  domain.Response[deleteChatJsonResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized | User doesn't belong to chat"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /deleteChat [post]
func (chatsHandler ChatsHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	chatIDJson := chatIDStruct{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chatIDJson)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	success, err := usecase.DeleteChat(ctx, userID, chatIDJson.ChatID, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, deleteChatJsonResponse{success})
}

// CreateGroupChat creates group chat
//
// @Summary creates group chat
// @ID CreateGroupChat
// @Accept json
// @Produce json
// @Param user body createGroupJson true "IDs of users to create group chat with"
// @Success 200 {object}  domain.Response[chatIDStruct]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized | Пользователь, с которым вы хотите создать дилаог, не найден | Чат с этим пользователем уже существует"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createGroupChat [post]
func (chatsHandler ChatsHandler) CreateGroupChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	groupRequest := createGroupJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&groupRequest)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, err := usecase.CreateGroupChat(ctx, userID, groupRequest.Users, groupRequest.GroupName, groupRequest.Description, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, chatIDStruct{ChatID: chatID})
}

// UpdateGroupChat updates group chat
//
// @Summary updates group chat
// @ID UpdateGroupChat
// @Accept json
// @Produce json
// @Param user body updateChatJson true "updated chat (если имя или описание не обновлялось, поле не слать вообще)"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /updateGroupChat [post]
func (chatsHandler ChatsHandler) UpdateGroupChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	updatedChat := updateChatJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedChat)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.UpdateGroupChat(ctx, userID, updatedChat.ChatID, updatedChat.NewName, updatedChat.NewDescription, chatsHandler.Chats)
	if err != nil {
		misc.WriteInternalErrorJson(ctx, w)
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}

// GetPopularChannels updates group chat
//
// @Summary gets 10 popular channels
// @ID GetPopularChannels
// @Produce json
// @Success 200 {object}  domain.Response[getPopularChannelsResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getPopularChannels [get]
func (chatsHandler ChatsHandler) GetPopularChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, _ := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	channels, err := usecase.GetPopularChannels(ctx, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "Что-то пошло не так"})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, channels)
}

// JoinChannel joins channel
//
// @Summary joins channel
// @ID JoinChannel
// @Accept json
// @Produce json
// @Param user body chatIDStruct true "id of channel"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /joinChannel [post]
func (chatsHandler ChatsHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := chatIDStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.JoinChannel(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, nil)
}

// LeaveChannel exits from channel
//
// @Summary exits from channel
// @ID LeaveChannel
// @Accept json
// @Produce json
// @Param user body chatIDStruct true "id of channel"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /leaveChannel [post]
func (chatsHandler ChatsHandler) LeaveChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := chatIDStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.LeaveChat(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, nil)
}

// CreateChannel creates channel
//
// @Summary creates channel
// @ID CreateChannel
// @Accept json
// @Produce json
// @Param user body createChannelJson true "IDs of users to create group chat with"
// @Success 200 {object}  domain.Response[chatIDStruct]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createChannel [post]
func (chatsHandler ChatsHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	channelRequest := createChannelJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&channelRequest)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, err := usecase.CreateChannel(ctx, userID, channelRequest.Name, channelRequest.Description, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	misc.WriteStatusJson(ctx, w, 200, chatIDStruct{ChatID: chatID})
}
