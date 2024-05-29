package delivery

import (
	"database/sql"
	"io/ioutil"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"
	repository "ProjectMessenger/internal/search/repository/db"
	"ProjectMessenger/internal/search/usecase"
	"github.com/mailru/easyjson"
)

type SearchHandler struct {
	ChatsHandler *delivery.ChatsHandler
	Search       usecase.SearchStore
}

func NewSearchHandler(chatsHandler *delivery.ChatsHandler, database *sql.DB) *SearchHandler {
	return &SearchHandler{
		ChatsHandler: chatsHandler,
		Search:       repository.NewSearchStorage(database, chatsHandler.Chats),
	}
}

func (SearchHandler *SearchHandler) SearchObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := SearchHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	searchRequestStruct := domain.SearchRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	err = easyjson.Unmarshal(body, &searchRequestStruct)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "json decode",
			Message: err.Error(),
			Segment: "SearchObjects, search_delivery.go",
		}
		logger.Error(customErr.Message, "err", customErr)
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: customErr.Message})
	}

	switch searchRequestStruct.Type {
	case "chat":
		foundChat := usecase.SearchChats(ctx, SearchHandler.Search, searchRequestStruct.Word, userID)
		misc.WriteStatusJson(ctx, w, 200, foundChat)
	case "contact":
		foundContact := usecase.SearchContacts(ctx, SearchHandler.Search, searchRequestStruct.Word, userID)
		misc.WriteStatusJson(ctx, w, 200, foundContact)
	case "channel":
		foundChannels := usecase.SearchChannels(ctx, SearchHandler.Search, searchRequestStruct.Word, userID)
		misc.WriteStatusJson(ctx, w, 200, foundChannels)
	case "message":
		foundMessages := usecase.SearchMessages(ctx, SearchHandler.Search, searchRequestStruct.Word, userID, searchRequestStruct.ChatID)
		misc.WriteStatusJson(ctx, w, 200, foundMessages)
	}
}
