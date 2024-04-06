package delivery

import (
	profileUsecase "ProjectMessenger/internal/profile/usecase"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/repository/db"
	"ProjectMessenger/internal/auth/repository/inMemory"
	"ProjectMessenger/internal/auth/usecase"
	"ProjectMessenger/internal/misc"
)

type AuthHandler struct {
	Sessions usecase.SessionStore
	Users    usecase.UserStore
}

func NewAuthHandler(dataBase *sql.DB, avatarPath string) *AuthHandler {
	handler := AuthHandler{
		Sessions: db.NewSessionStorage(dataBase),
		Users:    db.NewUserStorage(dataBase, avatarPath),
	}
	return &handler
}

func NewAuthMemoryStorage() *AuthHandler {
	handler := AuthHandler{
		Sessions: inMemory.NewSessionStorage(),
		Users:    inMemory.NewUserStorage(),
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
	ctx := r.Context()
	session, err := r.Cookie("session_id")
	if !errors.Is(err, http.ErrNoCookie) {
		sessionExists, _ := usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
		if sessionExists {
			misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "session already exists"})
			return
		}
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
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
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
	}

	sessionID, err := usecase.LoginUser(ctx, jsonUser, authHandler.Users, authHandler.Sessions)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
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
	misc.WriteStatusJson(ctx, w, 200, nil)
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
	ctx := r.Context()
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "no session to logout"})
		return
	}
	sessionExists, _ := usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	if !sessionExists {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "no session to logout"})
		return
	}

	usecase.LogoutUser(ctx, session.Value, authHandler.Sessions)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	misc.WriteStatusJson(ctx, w, 200, nil)
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
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
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
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
	}

	sessionID, userID, err := usecase.RegisterAndLoginUser(ctx, jsonUser, authHandler.Users, authHandler.Sessions)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.Error()})
		return
	}

	ok := profileUsecase.AddToAllContacts(ctx, userID, authHandler.Users)
	if !ok {
		logger.Error("Register: contacts failed", "userID", userID)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	misc.WriteStatusJson(ctx, w, 200, nil)
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
	ctx := r.Context()
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	}
	if authorized {
		misc.WriteStatusJson(ctx, w, 200, nil)
	} else {
		misc.WriteStatusJson(ctx, w, 401, domain.Error{Error: "Person not authorized"})
	}
}

func (authHandler *AuthHandler) CheckAuthNonAPI(w http.ResponseWriter, r *http.Request) (authorized bool, userID uint) {
	ctx := r.Context()
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, userID = usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	}
	if !authorized {
		misc.WriteStatusJson(ctx, w, 401, domain.Error{Error: "Person not authorized"})
	}
	return authorized, userID
}
