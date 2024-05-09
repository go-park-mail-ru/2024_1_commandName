package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/chats/repository/db"
	"ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
	"github.com/prometheus/client_golang/prometheus"
)

type ChatsHandler struct {
	AuthHandler       *authdelivery.AuthHandler
	Chats             usecase.ChatStore
	prometheusMetrics *PrometheusMetrics
}

type chatIDIsNewJsonResponse struct {
	ChatID    uint `json:"chat_id"`
	IsNewChat bool `json:"is_new_chat"`
}

type messagesByChatIDRequest struct {
	ChatID uint `json:"chat_id"`
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

type createGroupJson struct {
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
	Users       []uint `json:"user_ids"`
}

type createChannelJson struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type deleteChatJsonResponse struct {
	SuccessfullyDeleted bool `json:"successfully_deleted"`
}

type updateChatJson struct {
	ChatID         uint    `json:"chat_id"`
	NewName        *string `json:"new_name"`
	NewDescription *string `json:"new_description"`
}

type getPopularChannelsResponse struct {
	Channels []domain.ChannelWithCounter `json:"channels"`
}

type PrometheusMetrics struct {
	ActiveSessionsCount prometheus.Gauge
	Hits                *prometheus.CounterVec
	Errors              *prometheus.CounterVec
	Methods             *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	chats_hits := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chats_hits",
			Help: "Total number of chats hits.",
		}, []string{"status", "path"},
	)

	chats_errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chats_errors",
			Help: "Number of errors some type.",
		}, []string{"error_type"},
	)

	chats_methods := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chats_called_methods",
			Help: "Number of called methods.",
		}, []string{"method"},
	)

	chats_requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "chats_http_request_duration_seconds",
			Help:    "Histogram of request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	prometheus.MustRegister(chats_hits, chats_errors, chats_methods, chats_requestDuration)

	return &PrometheusMetrics{
		Hits:            chats_hits,
		Errors:          chats_errors,
		Methods:         chats_methods,
		requestDuration: chats_requestDuration,
	}
}

func NewChatsHandler(authHandler *authdelivery.AuthHandler, dataBase *sql.DB) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler:       authHandler,
		Chats:             db.NewChatsStorage(dataBase),
		prometheusMetrics: NewPrometheusMetrics(),
	}
}

