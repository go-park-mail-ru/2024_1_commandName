package delivery

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"

	feedbackrepo "ProjectMessenger/internal/feedback/repository/db"
	usecase "ProjectMessenger/internal/feedback/usecase"
)

type GetFeedbackResponse struct {
	IsAnswered bool `json:"isAnswered"`
}

type FeedbackHandler struct {
	Feedback     *feedbackrepo.Feedback
	ChatsHandler *delivery.ChatsHandler
}

func NewMessagesHandler(chatsHandler *delivery.ChatsHandler, database *sql.DB) *FeedbackHandler {
	return &FeedbackHandler{
		Feedback:     feedbackrepo.NewFeedbackStorage(database),
		ChatsHandler: chatsHandler,
	}
}

func (FeedbackHandler *FeedbackHandler) CheckGlobalUserExpNeeded(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	isAnswered := usecase.IsUserAnsweredOnGlobalCSAT(ctx, FeedbackHandler.Feedback, userID)
	misc.WriteStatusJson(ctx, w, 200, GetFeedbackResponse{IsAnswered: isAnswered})
}

func (FeedbackHandler *FeedbackHandler) CheckFitchUserExpNeeded(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	isAnswered := usecase.IsUserAnsweredOnGlobalCSAT(ctx, FeedbackHandler.Feedback, userID)
	misc.WriteStatusJson(ctx, w, 200, GetFeedbackResponse{IsAnswered: isAnswered})
}

// GetChatMessages returns messages of some chat
//
// @Summary GetChatMessages
// @ID getChatMessages
// @Accept application/json
// @Produce application/json
// @Param user body  RequestChatIDBody true "ID of chat"
// @Success 200 {object}  domain.Response[domain.Messages]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChatMessages [post]
func (messageHandler *MessageHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, _ := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var RequestChatID RequestChatIDBody
	err := decoder.Decode(&RequestChatID)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	limit := 100
	messages := usecase.GetChatMessages(r.Context(), limit, RequestChatID.ChatID, messageHandler.Messages)
	misc.WriteStatusJson(ctx, w, 200, domain.Messages{Messages: messages})
}
