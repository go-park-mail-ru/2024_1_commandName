package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"ProjectMessenger/domain"
	authusecase "ProjectMessenger/internal/auth/usecase"
	"ProjectMessenger/internal/misc"
	"gopkg.in/yaml.v3"
)

type Messages struct {
	db                  *sql.DB
	pathToStorageFolder string
}

func NewMessageStorage(db *sql.DB, path string) *Messages {
	return &Messages{db: db, pathToStorageFolder: path}
}

func (m *Messages) SetMessage(ctx context.Context, message domain.Message) (messageSaved domain.Message) {
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	query := "INSERT INTO chat.message (user_id, chat_id, message, edited_at, created_at, sticker_path) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	var messageID uint
	err := m.db.QueryRowContext(ctx, query, message.UserID, message.ChatID, message.Message, message.EditedAt, message.CreatedAt, message.StickerPath).Scan(&messageID)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "SetMessage, messages.go",
		}
		fmt.Println(customErr.Error())
		return domain.Message{}
	}
	message.ID = messageID
	logger.Debug("SetMessage: success", "msg", message)
	return message
}

func (m *Messages) SetFile(ctx context.Context, multipartFile multipart.File, userID uint, messageID uint, request domain.FileFromUser, userStorage authusecase.UserStore, fileHandler *multipart.FileHeader) error {
	_, found := userStorage.GetByUserID(ctx, userID)
	if !found {
		customErr := domain.CustomError{
			Type:    "find user by ID",
			Message: "user not found",
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
		return customErr
	}
	buff := make([]byte, 8192)
	if _, err := multipartFile.Read(buff); err != nil {
		customErr := domain.CustomError{
			Type:    "read multipart file",
			Message: err.Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
		return customErr
	}
	seek, err := multipartFile.Seek(0, io.SeekStart)
	if err != nil || seek != 0 {
		customErr := domain.CustomError{
			Type:    "seek multipart file",
			Message: err.Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
		return customErr
	}
	mimeType := http.DetectContentType(buff)
	fmt.Println(mimeType) // TODO CHECK TYPE AND SIZE

	filePath, err := m.StoreFile(ctx, multipartFile, fileHandler)
	query := "INSERT INTO chat.file (message_id, file_path, type, originalname) VALUES($1, $2, $3, $4)"
	row := m.db.QueryRowContext(ctx, query, messageID, filePath, request.AttachmentType, fileHandler.Filename)
	fmt.Println("INSERTING", userID, messageID, filePath)
	if row.Err() != nil {
		fmt.Println("ERR:")
		customErr := domain.CustomError{
			Type:    "database",
			Message: row.Err().Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
		return customErr
	}
	return nil
}

func (m *Messages) GetFilePathByMessageID(ctx context.Context, messageID uint) (filePath []string) {
	query := "SELECT file_path FROM chat.file WHERE message_id =$1"
	rows, err := m.db.QueryContext(ctx, query, messageID)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "GetFilePathByMessageID",
			Message: err.Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
	}

	filePath = make([]string, 0)
	for rows.Next() {
		path := ""
		err = rows.Scan(&path)
		filePath = append(filePath, path)
		if err != nil {
			customErr := domain.CustomError{
				Type:    "GetFilePathByMessageID",
				Message: err.Error(),
				Segment: "SetFiles, messages.go",
			}
			fmt.Println(customErr.Error())
		}
	}
	return filePath
}

func (m *Messages) GetFileByPath(filePath string) (file *os.File, fileInfo os.FileInfo) {
	file, err := os.Open(filePath)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "GetFileByPath",
			Message: err.Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
	}
	fileInfo, err = file.Stat()
	if err != nil {
		customErr := domain.CustomError{
			Type:    "GetFileByPath",
			Message: err.Error(),
			Segment: "SetFiles, messages.go",
		}
		fmt.Println(customErr.Error())
	}
	return file, fileInfo
}

func (m *Messages) StoreFile(ctx context.Context, multipartFile multipart.File, fileHandler *multipart.FileHeader) (filePath string, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	originalName := fileHandler.Filename
	fileNameSlice := strings.Split(originalName, ".")
	if len(fileNameSlice) < 2 {
		logger.Info("StoreAvatar filename without extension")
		return "", fmt.Errorf("Файл не имеет расширения")
	}
	extension := fileNameSlice[len(fileNameSlice)-1]

	filename := misc.RandStringRunes(20)
	filePath = "files/" + filename + "." + extension

	f, err := os.OpenFile(m.pathToStorageFolder+filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "os open file",
			Message: err.Error(),
			Segment: "method StoreFile, messages.go",
		}
		fmt.Println(customErr.Error())
		logger.Error("StoreFile failed to open a file", "path", m.pathToStorageFolder+filePath)
		return "", fmt.Errorf("internal error")
	}
	defer f.Close()

	_, err = io.Copy(f, multipartFile)
	if err != nil {
		logger.Error("StoreFile failed to copy file", "path", m.pathToStorageFolder+filePath)
		return "", fmt.Errorf("internal error")
	}
	logger.Debug("StoreFile success", "path", m.pathToStorageFolder+filePath)
	return filePath, nil
}

func (m *Messages) GetAllStickers(ctx context.Context) (stickers []domain.Sticker) {
	query := "SELECT id, description, type, file_path FROM chat.sticker ORDER BY id"
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetAllStickers, messages.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}

	stickers = make([]domain.Sticker, 0)
	for rows.Next() {
		sticker := domain.Sticker{}
		err = rows.Scan(&sticker.StickerID, &sticker.StickerDesc, &sticker.StickerType, &sticker.StickerPath)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetAllStickers, messages.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		stickers = append(stickers, sticker)
	}
	return stickers
}

