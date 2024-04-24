package delivery

import (
	"database/sql"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"

	//chatsInMemoryRepository "ProjectMessenger/internal/chats/repository/inMemory"
	repository "ProjectMessenger/internal/search/repository/db"
	"ProjectMessenger/internal/search/usecase"
)

type RequestChatIDBody struct {
	ChatID uint `json:"chatID"`
}

type SearchHandler struct {
	ChatsHandler *delivery.ChatsHandler
	Search       usecase.SearchStore
}

func NewSearchHandler(chatsHandler *delivery.ChatsHandler, database *sql.DB) *SearchHandler {
	return &SearchHandler{
		ChatsHandler: chatsHandler,
		Search:       repository.NewSearchStorage(database),
	}
}

// SendMessage method to send messages
//
// @Summary SendMessage
// @Description Сначала по этому URL надо произвести upgrade до вебсокета, потом слать json сообщений
// @ID sendMessage
// @Accept application/json
// @Produce application/json
// @Param user body  domain.Message true "message that was sent"
// @Success 200 {object}  domain.Response[int]
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error | could not upgrade connection"
// @Router /sendMessage [post]
func (SearchHandler *SearchHandler) SearchChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := SearchHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	upgrader := repository.UpgradeConnection()

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("SearchChats: upgrade failed", "err", err.Error())
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}
	user, found := SearchHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("could not upgrade connection :user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}
	SearchHandler.Search.HandleWebSocket(ctx, connection, user)
}
