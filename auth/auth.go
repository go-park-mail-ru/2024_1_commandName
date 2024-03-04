package auth

import (
	"ProjectMessenger/models"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

type Messenger struct {
	chats map[int]*models.Chat
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

	if len(api.users) > 3 {
		messenger := NewMessenger()
		messenger.fillDB(api.users)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	//Messenger.fillDB(api.users)
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
	api.sessions = make(map[string]uint)
}

func (m *Messenger) fillDB(users map[string]*models.Person) {
	m.chats = make(map[int]*models.Chat)
	for username, person := range users {
		fmt.Printf("Username: %s, ID: %d\n", username, person.ID)
	}
	if len(users) > 3 {
		messagesChat1 := make([]*models.Message, 0)
		messagesChat1 = append(messagesChat1,
			&models.Message{ID: 1, ChatID: 1, UserID: users["admin1"].ID, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
			&models.Message{ID: 2, ChatID: 1, UserID: users["admin2"].ID, Message: "Балдёж балдёж", Edited: false},
		)
		chat1 := models.Chat{Name: "noName", ID: 1, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1}
		m.chats[chat1.ID] = &chat1

		messagesChat2 := make([]*models.Message, 0)
		messagesChat2 = append(messagesChat2,
			&models.Message{ID: 1, ChatID: 2, UserID: users["admin3"].ID, Message: "Пойдём в столовку?", Edited: false},
			&models.Message{ID: 2, ChatID: 2, UserID: users["admin"].ID, Message: "Уже бегу", Edited: false},
		)
		chat2 := models.Chat{Name: "noName", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat2}
		m.chats[chat2.ID] = &chat2
		fmt.Println("MESSAGES Chat 1:")
		for _, message := range messagesChat1 {

			creatorUsername := users[findUser(message.UserID, users)]
			fmt.Printf("Message ID: %d\n", message.ID)
			fmt.Printf("Creator: %s (ID: %d)\n", creatorUsername, message.UserID)
			fmt.Printf("Message: %s\n", message.Message)
			fmt.Println("---------------------")
		}

		fmt.Println("MESSAGES Chat 2:")
		for _, message := range messagesChat2 {
			creatorUsername := users[findUser(message.UserID, users)]
			fmt.Printf("Message ID: %d\n", message.ID)
			fmt.Printf("Creator: %s (ID: %d)\n", creatorUsername, message.UserID)
			fmt.Printf("Message: %s\n", message.Message)
			fmt.Println("---------------------")
		}
	}
	fmt.Println("CHATS:", m.chats)

}

func (m *Messenger) getChats() {
	chats := m.chats

}

func NewMessenger() *Messenger {
	return &Messenger{
		chats: map[int]*models.Chat{},
	}
}

func findUser(ID uint, users map[string]*models.Person) string {
	for _, user := range users {
		if user.ID == ID {
			return user.Username
		}
	}
	return ""
}
