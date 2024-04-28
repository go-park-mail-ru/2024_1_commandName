package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/chats/repository/db"
	"ProjectMessenger/internal/chats/usecase"
	ws "ProjectMessenger/internal/messages/repository/db"
	translaterepo "ProjectMessenger/internal/translate/repository/db"
	tl "ProjectMessenger/internal/translate/usecase"
	"github.com/gorilla/websocket"
)

type Search struct {
	db          *sql.DB
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
	Chats       usecase.ChatStore
	WebSocket   *ws.Websocket
	Translate   tl.TranslateStore
}

func UpgradeConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	return upgrader
}

func (s *Search) AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) context.Context {
	s.mu.Lock()
	s.WebSocket.Connections[userID] = connection
	s.mu.Unlock()
	ctx = context.WithValue(ctx, "ws userID", userID)
	logger := slog.With("requestID", ctx.Value("traceID")).With("ws userID", ctx.Value("ws userID"))
	logger.Debug("established ws")
	return ctx
}

func (s *Search) DeleteConnection(userID uint) {
	s.mu.Lock()
	delete(s.WebSocket.Connections, userID)
	s.mu.Unlock()
}

func (s *Search) GetConnection(userID uint) *websocket.Conn {
	s.mu.RLock()
	conn := s.WebSocket.Connections[userID]
	s.mu.RUnlock()
	return conn
}

func (s *Search) SendMessageToUser(userID uint, message []byte) error {
	fmt.Println("GET CONN FOR USER ", userID)
	connection := s.GetConnection(userID)
	if connection == nil {
		err := errors.New("no connection found for user")
		customErr := &domain.CustomError{
			Type:    "websocket",
			Message: err.Error(),
			Segment: "method SendMessageToUser, search.go",
		}
		fmt.Println(customErr.Error())
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (s *Search) SearchChats(ctx context.Context, word string, userID uint) (foundChatsStructure domain.ChatSearchResponse) {
	wordsArr := strings.Split(word, " ")
	translatedWordsArr := s.TranslateWordWithTranslator(wordsArr)
	translatedWordsWithRuneArr := s.TranslateWordWithRune(wordsArr)
	translatedWordsWithSyllableArr := s.TranslateWordWithSyllable(wordsArr)

	minLength := len(wordsArr)
	if len(translatedWordsArr) < minLength {
		minLength = len(translatedWordsArr)
	}
	if len(translatedWordsWithRuneArr) < minLength {
		minLength = len(translatedWordsWithRuneArr)
	}
	if len(translatedWordsWithSyllableArr) < minLength {
		minLength = len(translatedWordsWithSyllableArr)
	}

	fmt.Println("Search for words: ", wordsArr, translatedWordsArr, translatedWordsWithRuneArr, translatedWordsWithSyllableArr)
	if len(translatedWordsArr) > 0 {
		requestToSearchTranslator := ""
		requestToSearchOriginal := ""
		requestToSearchRune := ""
		requestToSearchSyllable := ""

		for i := 0; i < minLength; i++ {
			requestToSearchTranslator += translatedWordsArr[i]
			requestToSearchOriginal += wordsArr[i]
			requestToSearchRune += translatedWordsWithRuneArr[i]
			requestToSearchSyllable += translatedWordsWithSyllableArr[i]

			rows, err := s.db.QueryContext(ctx,
				`SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id
					FROM chat.chat c
					JOIN chat.chat_user cu ON c.id = cu.chat_id
					WHERE (name ILIKE $1 || '%' OR name ILIKE $2 || '%' OR name ILIKE $3 || '%' OR name ILIKE $4 || '%') AND cu.user_id = $5`, requestToSearchTranslator, requestToSearchOriginal, requestToSearchRune, requestToSearchSyllable, userID)
			if err != nil {
				customErr := &domain.CustomError{
					Type:    "database",
					Message: err.Error(),
					Segment: "method searchChats, search.go",
				}
				fmt.Println(customErr.Error())
				return foundChatsStructure
			}
			matchedChats := make([]domain.Chat, 0)
			for rows.Next() {
				var mChat domain.Chat
				err = rows.Scan(&mChat.ID, &mChat.Type, &mChat.Name, &mChat.Description, &mChat.AvatarPath, &mChat.CreatedAt, &mChat.LastActionDateTime, &mChat.CreatorID)
				if err != nil {
					customErr := &domain.CustomError{
						Type:    "database",
						Message: err.Error(),
						Segment: "method searchChats, search.go",
					}
					fmt.Println(customErr.Error())
					return foundChatsStructure
				}
				mChat.Messages = append(mChat.Messages, s.Chats.GetMessagesByChatID(ctx, mChat.ID)...)
				matchedChats = append(matchedChats, mChat)
				foundChatsStructure.Chats = append(foundChatsStructure.Chats, mChat)
			}
			if err = rows.Err(); err != nil {
				customErr := &domain.CustomError{
					Type:    "database",
					Message: err.Error(),
					Segment: "method searchChats, search.go",
				}
				fmt.Println(customErr.Error())
				return foundChatsStructure
			}
		}
	}
	foundChatsStructure.UserID = userID
	return foundChatsStructure
}

func ConvertToJSONResponse(chats []domain.Chat, userID uint) (jsonResponse []byte) {
	var chatSearchResponse domain.ChatSearchResponse
	chatSearchResponse.Chats = chats
	chatSearchResponse.UserID = userID
	jsonResponse, err := json.Marshal(chatSearchResponse)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "json",
			Message: err.Error(),
			Segment: "method ConvertToJSONResponse, search.go",
		}
		fmt.Println(customErr.Error())
	}
	return jsonResponse
}

