package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"ProjectMessenger/internal/auth/repository/InMemory"
	chatrepo "ProjectMessenger/internal/chats/repository/inMemory"
	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/usecase"
	chatusecase "ProjectMessenger/internal/chats/usecase"
	"ProjectMessenger/internal/misc"
)

type AuthHandler struct {
	Sessions usecase.SessionStore
	Users    usecase.UserStore
	Chats    chatusecase.ChatStore
}

func NewAuthHandler() *AuthHandler {
	handler := AuthHandler{
		Sessions: InMemory.NewSessionStorage(),
		Users:    InMemory.NewUserStorage(),
		Chats:    chatrepo.NewChatsStorage(),
	}
	return &handler
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
		sessionExists, _ := usecase.CheckAuthorized(session.Value, authHandler.Sessions)
		if sessionExists {
			misc.WriteStatusJson(w, 400, domain.Error{Error: "session already exists"})
			return
		}
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
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

	sessionID, err := usecase.LoginUser(jsonUser, authHandler.Users, authHandler.Sessions)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	misc.WriteStatusJson(w, 200, nil)
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
		misc.WriteStatusJson(w, 400, domain.Error{Error: "no session to logout"})
		return
	}

	sessionExists, _ := usecase.CheckAuthorized(session.Value, authHandler.Sessions)
	if !sessionExists {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "no session to logout"})
		return
	}

	usecase.LogoutUser(session.Value, authHandler.Sessions)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	misc.WriteStatusJson(w, 200, nil)
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
		misc.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
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
		misc.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
	}

	sessionID, err := usecase.RegisterAndLoginUser(jsonUser, authHandler.Users, authHandler.Sessions)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	misc.WriteStatusJson(w, 200, nil)
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
		authorized, _ = usecase.CheckAuthorized(session.Value, authHandler.Sessions)
	}
	if authorized {
		misc.WriteStatusJson(w, 200, nil)
	} else {
		misc.WriteStatusJson(w, 401, domain.Error{Error: "Person not authorized"})
	}
}
