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
	IsNeededToShow bool `json:"isNeeded"`
}

type QuestionsResponse struct {
	Questions []domain.Question `json:"questions"`
}

type OneQuestionRequest struct {
	QuestionID int `json:"question_id"`
}

type GetQuestionsRequest struct {
	UserID uint `json:"user_id"`
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

func (FeedbackHandler *FeedbackHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	decoder := json.NewDecoder(r.Body)
	request := OneQuestionRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	isNeededToShow := usecase.IsReturnNeeded(ctx, FeedbackHandler.Feedback, userID, request.QuestionID)
	misc.WriteStatusJson(ctx, w, 200, GetFeedbackResponse{IsNeededToShow: isNeededToShow})
}

// GetQuestions получает список доступных опросов для пользователя
//
// @Summary
// @ID GetQuestions
// @Produce json
// @Success 200 {object}  domain.Response[QuestionsResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getQuestions [get]
func (FeedbackHandler *FeedbackHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	questions := usecase.ReturnQuestions(ctx, FeedbackHandler.Feedback, userID)
	misc.WriteStatusJson(ctx, w, 200, QuestionsResponse{Questions: questions})
}

func (FeedbackHandler *FeedbackHandler) Set(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	decoder := json.NewDecoder(r.Body)
	request := GetQuestionsRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	questions := usecase.ReturnQuestions(ctx, FeedbackHandler.Feedback, userID)
	misc.WriteStatusJson(ctx, w, 200, QuestionsResponse{Questions: questions})
}
