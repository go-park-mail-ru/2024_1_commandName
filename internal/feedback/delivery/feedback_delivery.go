package delivery

import (
	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"
	"database/sql"
	"encoding/json"
	"net/http"

	feedbackrepo "ProjectMessenger/internal/feedback/repository/db"
	usecase "ProjectMessenger/internal/feedback/usecase"
)

type GetFeedbackResponse struct {
	IsNeededToShow bool `json:"isNeeded"`
}

type getQuestionsRequest struct {
	QuestionID int `json:"question_id"`
	Grade      int `json:"grade"`
}

type QuestionsResponse struct {
	Questions []domain.Question `json:"questions"`
}

type OneQuestionRequest struct {
	QuestionID int `json:"question_id"`
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

func (FeedbackHandler *FeedbackHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	questions := usecase.ReturnQuestions(ctx, FeedbackHandler.Feedback, userID)
	misc.WriteStatusJson(ctx, w, 200, QuestionsResponse{Questions: questions})
}

// SetAnswer получает список доступных опросов для пользователя
//
// @Summary
// @ID SetAnswer
// @Accept json
// @Produce json
// @Param user body  getQuestionsRequest true "user answer to question"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /setAnswer [get]
func (FeedbackHandler *FeedbackHandler) SetAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	decoder := json.NewDecoder(r.Body)
	userAnswer := getQuestionsRequest{}
	err := decoder.Decode(&userAnswer)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	usecase.SetAnswer(ctx, userID, userAnswer.QuestionID, userAnswer.Grade, FeedbackHandler.Feedback)
	misc.WriteStatusJson(ctx, w, 200, nil)
}

func (FeedbackHandler *FeedbackHandler) GetStatistic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, _ := FeedbackHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	statistic := usecase.GetAllStatistic(ctx, FeedbackHandler.Feedback)

	misc.WriteStatusJson(ctx, w, 200, statistic)
}
