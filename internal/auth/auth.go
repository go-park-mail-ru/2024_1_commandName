package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
	"ProjectMessenger/internal/repository/inMemory"
)

type SessionStore interface {
	GetUserIDbySessionID(sessionID string) (userID uint, sessionExists bool)
	CreateSession(userID uint) (sessionID string)
	DeleteSession(sessionID string)
}

type UserStore interface {
	GetByUsername(username string) (user domain.Person, found bool)
	CreateUser(user domain.Person) (userID uint, err error)
}

type ChatStore interface {
	GetChatsByID(userID uint) []domain.Chat
}

type AuthHandler struct {
	sessions SessionStore
	users    UserStore
	chats    ChatStore
}

func NewAuthHandler() *AuthHandler {
	handler := AuthHandler{
		sessions: inMemory.NewSessionStorage(),
		users:    inMemory.NewUserStorage(),
		chats:    inMemory.NewChatsStorage(),
	}
	return &handler
}

type MyHandler struct {
	sessions map[string]*domain.Person
	users    map[string]*domain.Person
	chats    map[int]*domain.Chat
	chatUser []*domain.ChatUser
	sessMU   sync.RWMutex
	usersMU  sync.RWMutex
	chatsMU  sync.RWMutex
	isDebug  bool
}

// Login logs user in
//
// @Summary logs user in
// @ID login
// @Accept application/json
// @Produce application/json
// @Param user body  domain.Person true "Person"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure | user not found | wrong password"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /login [post]
func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if !errors.Is(err, http.ErrNoCookie) {
		_, sessionExists := authHandler.sessions.GetUserIDbySessionID(session.Value)
		if sessionExists {
			err := domain.WriteStatusJson(w, 400, domain.Error{Error: "session already exists"})
			if err != nil {
				domain.WriteInternalErrorJson(w)
				return
			}
			return
		}
	}
	if r.Method != http.MethodPost {
		err := domain.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
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
	var jsonUser domain.Person
	err = decoder.Decode(&jsonUser)
	if err != nil {
		http.Error(w, "wrong json structure", 400)
		return
	}
	if jsonUser.Username == "" {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}
	user, userFound := authHandler.users.GetByUsername(jsonUser.Username)
	if !userFound {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "Пользователь не найден"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	inputPassword := jsonUser.Password
	inputHash := misc.GenerateHash(inputPassword, user.PasswordSalt)
	if user.Password != inputHash {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "Неверный пароль"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	sessionID := authHandler.sessions.CreateSession(user.ID)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	err = domain.WriteStatusJson(w, 200, nil)
	if err != nil {
		domain.WriteInternalErrorJson(w)
		return
	}
}

// Logout logs user out
//
// @Summary logs user out
// @ID logout
// @Produce json
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "no session to logout"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /logout [get]
func (authHandler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "no session to logout"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	_, sessionExists := authHandler.sessions.GetUserIDbySessionID(session.Value)
	if !sessionExists {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "no session to logout"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	authHandler.sessions.DeleteSession(session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	err = domain.WriteStatusJson(w, 200, nil)
	if err != nil {
		domain.WriteInternalErrorJson(w)
		return
	}
}

// Register registers user
//
// @Summary registers user
// @ID register
// @Accept json
// @Produce json
// @Param user body  domain.Person true "Person"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "user already exists | required field empty | wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /register [post]
func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := domain.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
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
	var jsonUser domain.Person
	err := decoder.Decode(&jsonUser)
	if err != nil {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
	}
	if jsonUser.Username == "" || jsonUser.Password == "" {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "required field is empty"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	_, userFound := authHandler.users.GetByUsername(jsonUser.Username)
	if userFound {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "Пользователь с таким именем уже существет"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}
	passwordHash, passwordSalt := misc.GenerateHashAndSalt(jsonUser.Password)
	jsonUser.Password = passwordHash
	jsonUser.PasswordSalt = passwordSalt

	userID, err := authHandler.users.CreateUser(jsonUser)
	if err != nil {
		domain.WriteInternalErrorJson(w)
		return
	}
	sessionID := authHandler.sessions.CreateSession(userID)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	err = domain.WriteStatusJson(w, 200, nil)
	if err != nil {
		domain.WriteInternalErrorJson(w)
		return
	}
}

// CheckAuth checks that user is authenticated
//
// @Summary checks that user is authenticated
// @ID checkAuth
// @Produce json
// @Success 200 {object}  domain.Response[int]
// @Failure 401 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /checkAuth [get]
func (authHandler *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = authHandler.sessions.GetUserIDbySessionID(session.Value)
	}

	if authorized {
		err = domain.WriteStatusJson(w, 200, nil)
	} else {
		err = domain.WriteStatusJson(w, 401, domain.Error{Error: "Person not authorized"})
	}
	if err != nil {
		domain.WriteInternalErrorJson(w)
		return
	}
}

// GetChats gets chats previews for user
//
// @Summary gets chats previews for user
// @ID GetChats
// @Produce json
// @Success 200 {object}  domain.Response[domain.Chats]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getChats [get]
func (authHandler *AuthHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		err := domain.WriteStatusJson(w, 400, domain.Error{Error: "Person not authorized"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}
	userID, authorized := authHandler.sessions.GetUserIDbySessionID(session.Value)
	if !authorized {
		err = domain.WriteStatusJson(w, 400, domain.Error{Error: "Person not authorized"})
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}

	chats := authHandler.chats.GetChatsByID(userID)
	err = domain.WriteStatusJson(w, 200, domain.Chats{Chats: chats})
	if err != nil {
		errResp := domain.Error{Error: err.Error()}
		err := domain.WriteStatusJson(w, 500, errResp)
		if err != nil {
			domain.WriteInternalErrorJson(w)
			return
		}
		return
	}
}