func NewRawChatsHandler(authHandler *authdelivery.AuthHandler, dataBase *sql.DB) *ChatsHandler {
	return &ChatsHandler{
		AuthHandler: authHandler,
		Chats:       db.NewRawChatsStorage(dataBase),
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
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("GetChats").Inc()
	start := time.Now()
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	chats := usecase.GetChatsForUser(ctx, userID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	misc.WriteStatusJson(ctx, w, 200, domain.Chats{Chats: chats})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/getChats").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
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
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("GetChat").Inc()
	start := time.Now()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		chatsHandler.prometheusMetrics.Hits.WithLabelValues("405", r.URL.String()).Inc()
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
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		chatsHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chat, err := usecase.GetChatByChatID(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	if err != nil {

		if err.Error() == "internal error" {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, chatJsonResponse{Chat: chat})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/getChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
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
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("CreatePrivateChat").Inc()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
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
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, isNewChat, err := usecase.CreatePrivateChat(ctx, userID, userIDFromRequest.ID, chatsHandler.Chats, chatsHandler.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, chatIDIsNewJsonResponse{ChatID: chatID, IsNewChat: isNewChat})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/createPrivateChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
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
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("DeleteChat").Inc()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
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
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	success, err := usecase.DeleteChat(ctx, userID, chatIDJson.ChatID, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, deleteChatJsonResponse{success})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/createPrivateChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// CreateGroupChat creates group chat
//
// @Summary creates group chat
// @ID CreateGroupChat
// @Accept json
// @Produce json
// @Param user body createGroupJson true "IDs of users to create group chat with"
// @Success 200 {object}  domain.Response[chatIDStruct]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized | Пользователь, с которым вы хотите создать дилаог, не найден | Чат с этим пользователем уже существует"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createGroupChat [post]
func (chatsHandler ChatsHandler) CreateGroupChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("CreateGroupChat").Inc()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	groupRequest := createGroupJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&groupRequest)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, err := usecase.CreateGroupChat(ctx, userID, groupRequest.Users, groupRequest.GroupName, groupRequest.Description, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, chatIDStruct{ChatID: chatID})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/createGroupChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// UpdateGroupChat updates group chat
//
// @Summary updates group chat
// @ID UpdateGroupChat
// @Accept json
// @Produce json
// @Param user body updateChatJson true "updated chat (если имя или описание не обновлялось, поле не слать вообще)"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /updateGroupChat [post]
func (chatsHandler ChatsHandler) UpdateGroupChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("UpdateGroupChat").Inc()
	ctx := r.Context()
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	updatedChat := updateChatJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updatedChat)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.UpdateGroupChat(ctx, userID, updatedChat.ChatID, updatedChat.NewName, updatedChat.NewDescription, chatsHandler.Chats)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
		misc.WriteInternalErrorJson(ctx, w)
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/updateGroupChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// GetPopularChannels updates group chat
//
// @Summary gets 10 popular channels
// @ID GetPopularChannels
// @Produce json
// @Success 200 {object}  domain.Response[getPopularChannelsResponse]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getPopularChannels [get]
func (chatsHandler ChatsHandler) GetPopularChannels(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("GetPopularChannels").Inc()
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	channels, err := usecase.GetPopularChannels(ctx, userID, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "Что-то пошло не так"})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, getPopularChannelsResponse{Channels: channels})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/updateGroupChat").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// JoinChannel joins channel
//
// @Summary joins channel
// @ID JoinChannel
// @Accept json
// @Produce json
// @Param user body chatIDStruct true "id of channel"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /joinChannel [post]
func (chatsHandler ChatsHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("JoinChannel").Inc()
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := chatIDStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.JoinChannel(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/joinChannel").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// LeaveChannel exits from channel
//
// @Summary exits from channel
// @ID LeaveChannel
// @Accept json
// @Produce json
// @Param user body chatIDStruct true "id of channel"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /leaveChannel [post]
func (chatsHandler ChatsHandler) LeaveChannel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("LeaveChannel").Inc()
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	chatIDStruct := chatIDStruct{}
	err := decoder.Decode(&chatIDStruct)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.LeaveChat(ctx, userID, chatIDStruct.ChatID, chatsHandler.Chats)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/leaveChannel").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

// CreateChannel creates channel
//
// @Summary creates channel
// @ID CreateChannel
// @Accept json
// @Produce json
// @Param user body createChannelJson true "IDs of users to create group chat with"
// @Success 200 {object}  domain.Response[chatIDStruct]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /createChannel [post]
func (chatsHandler ChatsHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("CreateChannel").Inc()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	channelRequest := createChannelJson{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&channelRequest)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	chatID, err := usecase.CreateChannel(ctx, userID, channelRequest.Name, channelRequest.Description, chatsHandler.Chats)
	if err != nil {
		if err.Error() == "internal error" {
			chatsHandler.prometheusMetrics.Errors.WithLabelValues("500").Inc()
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		logger.Error(err.Error())
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, chatIDStruct{ChatID: chatID})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/createChannel").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}

func (chatsHandler ChatsHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	chatsHandler.prometheusMetrics.Methods.WithLabelValues("GetMessages").Inc()
	ctx := r.Context()
	authorized, userID := chatsHandler.AuthHandler.CheckAuthNonAPI(w, r) // нули ли проверять userID на то, что он состоит в запрашиваемом чате?
	if !authorized {
		return
	}
	fmt.Println(userID)

	messageByChatIDRequest := messagesByChatIDRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&messageByChatIDRequest)
	if err != nil {
		chatsHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	messages := usecase.GetMessagesByChatID(ctx, chatsHandler.Chats, messageByChatIDRequest.ChatID)
	misc.WriteStatusJson(ctx, w, 200, domain.Messages{Messages: messages})
	duration := time.Since(start)
	chatsHandler.prometheusMetrics.requestDuration.WithLabelValues("/getMessages").Observe(duration.Seconds())
	chatsHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
}
