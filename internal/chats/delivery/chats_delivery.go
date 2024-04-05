package delivery

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	db "ProjectMessenger/internal/chats/repository/db"
	"ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
)

type ChatsHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Chats       usecase.ChatStore
}

type getChatStruct struct {
	ChatID uint `json:"chat_id"`
}

type chatJson struct {
	Chat domain.Chat `json:"chat"`
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
// @Param user body  getChatStruct true "id of chat to get"
// @Success 200 {object}  domain.Response[chatJson]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChat [post]
func (chatsHandler ChatsHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := getChatStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		http.Error(w, "wrong json structure", 400)
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

	misc.WriteStatusJson(ctx, w, 200, chatJson{Chat: chat})
}

// CreateDialogue creates dialogue
//
// @Summary creates dialogue
// @ID CreateDialogue
// @Accept json
// @Produce json
// @Param user body  domain.Person true "Person"
// @Success 200 {object}  domain.Response[domain.Chats]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createDialogue [get]
func (chatsHandler ChatsHandler) CreateDialogue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	chats := usecase.GetChatsForUser(ctx, userID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	misc.WriteStatusJson(ctx, w, 200, domain.Chats{Chats: chats})
}
