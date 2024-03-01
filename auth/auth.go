package auth

import (
	"ProjectMessenger/models"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/swaggo/http-swagger"
)

var (
	letterDigitRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

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

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*models.Person
}

func NewMyHandler() *MyHandler {
	adminHash, adminSalt := generateHashAndSalt("admin")
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*models.Person{
			"admin": {ID: 1, Username: "admin", Email: "admin@mail.ru", Name: "Ivan", Surname: "Ivanov",
				About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, PasswordHash: adminHash},
		},
	}
}

// Login logs user in
//
// @Summary logs user in
// @ID login
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object}  models.Response
// @Failure 405 {object}  models.ErrorResponse "Use POST"
// @Failure 400 {object}  models.ErrorResponse "Username or password wrong"
// @Router /login [post]
func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		return
	}
	username := r.FormValue("username")
	if username == "" {
		models.WriteStatusJson(w, 400, models.Error{Error: "username is not present in request"})
		return
	}
	user, userFound := api.users[r.FormValue("username")]
	if !userFound {
		models.WriteStatusJson(w, 400, models.Error{Error: "user not found"})
		return
	}

	inputPassword := r.FormValue("password")
	inputHash := generateHash(inputPassword, user.PasswordSalt)

	if user.PasswordHash != inputHash {
		models.WriteStatusJson(w, 400, models.Error{Error: "wrong password"})
		return
	}

	SID := randStringRunes(32)

	api.sessions[SID] = user.ID

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	models.WriteStatusJson(w, 200, nil)
}

// Logout logs user out
//
// @Summary logs user out
// @ID logout
// @Produce json
// @Success 200 {object}  models.Response
// @Failure 400 {object}  models.ErrorResponse "no session to logout"
// @Router /logout [get]
func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		return
	}
	if _, ok := api.sessions[session.Value]; !ok {
		models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		return
	}

	delete(api.sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	models.WriteStatusJson(w, 200, nil)
}

// Register registers user
//
// @Summary registers user
// @ID register
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Success 200 {object}  models.Response
// @Failure 405 {object}  models.ErrorResponse "Use POST"
// @Failure 400 {object}  models.ErrorResponse "Username exists or field required field empty"
// @Router /register [post]
func (api *MyHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	if username == "" || password == "" || email == "" {
		models.WriteStatusJson(w, 400, models.Error{Error: "required field is empty"})
		return
	}
	_, userFound := api.users[username]
	if userFound {
		models.WriteStatusJson(w, 400, models.Error{Error: "user already exists"})
		return
	}
	ID := uint(len(api.users) + 1)
	passwordHash, passwordSalt := generateHashAndSalt(password)
	newUser := &models.Person{
		ID:           ID,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		Username:     username,
		Email:        email,
	}
	api.users[newUser.Username] = newUser
	SID := randStringRunes(32)

	api.sessions[SID] = newUser.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	models.WriteStatusJson(w, 200, nil)
}

func (api *MyHandler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		models.WriteStatusJson(w, 200, nil)
	} else {
		models.WriteStatusJson(w, 401, models.Error{Error: "User not authorized"})
	}
}

// @Title Messenger authorization API
// @Version 1.0
// @BasePath /

func Start() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)
	r.HandleFunc("/register", api.Register)

	http.ListenAndServe(":8080", r)
}
