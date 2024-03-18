package auth

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/models"
	_ "github.com/lib/pq"
)

var (
	letterDigitRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type ChatMe struct {
	db      *sql.DB
	isDebug bool
}

func CreateChatMeHandler(isDebug bool) (*ChatMe, error) {
	connStrToDataBase := "user=postgres dbname=Messenger password=Artem557 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStrToDataBase)
	if err != nil {
		return nil, err
	}

	chatMe := &ChatMe{
		db:      db,
		isDebug: isDebug,
	}
	chatMe.fillDB()

	return chatMe, err
	//fillUsers()
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
func (api *ChatMe) Login(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	session, err := r.Cookie("session_id")
	fmt.Println("COOKIE: ", session)
	if !errors.Is(err, http.ErrNoCookie) {
		ok := true
		if ok, err = api.isSessionExistByValue(session.Value); ok {
			//if _, ok, err = api.getSessionByCookieValue(session.Value); ok { // DONE
			fmt.Println("err =", err)
			err = models.WriteStatusJson(w, 400, models.Error{Error: "session already exists"})
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
	//user, userFound := api.users[jsonUser.Username]
	user, userFound, err := api.getUserByUsername(jsonUser.Username)
	if err != nil {
		// TODO
	}
	fmt.Println("user = ", user, err, userFound)
	if !userFound {
		err = models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь не найден"})
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
	//.sessions[SID] = user // ALREADY MADE
	err = api.setSessionBySessionID(SID, user)
	if err != nil {
		/// TODO
	}
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
func (api *ChatMe) Logout(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	session, err := r.Cookie("session_id")
	fmt.Println("COOKIE: ", session)
	if errors.Is(err, http.ErrNoCookie) {
		err := models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	var ok bool
	//if _, ok := api.sessions[session.Value]; !ok {
	if ok, err = api.isSessionExistByValue(session.Value); !ok { // DONE
		err = models.WriteStatusJson(w, 400, models.Error{Error: "no session to logout"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}

	//delete(api.sessions, session.Value) // ALREADY DONE
	err = api.deleteSessionByCookieValue(session.Value)
	if err != nil {
		fmt.Println(err)
		// TODO
	}

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
func (api *ChatMe) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("register")
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
	//_, userFound := api.users[jsonUser.Username]
	_, userFound, err := api.getUserByUsername(jsonUser.Username) // DONE
	fmt.Println("err = ", err)
	if userFound {
		err = models.WriteStatusJson(w, 400, models.Error{Error: "Пользователь с таким именем уже существет"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	//jsonUser.ID = uint(len(api.users) + 1) // ALREADY DONE
	countOfUsers, err := api.getCountOfUsers()
	if err != nil {
		// TODO
	}
	jsonUser.ID = uint(countOfUsers) + 1
	passwordHash, passwordSalt := generateHashAndSalt(jsonUser.Password)
	jsonUser.Password = passwordHash
	jsonUser.PasswordSalt = passwordSalt

	err = api.setUserByUsername(jsonUser)
	if err != nil {
		fmt.Println("err in setUser: ", err)
		//TODO
	}
	//api.users[jsonUser.Username] = &jsonUser // ALREADY DONE
	sessionID := randStringRunes(32)
	err = api.setSessionBySessionID(sessionID, jsonUser)
	if err != nil {
		fmt.Println("err in setSession: ", err)
		//TODO
	}
	//api.sessions[sessionID] = &jsonUser // ALREADY DONE

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
func (api *ChatMe) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if api.isDebug {
		if setDebugHeaders(w, r) {
			return
		}
	}

	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, err = api.isSessionExistByValue(session.Value)
		//_, authorized = api.sessions[session.Value] // ALREADY DONE
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

// GetChats gets chats previews for user
//
// @Summary gets chats previews for user
// @ID GetChats
// @Produce json
// @Success 200 {object}  models.Response[models.Chats]
// @Failure 400 {object}  models.Response[models.Error] "Person not authorized"
// @Router /getChats [get]
func (api *ChatMe) GetChats(w http.ResponseWriter, r *http.Request) {
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
	//user := api.sessions[session.Value] // ALREADY DONE
	//user, authorized, err := api.getSessionByCookieValue(session.Value)
	authorized, err := api.isSessionExistByValue(session.Value)
	if authorized == false {
		err = models.WriteStatusJson(w, 400, models.Error{Error: "Person not authorized"})
		if err != nil {
			http.Error(w, "internal server error", 500)
			return
		}
		return
	}
	///////
	user, err := api.getUserByValue(session.Value)
	if err != nil {
		// TODO
		fmt.Println(err)
	}
	chats, err := api.getChatsByID(user.ID) // TODO
	if err != nil {
		//TODO
		fmt.Println(err)
	}
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

/*
func (api *ChatMe) getChasByID(userID uint) []*models.Chat {
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

func (api *MyHandler) getChatUsersByChatID(chatID int) []*models.ChatUser {
	usersOfChat := make([]*models.ChatUser, 0)
	for i := range api.chatUser {
		if api.chatUser[i].ChatID == chatID {
			usersOfChat = append(usersOfChat, api.chatUser[i])
		}
	}
	return usersOfChat
}*/
