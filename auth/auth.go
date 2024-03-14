package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/models"
)

var (
	letterDigitRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type MyHandler struct {
	sessions map[string]*models.Person
	users    map[string]*models.Person
	chats    map[int]*models.Chat
	chatUser []*models.ChatUser
	sessMU   sync.RWMutex
	usersMU  sync.RWMutex
	chatsMU  sync.RWMutex
	isDebug  bool
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterDigitRunes[rand.Intn(len(letterDigitRunes))]
	}
	return string(b)
}

func generateHashAndSalt(password string) (hash string, salt string) {
	salt = randStringRunes(8)
	hasher := sha512.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil)), salt
}

func generateHash(password string, salt string) (hash string) {
	hasher := sha512.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func setDebugHeaders(w http.ResponseWriter, r *http.Request) (needToReturn bool) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "http://localhost:3000")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	header.Add("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		needToReturn = true
	}

	return needToReturn
}

func NewMyHandler(isDebug bool) *MyHandler {
	handler := &MyHandler{
		sessions: make(map[string]*models.Person, 10),
		chats:    make(map[int]*models.Chat),
		chatUser: make([]*models.ChatUser, 0),
		isDebug:  isDebug,
		sessMU:   sync.RWMutex{},
		chatsMU:  sync.RWMutex{},
		usersMU:  sync.RWMutex{},
	}
	handler.users = handler.fillUsers()
	handler.fillDB()
	return handler
}

// Login logs user in
//
// @Summary logs user in
// @ID login
// @Accept application/json
// @Produce application/json
// @Param user body  models.Person true "Person"
// @Success 200 {object}  models.Response[int]
// @Failure 405 {object}  models.Response[models.Error] "use POST"
// @Failure 400 {object}  models.Response[models.Error] "wrong json structure | user not found | wrong password"
// @Failure 500 {object}  models.Response[models.Error] "Internal server error"
// @Router /login [post]
func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	session, err := r.Cookie("session_id")
	if !errors.Is(err, http.ErrNoCookie) {
		if _, ok := api.sessions[session.Value]; ok {
			err := models.WriteStatusJson(w, 400, models.Error{Error: "session already exists"})
			if err != nil {
				models.WriteInternalErrorJson(w)
				return
			}
			return
		}
	}
	if r.Method != http.MethodPost {
		err := models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser models.Person
	err = decoder.Decode(&jsonUser)
	if err != nil {
		http.Error(w, "wrong json structure", 400)
		return
	}
	if jsonUser.Username == "" {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "wrong json structure"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	user, userFound := api.users[jsonUser.Username]
	if !userFound {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь не найден"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}

	inputPassword := jsonUser.Password
	inputHash := generateHash(inputPassword, user.PasswordSalt)
	if user.Password != inputHash {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Неверный пароль"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}

	SID := randStringRunes(32)
	api.sessions[SID] = user
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		models.WriteInternalErrorJson(w)
		return
	}
}

// Logout logs user out
//
// @Summary logs user out
// @ID logout
// @Produce json
// @Success 200 {object}  models.Response[int]
// @Failure 400 {object}  models.Response[models.Error] "no session to logout"
// @Failure 500 {object}  models.Response[models.Error] "Internal server error"
// @Router /logout [get]
func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	if _, ok := api.sessions[session.Value]; !ok {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	api.sessMU.Lock()
	delete(api.sessions, session.Value)
	api.sessMU.Unlock()
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		models.WriteInternalErrorJson(w)
		return
	}
}

// Register registers user
//
// @Summary registers user
// @ID register
// @Accept json
// @Produce json
// @Param user body  models.Person true "Person"
// @Success 200 {object}  models.Response[int]
// @Failure 405 {object}  models.Response[models.Error] "use POST"
// @Failure 400 {object}  models.Response[models.Error] "user already exists | required field empty | wrong json structure"
// @Failure 500 {object}  models.Response[models.Error] "Internal server error"
// @Router /register [post]
func (api *MyHandler) Register(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	if r.Method != http.MethodPost {
		err := models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser models.Person
	err := decoder.Decode(&jsonUser)
	if err != nil {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "wrong json structure"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
	}
	if jsonUser.Username == "" || jsonUser.Password == "" {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "required field is empty"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	api.usersMU.Lock()
	_, userFound := api.users[jsonUser.Username]
	if userFound {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь с таким именем уже существет"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			api.usersMU.Unlock()
			return
		}
		api.usersMU.Unlock()
		return
	}
	jsonUser.ID = uint(len(api.users) + 1)
	passwordHash, passwordSalt := generateHashAndSalt(jsonUser.Password)
	jsonUser.Password = passwordHash
	jsonUser.PasswordSalt = passwordSalt

	api.users[jsonUser.Username] = &jsonUser
	api.usersMU.Unlock()
	sessionID := randStringRunes(32)
	api.sessMU.Lock()
	api.sessions[sessionID] = &jsonUser
	api.sessMU.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		models.WriteInternalErrorJson(w)
		return
	}
}

// CheckAuth checks that user is authenticated
//
// @Summary checks that user is authenticated
// @ID checkAuth
// @Produce json
// @Success 200 {object}  models.Response[int]
// @Failure 401 {object}  models.Response[models.Error] "Person not authorized"
// @Failure 500 {object}  models.Response[models.Error] "Internal server error"
// @Router /checkAuth [get]
func (api *MyHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		err = models.WriteStatusJson(w, 200, nil)
	} else {
		err = models.WriteStatusJson(w, 401, models.Error{Error: "Person not authorized"})
	}
	if err != nil {
		models.WriteInternalErrorJson(w)
		return
	}
}

func (api *MyHandler) ClearUserData() {
	api.users = make(map[string]*models.Person)
	api.sessions = make(map[string]*models.Person)
}

// GetChats gets chats previews for user
//
// @Summary gets chats previews for user
// @ID GetChats
// @Produce json
// @Success 200 {object}  models.Response[models.Chats]
// @Failure 400 {object}  models.Response[models.Error] "Person not authorized"
// @Failure 500 {object}  models.Response[models.Error] "Internal server error"
// @Router /getChats [get]
func (api *MyHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Person not authorized"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	user := api.sessions[session.Value]
	if user == nil {
		err = models.WriteStatusJson(w, 400, models.Error{Error: "Person not authorized"})
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
	chats := api.getChatsByID(user.ID)
	err = models.WriteStatusJson(w, 200, models.Chats{Chats: chats})
	if err != nil {
		errResp := models.Error{Error: err.Error()}
		err := models.WriteStatusJson(w, 500, errResp)
		if err != nil {
			models.WriteInternalErrorJson(w)
			return
		}
		return
	}
}

func (api *MyHandler) getChatsByID(userID uint) []*models.Chat {
	userChats := make(map[int]*models.Chat)
	api.chatsMU.Lock()
	for _, cUser := range api.chatUser {
		if cUser.UserID == userID {
			chat, ok := api.chats[cUser.ChatID]
			if ok {
				userChats[cUser.ChatID] = chat
			}
		}
	}
	api.chatsMU.Unlock()

	var chats []*models.Chat
	for _, chat := range userChats {
		chats = append(chats, chat)
	}
	return chats
}

func (api *MyHandler) getChatUsersByChatID(chatID int) []*models.ChatUser {
	usersOfChat := make([]*models.ChatUser, 0)
	for i := range api.chatUser {
		if api.chatUser[i].ChatID == chatID {
			usersOfChat = append(usersOfChat, api.chatUser[i])
		}
	}
	return usersOfChat
}