func (s *Search) AddSearchIndexes(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_c ON chat.chat (id)")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_user_id ON chat.chat_user (user_id)")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_cu ON chat.chat_user (chat_id);")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
}

func (s *Search) DeleteSearchIndexes(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_c CASCADE")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_user_id CASCADE")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_cu CASCADE")
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddSearchIndexes, search.go",
		}
		fmt.Println(customErr.Error())
	}
}

func (s *Search) SendMatchedSearchResponse(response domain.ChatSearchResponse) {
	jsonResp := ConvertToJSONResponse(response.Chats, response.UserID)
	err := s.WebSocket.SendMessageToUser(response.UserID, jsonResp)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    ".WebSocket.SendMessageToUser",
			Message: err.Error(),
			Segment: "method SendMatchedSearchResponse, search.go",
		}
		fmt.Println(customErr.Error())
	}
}

func (s *Search) TranslateWordWithRune(words []string) (translatedWords []string) {
	letterMap := map[string]string{
		"а": "a",
		"б": "b",
		"в": "v",
		"г": "g",
		"д": "d",
		"е": "e",
		"ё": "yo",
		"ж": "zh",
		"з": "z",
		"и": "i",
		"й": "y",
		"к": "k",
		"л": "l",
		"м": "m",
		"н": "n",
		"о": "o",
		"п": "p",
		"р": "r",
		"с": "s",
		"т": "t",
		"у": "u",
		"ф": "f",
		"х": "kh",
		"ц": "ts",
		"ч": "ch",
		"ш": "sh",
		"щ": "shch",
		"ъ": "",
		"ы": "y",
		"ь": "",
		"э": "e",
		"ю": "yu",
		"я": "ya",
	}

	for _, word := range words {
		enWord := ""
		for _, char := range word {
			enLetter := letterMap[string(char)]
			if enLetter == "" {
				enWord += string(char)
			} else {
				enWord += enLetter
			}
		}
		translatedWords = append(translatedWords, enWord)
	}
	return translatedWords
}

func (s *Search) TranslateWordWithTranslator(words []string) (translatedWords []string) {
	var request domain.TranslateRequest
	request.TargetLanguageCode = "en"
	request.Messages = words
	request.FolderID = s.Translate.GetFolderID()
	response := s.Translate.Translate(request).Translations

	for i := 0; i < len(response); i++ {
		translatedWords = append(translatedWords, response[i].Text)
	}
	return translatedWords
}

func (s *Search) TranslateWordWithSyllable(words []string) (translatedWords []string) {
	vowels := "еыаоэяию"
	currSyllable := ""
	counterOfLetters := 0
	counterOfVowels := 0
	currCount := 0
	magicNumber := int32(100000) // char

	syllToTranslate := make([]string, 0)
	for _, word := range words {
		for _, char := range word {
			counterOfLetters++
			if char == magicNumber {
				fmt.Println(char)
			}
		}

		for _, char := range word {
			currCount++
			if strings.ContainsRune(vowels, char) {
				counterOfVowels++
			}
			if counterOfVowels < 2 {
				currSyllable += string(char)
			} else {
				syllToTranslate = append(syllToTranslate, currSyllable)
				currSyllable = ""
				currSyllable += string(char)
				counterOfVowels = 1
			}
			if counterOfLetters == currCount {
				syllToTranslate = append(syllToTranslate, currSyllable)
			}
		}
		syllToTranslate = append(syllToTranslate, " ")
	}
	translatedWords = s.TranslateWordWithTranslator(syllToTranslate)
	return translatedWords
}

func NewSearchStorage(database *sql.DB) *Search {
	slog.Info("created search storage")
	var YandexConfig domain.YandexConfig
	YandexConfig.TranslateKey = "Bearer t1.9euelZrJyJLOxsbGz8-LjpKejpHNzO3rnpWanMyVzMzLyJuXnJSQzZSQzJnl8_cIUUhO-e9hHS1M_t3z90h_RU7572EdLUz-zef1656VmsjLnMielonJlpuYjpuZm42W7_zF656VmsjLnMielonJlpuYjpuZm42W.iPQb5uK9StcUYyZioMcsyl4-0D2JnPTpJ1BUjDTpbf6pUleDTHbe65VZSJwILJx3nkX2v0UcUuHXcmAX0gRNCw"
	YandexConfig.Url = "https://translate.api.cloud.yandex.net/translate/v2/translate"
	YandexConfig.FolderID = "b1gq4i9e5unl47m0kj5f"
	YandexConfig.Header = "application/json"
	YandexConfig.Method = "POST"
	return &Search{
		db:          database,
		Connections: make(map[uint]*websocket.Conn),
		Chats:       db.NewChatsStorage(database),
		WebSocket:   ws.NewWsStorage(database),
		Translate:   translaterepo.NewTranslateStorage(database, YandexConfig),
	}
}