func (m *Messages) GetStickerPathByID(ctx context.Context, stickerID uint) (filePah string) {
	query := "SELECT file_path FROM chat.sticker WHERE id = $1"
	row := m.db.QueryRowContext(ctx, query, stickerID)
	err := row.Scan(&filePah)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filePah
}

func (m *Messages) GetChatMessages(ctx context.Context, chatID uint, limit int) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)
	rows, err := m.db.QueryContext(ctx, "SELECT message.id, user_id, chat_id, message.message, COALESCE(message.created_at, '2000-01-01 00:00:00'), COALESCE(message.edited_at, '2000-01-01 00:00:00'), username, COALESCE(originalname, '') AS originalname, COALESCE(file_path, '') AS file_path, COALESCE(type, '') AS type, COALESCE(sticker_path, '') AS sticker_path FROM chat.message JOIN auth.person ON message.user_id = person.id LEFT JOIN chat.file f on message.id = f.message_id WHERE chat_id = $1 ORDER BY chat.message.created_at", chatID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetChatMessages, messages.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var mess domain.Message
		mess.File = &domain.FileInMessage{}
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.CreatedAt, &mess.EditedAt, &mess.SenderUsername, &mess.File.OriginalName, &mess.File.Path, &mess.File.Type, &mess.StickerPath); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetMessagesByChatID, profile.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		if mess.File.Path == "" {
			mess.File = nil
		}
		chatMessagesArr = append(chatMessagesArr, mess)
	}
	return chatMessagesArr
}

func (m *Messages) GetMessage(ctx context.Context, messageID uint) (message domain.Message, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	message = domain.Message{}
	err = m.db.QueryRowContext(ctx, "SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), message.created_at FROM chat.message WHERE id = $1", messageID).Scan(
		&message.ID, &message.UserID, &message.ChatID, &message.Message, &message.Edited, &message.EditedAt, &message.CreatedAt)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("EditMessage didn't found message", "messageID", messageID)
			return message, fmt.Errorf("Такого сообщения не существует")
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetChatsForUser, chats.go",
		}
		logger.Error(err.Error(), "EditMessage db error", customErr.Message)
		return message, customErr
	}
	return message, nil
}

