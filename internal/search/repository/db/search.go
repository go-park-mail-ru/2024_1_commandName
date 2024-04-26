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
	connection := s.GetConnection(userID)
	if connection == nil {
		return errors.New("No connection found for user")
	}
	return connection.WriteMessage(websocket.TextMessage, message)
}

func (s *Search) SearchChats(ctx context.Context, word string, userID uint) (foundChatsStructure domain.ChatSearchResponse) {
	wordsArr := strings.Split(word, " ")
	translatedWordsArr := s.TranslateWordWithTranslator(wordsArr)

	/*
		requestToSearch := ""
		for i:=0; i < len(translatedWordsArr); i++{
			requestToSearch += translatedWordsArr[i]
			rows, err := s.db.QueryContext(ctx,
				`SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id
					FROM chat.chat c
					JOIN chat.chat_user cu ON c.id = cu.chat_id
					WHERE name ILIKE $1 || '%' AND cu.user_id = $2`, requestToSearch, userID)
			if err != nil {
				//TODO
				fmt.Println("err:", err)
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
				mChat.Messages = s.Chats.GetMessagesByChatID(ctx, mChat.ID)
				if mChat.Messages != nil {
					mChat.Users = s.Chats.GetChatUsersByChatID(ctx, mChat.ID)
				}

				if mChat.Users != nil {
					matchedChats = append(matchedChats, mChat)
				}
			}
			if err = rows.Err(); err != nil {
				customErr := &domain.CustomError{
					Type:    "database",
					Message: err.Error(),
					Segment: "method searchChats, search.go",
				}
				fmt.Println("ERROR: ", customErr.Error())
				return foundChatsStructure
			}
		}
	*/

	rows, err := s.db.QueryContext(ctx,
		`SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id 
				FROM chat.chat c
				JOIN chat.chat_user cu ON c.id = cu.chat_id 
				WHERE name ILIKE $1 || '%' AND cu.user_id = $2`, word, userID)
	if err != nil {
		//TODO
		fmt.Println("err:", err)
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
		mChat.Messages = s.Chats.GetMessagesByChatID(ctx, mChat.ID)
		if mChat.Messages != nil {
			mChat.Users = s.Chats.GetChatUsersByChatID(ctx, mChat.ID)
		}

		if mChat.Users != nil {
			matchedChats = append(matchedChats, mChat)
		}
	}
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method searchChats, search.go",
		}
		fmt.Println("ERROR: ", customErr.Error())
		return foundChatsStructure
	}

	var chatSearchResponse domain.ChatSearchResponse
	chatSearchResponse.Chats = matchedChats
	chatSearchResponse.UserID = userID

	return chatSearchResponse
}

func ConvertToJSONResponse(chats []domain.Chat, userID uint) (jsonResponse []byte) {
	var chatSearchResponse domain.ChatSearchResponse
	chatSearchResponse.Chats = chats
	chatSearchResponse.UserID = userID
	jsonResponse, err := json.Marshal(chatSearchResponse)
	if err != nil {
		fmt.Println("err encoding JSON:", err)
	}
	return jsonResponse
}

func (s *Search) AddSearchIndexes(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_c ON chat.chat (id)")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_user_id ON chat.chat_user (user_id)")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_chat_id_cu ON chat.chat_user (chat_id);")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
}

func (s *Search) DeleteSearchIndexes(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_c CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_user_id CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
	_, err = s.db.ExecContext(ctx, "DROP INDEX IF EXISTS idx_chat_id_cu CASCADE")
	if err != nil {
		fmt.Println("err:", err)
		//TODO
	}
}

func (s *Search) SendMatchedSearchResponse(response domain.ChatSearchResponse) {

	jsonResp := ConvertToJSONResponse(response.Chats, response.UserID)
	err := s.WebSocket.SendMessageToUser(response.UserID, jsonResp)
	if err != nil {
		//TODO
		fmt.Println("ERROR:", err)
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
	syllToTranslate := make([]string, 0)

	for _, word := range words {
		for _, char := range word {
			counterOfLetters++
			fmt.Println(string(char))
		}

		for _, char := range word {
			currCount++
			if strings.ContainsRune(vowels, char) {
				counterOfVowels++
			}
			if counterOfVowels < 2 {
				fmt.Println("add", string(char))
				currSyllable += string(char)
			} else {
				syllToTranslate = append(syllToTranslate, currSyllable)
				currSyllable = ""
				currSyllable += string(char)
				counterOfVowels = 1
			}
			if counterOfLetters == currCount {
				fmt.Println(len(word), counterOfLetters)
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
	YandexConfig.TranslateKey = "Bearer t1.9euelZqelYrMyciLnJDHj5PKzpyclO3rnpWanMyVzMzLyJuXnJSQzZSQzJnl8_dlHlBO-e80ShNo_d3z9yVNTU757zRKE2j9zef1656VmozOzZPGlMidmZTHjcjNk86e7_zF656VmozOzZPGlMidmZTHjcjNk86e.dbhRbkheLJfmVeunG45CqjxpeIosd9qEl3g0HlRvQSQBnn3QvPOBklVEm5GxoOUKTBWvWJIxBTsOXvRb9fOIDA"
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
