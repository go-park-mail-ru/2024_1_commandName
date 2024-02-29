package auth

import (
	"ProjectMessenger/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*models.Person
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*models.Person{
			"rvasily": {1, "rvasily", "vasily@mail.ru", "Vasily", "Romanov", "Developer", "", time.Now(), time.Now(), "avatarPath", "love"},
			//"aChernikov": {2, "aChernikov", "Artem557"},
		},
	}

}

// http://127.0.0.1:8080/login?login=rvasily&password=love
// http://127.0.0.1:8080/login?username=rvasily&password=love
// http://127.0.0.1:8080/register?username=Artem&password=Artem557&email=List.kedra.79

func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {

	user, ok := api.users[r.FormValue("username")]
	if !ok {
		http.Error(w, `no user`, 404)
		return
	}

	if user.Password != r.FormValue("password") {
		http.Error(w, `bad pass`, 400)
		return
	}

	SID := RandStringRunes(32)

	api.sessions[SID] = user.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte(SID))

}

func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `no sess`, 401)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, `no sess`, 401)
		return
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (api *MyHandler) Registration(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	_, ok := api.users[username]
	if ok {
		http.Error(w, `User with this username already exist. Choose another`, 404)
		return
	}
	pass := r.FormValue("password")
	email := r.FormValue("email")
	ID := uint(len(api.users) + 1)
	newUser := &models.Person{
		ID:       ID,
		Password: pass,
		Username: username,
		Email:    email,
	}
	api.users[newUser.Username] = newUser
	SID := RandStringRunes(32)

	api.sessions[SID] = newUser.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("You are registrated"))
}

func (api *MyHandler) Root(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
}

func Start() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)
	r.HandleFunc("/register", api.Registration)

	http.ListenAndServe(":8080", r)
}
