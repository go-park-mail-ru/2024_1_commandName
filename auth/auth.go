package auth

import (
	"ProjectMessenger/models"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/swaggo/http-swagger"
)

var (
	letterDigitRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*models.Person
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

func NewMyHandler() *MyHandler {
	adminHash, adminSalt := generateHashAndSalt("admin")
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*models.Person{
			"admin": {ID: 1, Username: "admin", Email: "admin@mail.ru", Name: "Ivan", Surname: "Ivanov",
				About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
		},
	}
}

// Login logs user in
//
// @Summary logs user in
// @ID login
// @Accept application/json
// @Produce application/json
// @Param user body  models.Person true "Person"
// @Success 200 {object}  models.Response
// @Failure 405 {object}  models.ErrorResponse "use POST"
// @Failure 400 {object}  models.ErrorResponse "wrong json structure | user not found | wrong password"
// @Router /login [post]
func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		err := models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		if err != nil {
			http.Error(w, "internal server error", 500)
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
		http.Error(w, "wrong json structure", 400)
		return
	}
	if jsonUser.Username == "" {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "wrong json structure"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	user, userFound := api.users[jsonUser.Username]
	if !userFound {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "user not found"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}

	inputPassword := jsonUser.Password
	inputHash := generateHash(inputPassword, user.PasswordSalt)
	if user.Password != inputHash {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "wrong password"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
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
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}
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
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	if _, ok := api.sessions[session.Value]; !ok {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}

	delete(api.sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		http.Error(w, "internal server error", 500)
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
// @Success 200 {object}  models.Response
// @Failure 405 {object}  models.ErrorResponse "use POST"
// @Failure 400 {object}  models.ErrorResponse "user already exists | required field empty | wrong json structure"
// @Router /register [post]
func (api *MyHandler) Register(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if r.Method != http.MethodPost {
		err := models.WriteStatusJson(w, 405, models.Error{Error: "use POST"})
		if err != nil {
			http.Error(w, "internal server error", 500)
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
			http.Error(w, "internal server error", 500)
			return
		}
	}
	if jsonUser.Username == "" || jsonUser.Password == "" {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "required field is empty"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	_, userFound := api.users[jsonUser.Username]
	if userFound {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "user already exists"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	jsonUser.ID = uint(len(api.users) + 1)
	passwordHash, passwordSalt := generateHashAndSalt(jsonUser.Password)
	jsonUser.Password = passwordHash
	jsonUser.PasswordSalt = passwordSalt

	api.users[jsonUser.Username] = &jsonUser
	sessionID := randStringRunes(32)

	api.sessions[sessionID] = jsonUser.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	err = models.WriteStatusJson(w, 200, nil)
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}
}

// CheckAuth checks that user is authenticated
//
// @Summary checks that user is authenticated
// @ID checkAuth
// @Produce json
// @Success 200 {object}  models.Response
// @Failure 401 {object}  models.ErrorResponse "Person not authorized"
// @Router /checkAuth [get]
func (api *MyHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

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
		http.Error(w, "internal server error", 500)
		return
	}
}

func Start() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/checkAuth", api.CheckAuth)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)
	r.HandleFunc("/register", api.Register)

	log.Fatal(http.ListenAndServe(":8080", r))
}
