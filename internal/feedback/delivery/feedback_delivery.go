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

type QuestionsResponse struct {
	Questions []domain.Question `json:"Questions"`
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

	isAnswered := usecase.IsReturnNeeded(ctx, FeedbackHandler.Feedback, userID, request.QuestionID)
	misc.WriteStatusJson(ctx, w, 200, GetFeedbackResponse{IsAnswered: isAnswered})
}

func (FeedbackHandler *FeedbackHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
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
