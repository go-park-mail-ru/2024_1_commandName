package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/misc"

	//chatsInMemoryRepository "ProjectMessenger/internal/chats/repository/inMemory"
	repository "ProjectMessenger/internal/messages/repository/db"
	"ProjectMessenger/internal/messages/usecase"
)

type requestChatIDBody struct {
	ChatID uint `json:"chatID"`
}

type editMessageRequest struct {
	MessageID      uint   `json:"message_id"`
	NewMessageText string `json:"new_message_text"`
}

type deleteMessageRequest struct {
	MessageID uint `json:"message_id"`
}

type MessageHandler struct {
	ChatsHandler *delivery.ChatsHandler
	Websocket    usecase.WebsocketStore
	Messages     *repository.Messages
}

func NewMessagesHandler(chatsHandler *delivery.ChatsHandler, database *sql.DB) *MessageHandler {
	return &MessageHandler{
		ChatsHandler: chatsHandler,
		Websocket:    repository.NewWsStorage(database),
		Messages:     repository.NewMessageStorage(database),
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
func (messageHandler *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	fmt.Println(userID)

	upgrader := repository.UpgradeConnection()

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("SendMessage: upgrade failed", "err", err.Error())
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}
	user, found := messageHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("could not upgrade connection :user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}
	usecase.HandleWebSocket(ctx, connection, user, messageHandler.Websocket, messageHandler.Messages, messageHandler.ChatsHandler.Chats)
}

func (messageHandler *MessageHandler) SetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	_, found := messageHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("could not upgrade connection :user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "could not upgrade connection"})
		return
	}

	var requestToSetFile domain.File
	err := r.ParseMultipartForm(10000)
	if err != nil {
		fmt.Println(err)
	}

	files := r.MultipartForm.File["files"]
	jsonString := r.MultipartForm.Value["json"]
	json.Unmarshal([]byte(jsonString[0]), &requestToSetFile)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Fprintf(w, "Uploaded File: %+v\n", fileHeader.Filename)
		fmt.Fprintf(w, "File Size: %+v\n", fileHeader.Size)
		fmt.Fprintf(w, "MIME Header: %+v\n", fileHeader.Header)

		usecase.SetFile(messageHandler.Messages, ctx, file, userID, requestToSetFile.MessageID, messageHandler.ChatsHandler.AuthHandler.Users, fileHeader)
	}
}

//func (messageHandler *MessageHandler)

// GetChatMessages returns messages of some chat
//
// @Summary GetChatMessages
// @ID getChatMessages
// @Accept application/json
// @Produce application/json
// @Param user body  requestChatIDBody true "ID of chat"
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
	var RequestChatID requestChatIDBody
	err := decoder.Decode(&RequestChatID)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	limit := 100
	messages := usecase.GetChatMessages(r.Context(), limit, RequestChatID.ChatID, messageHandler.Messages)
	misc.WriteStatusJson(ctx, w, 200, domain.Messages{Messages: messages})
}

// EditMessage edits message
//
// @Summary EditMessage
// @ID editMessage
// @Accept application/json
// @Produce application/json
// @Param user body  editMessageRequest true "ID of chat"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /editMessage [post]
func (messageHandler *MessageHandler) EditMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var json editMessageRequest
	err := decoder.Decode(&json)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	err = usecase.EditMessage(ctx, userID, json.MessageID, json.NewMessageText, messageHandler.Messages)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}

// DeleteMessage deletes message
//
// @Summary DeleteMessage
// @ID deleteMessage
// @Accept application/json
// @Produce application/json
// @Param user body  deleteMessageRequest true "ID of message to delete"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /deleteMessage [post]
func (messageHandler *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var json deleteMessageRequest
	err := decoder.Decode(&json)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	err = usecase.DeleteMessage(ctx, userID, json.MessageID, messageHandler.Messages)
	if err != nil {
		if err == fmt.Errorf("internal error") {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}
