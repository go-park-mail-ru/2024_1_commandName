package delivery

import (
	"database/sql"
	"encoding/json"
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

type deleteChatJsonResponse struct {
	SuccessfullyDeleted bool `json:"successfully_deleted"`
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
	}
	userBelongsToChat := usecase.CheckUserBelongsToChat(ctx, chatIDJson.ChatID, userID, chatsHandler.Chats)
	if !userBelongsToChat {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "User doesn't belong to chat"})
		return
	}
	success, err := usecase.DeletePrivateChat(ctx, chatIDJson.ChatID, userID, chatsHandler.Chats)
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
