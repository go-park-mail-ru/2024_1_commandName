package delivery

import (
	"errors"
	"net/http"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	authusecase "ProjectMessenger/internal/auth/usecase"
	"ProjectMessenger/internal/chats/repository/inMemory"
	"ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
)

type ChatsHandler struct {
	AuthHandler *authdelivery.AuthHandler
	Chats       usecase.ChatStore
}

func NewChatsHandler(authHandler *authdelivery.AuthHandler) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler: authHandler,
		Chats:       inMemory.NewChatsStorage(),
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
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "Person not authorized"})
		return
	}

	authorized, userID := authusecase.CheckAuthorized(session.Value, chatsHandler.AuthHandler.Sessions)

	if !authorized {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "Person not authorized"})
		return
	}

	chats := usecase.GetChatsForUser(userID, chatsHandler.AuthHandler.Chats)
	misc.WriteStatusJson(w, 200, domain.Chats{Chats: chats})
}
