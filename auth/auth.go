package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"
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
	adminHash, adminSalt := generateHashAndSalt("Admin123.")
	handler := &MyHandler{
		sessions: make(map[string]*models.Person, 10),
		users: map[string]*models.Person{
			"admin": {ID: 1, Username: "admin", Email: "admin@mail.ru", Name: "Ivan", Surname: "Naumov",
				About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
			"ArtemkaChernikov": {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
				About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
			"ArtemZhuk": {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
				About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
			"AlexanderVolohov": {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
				About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
			"mentor": {ID: 4, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
				About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
				PasswordSalt: adminSalt, Password: adminHash},
		},
		chats:    make(map[int]*models.Chat),
		chatUser: make([]*models.ChatUser, 0),
		isDebug:  isDebug,
	}
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
				http.Error(w, "internal server error", 500)
				return
			}
			return
		}
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
	err = decoder.Decode(&jsonUser)
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
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь не найден"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}

	inputPassword := jsonUser.Password
	inputHash := generateHash(inputPassword, user.PasswordSalt)
	if user.Password != inputHash {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Неверный пароль"})
		if err != nil {
			http.Error(w, "internal server error", 500)
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
// @Success 200 {object}  models.Response[int]
// @Failure 400 {object}  models.Response[models.Error] "no session to logout"
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
// @Success 200 {object}  models.Response[int]
// @Failure 405 {object}  models.Response[models.Error] "use POST"
// @Failure 400 {object}  models.Response[models.Error] "user already exists | required field empty | wrong json structure"
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
		err := models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь с таким именем уже существет"})
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

	api.sessions[sessionID] = &jsonUser

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
// @Success 200 {object}  models.Response[int]
// @Failure 401 {object}  models.Response[models.Error] "Person not authorized"
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
		http.Error(w, "internal server error", 500)
		return
	}
}

func (api *MyHandler) ClearUserData() {
	api.users = make(map[string]*models.Person)
	api.sessions = make(map[string]*models.Person)
}

func (api *MyHandler) fillDB() {
	messagesChat1 := make([]*models.Message, 0)
	messagesChat1 = append(messagesChat1,
		&models.Message{ID: 1, ChatID: 1, UserID: api.users["mentor"].ID, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
		//&models.Message{ID: 2, ChatID: 1, UserID: api.users["admin1"].ID, Message: "Балдёж балдёж", Edited: false},
	)
	chat1 := models.Chat{Name: "noName", ID: 1, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1}
	api.chats[chat1.ID] = &chat1

	messagesChat2 := make([]*models.Message, 0)
	messagesChat2 = append(messagesChat2,
		&models.Message{ID: 1, ChatID: 2, UserID: api.users["ArtemkaChernikov"].ID, Message: "Пойдём в столовку?", Edited: false},
		//&models.Message{ID: 2, ChatID: 2, UserID: api.users["admin3"].ID, Message: "Уже бегу", Edited: false},
	)
	chat2 := models.Chat{Name: "noName", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat2}
	api.chats[chat2.ID] = &chat2

	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 1, UserID: 1})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 1, UserID: 5})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 1, UserID: 2})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 2, UserID: 3})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 2, UserID: 4})

}

// GetChats gets chats previews for user
//
// @Summary gets chats previews for user
// @ID GetChats
// @Produce json
// @Success 200 {object}  models.Response[models.Chats]
// @Failure 400 {object}  models.Response[models.Error] "Person not authorized"
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
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	user := api.sessions[session.Value]
	if user == nil {
		err = models.WriteStatusJson(w, 400, models.Error{Error: "Person not authorized"})
		if err != nil {
			http.Error(w, "internal server error", 500)
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
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
}

func (api *MyHandler) getChatsByID(userID uint) []*models.Chat {
	userChats := make(map[int]*models.Chat)
	for _, cUser := range api.chatUser {
		if cUser.UserID == userID {
			chat, ok := api.chats[cUser.ChatID]
			if ok {
				userChats[cUser.ChatID] = chat
			}
		}
	}

	var chats []*models.Chat
	for _, chat := range userChats {
		chats = append(chats, chat)
	}
	return chats
}
