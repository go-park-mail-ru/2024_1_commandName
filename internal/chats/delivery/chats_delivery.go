package delivery

import (
	"database/sql"
	"net/http"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	db "ProjectMessenger/internal/chats/repository/db"
	chatsInMemoryRepository "ProjectMessenger/internal/chats/repository/inMemory"
	"ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
)

type ChatsHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Chats       usecase.ChatStore
}

func NewChatsHandler(authHandler *authdelivery.AuthHandler, dataBase *sql.DB) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler: authHandler,
		Chats:       db.NewChatsStorage(dataBase),
	}
}

func NewChatsHandlerMemory(authHandler *authdelivery.AuthHandler) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler: authHandler,
		Chats:       chatsInMemoryRepository.NewChatsStorage(),
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

	chats := usecase.GetChatsForUser(ctx, userID, chatsHandler.AuthHandler.Chats)
	misc.WriteStatusJson(w, 200, domain.Chats{Chats: chats})
}
