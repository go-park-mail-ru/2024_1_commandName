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
	ChatID uint `json:"chat_id"`
}

type editMessageRequest struct {
	MessageID      uint   `json:"message_id"`
	NewMessageText string `json:"new_message_text"`
}

type deleteMessageRequest struct {
	MessageID uint `json:"message_id"`
}

type uploadFileDocs struct {
	Info domain.FileFromUser `json:"jsonq"`
}

type MessageHandler struct {
	ChatsHandler *delivery.ChatsHandler
	Websocket    usecase.WebsocketStore
	Messages     *repository.Messages
}

func NewMessagesHandler(chatsHandler *delivery.ChatsHandler, database *sql.DB, path string) *MessageHandler {
	return &MessageHandler{
		ChatsHandler: chatsHandler,
		Websocket:    repository.NewWsStorage(database),
		Messages:     repository.NewMessageStorage(database, path),
	}
}

func (messageHandler *MessageHandler) GetAllStickers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	_, found := messageHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "user wasn't found"})
		return
	}

	stickers := usecase.GetAllStickers(ctx, messageHandler.Messages)
	misc.WriteStatusJson(ctx, w, 200, stickers)
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
	usecase.HandleWebSocket(ctx, connection, user, messageHandler.Websocket, messageHandler.Messages, messageHandler.ChatsHandler.Chats, messageHandler.ChatsHandler.AuthHandler.Users, messageHandler.ChatsHandler.AuthHandler.Firebase)
}

// SetFile sets array of files
//
// @Summary sets array of files
// @ID SetFile
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "file to upload"
// @Param json formData string true "json with data (message_text,chat_id,type); type должен быть file"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /uploadFiles [post]
func (messageHandler *MessageHandler) SetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	_, found := messageHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "user wasn't found"})
		return
	}

	var requestToSetFile domain.FileFromUser
	err := r.ParseMultipartForm(10000)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "ParseMultipartForm",
			Message: err.Error(),
			Segment: "SetFile, messages_delivery.go",
		}
		fmt.Println(customErr.Error())
	}

	files := r.MultipartForm.File["files"]
	jsonString := r.MultipartForm.Value["json"]
	json.Unmarshal([]byte(jsonString[0]), &requestToSetFile)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			misc.WriteInternalErrorJson(ctx, w)
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "Uploaded FileFromUser: %+v\n", fileHeader.Filename)
		//fmt.Fprintf(w, "FileFromUser Size: %+v\n", fileHeader.Size)
		//fmt.Fprintf(w, "MIME Header: %+v\n", fileHeader.Header)

		usecase.SetFile(ctx, file, userID, fileHeader, requestToSetFile, messageHandler.Messages, messageHandler.ChatsHandler.AuthHandler.Users, messageHandler.Websocket, messageHandler.ChatsHandler.Chats, messageHandler.ChatsHandler.AuthHandler.Firebase)
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}

func (messageHandler *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	authorized, userID := messageHandler.ChatsHandler.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	user, found := messageHandler.ChatsHandler.AuthHandler.Users.GetByUserID(ctx, userID)
	if !found {
		logger.Info("user wasn't found")
		misc.WriteStatusJson(ctx, w, 500, domain.Error{Error: "user wasn't found"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var fileRequest domain.FileFromUser
	err := decoder.Decode(&fileRequest)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "decoder.Decode",
			Message: err.Error(),
			Segment: "SendSticker, messages_delivery.go",
		}
		fmt.Println(customErr.Error())
	}
	usecase.SendSticker(ctx, messageHandler.Messages, messageHandler.Websocket, messageHandler.ChatsHandler.Chats, fileRequest, user, messageHandler.ChatsHandler.AuthHandler.Users, messageHandler.ChatsHandler.AuthHandler.Firebase)
}

// GetChatMessages returns messages of some chat
// GetMessages returns messages of some chat
//
// @Summary GetMessages
// @ID getChatMessages
// @Accept application/json
// @Produce application/json
// @Param user body  requestChatIDBody true "ID of chat"
// @Success 200 {object}  domain.Response[domain.Messages]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getMessages [post]
func (messageHandler *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
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