func (m *Messages) UpdateMessageText(ctx context.Context, message domain.Message) (err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	_, err = m.db.ExecContext(ctx, "UPDATE chat.message SET message = $1, edited = $2, edited_at = $3 WHERE id = $4", message.Message, message.Edited, message.EditedAt, message.ID)
	if err != nil {
		customErr := domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "UpdateMessageText, messages.go",
		}
		logger.Error("UpdateMessageText db error", "messageID", message.ID)
		return customErr
	}
	return nil
}

func (m *Messages) DeleteMessage(ctx context.Context, messageID uint) error {
	logger := slog.With("requestID", ctx.Value("traceID"))
	_, err := m.db.ExecContext(ctx, "DELETE FROM chat.message WHERE id = $1", messageID)
	if err != nil {
		logger.Error("DeleteMessage db error", "messageID", messageID)
		return fmt.Errorf("internal error")
	}
	return nil
}

func (m *Messages) SummarizeMessage(message domain.SummarizeMessageRequest) domain.TranslateResponse {
	stream := false
	temperature := 0.6
	maxTokens := "2000"

	reqData := domain.APIRequest{
		ModelURI: "gpt://b1gq4i9e5unl47m0kj5f/yandexgpt/latest",
		CompletionOptions: domain.CompletionOptions{
			Stream:      stream,
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
		Messages: []domain.SummarizeMessageRequest{
			{
				Role: "system",
				Text: "Выдели очень кратко основные мысли из сообщения от" + message.Username,
			},
			{
				Role: "user",
				Text: message.Text,
			},
		},
	}
	jsonRequest, err := json.Marshal(reqData)

	var YandexConfig domain.YandexConfig
	cfg := LoadConfig()
	YandexConfig.TranslateKey = cfg.Gpt.TrKey
	YandexConfig.Url = cfg.Gpt.Url
	YandexConfig.FolderID = cfg.Gpt.FolderID
	YandexConfig.Header = cfg.Gpt.Header
	YandexConfig.Method = cfg.Gpt.Method

	client := &http.Client{}
	req, err := http.NewRequest(YandexConfig.Method, YandexConfig.Url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "http new request",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	req.Header.Add("Content-Type", YandexConfig.Header)
	req.Header.Add("Authorization", YandexConfig.TranslateKey)

	fmt.Println(req.Method)

	resp, err := client.Do(req)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "http do request",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "read response",
			Message: err.Error(),
			Segment: "method Translate, translate.go",
		}
		fmt.Println(customErr.Error())
	}
	gptResponse := ParseSummarizeResponse(body)
	var summarizeResp domain.TranslateResponse
	gptToTranslations := domain.Translations{Text: gptResponse.Result.Alternatives[0].Message.Text}
	summarizeResp.Translations = append(summarizeResp.Translations, gptToTranslations)
	return summarizeResp
}

func LoadConfig() domain.Config {
	envPath := os.Getenv("GOCHATME_HOME")
	slog.Debug("env home =" + envPath)
	f, err := os.Open(envPath + "config.yml")
	slog.Debug("trying to open " + envPath + "config.yml")
	if err != nil {
		slog.Error("load config failed", "err", err)
		fmt.Errorf("load config failed").Error()
	}
	defer f.Close()

	var cfg domain.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func ParseSummarizeResponse(jsonResponse []byte) (response domain.SummarizeMessageResponse) {
	err := json.Unmarshal(jsonResponse, &response)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "json Unmarshal",
			Message: err.Error(),
			Segment: "method ParseSummarizeResponse, messages.go",
		}
		fmt.Println(customErr.Error())
	}
	fmt.Println(response.Result.Alternatives[0].Message.Text)
	return response
}
