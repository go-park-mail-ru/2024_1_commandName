package main

import (
	"ProjectMessenger/auth"
	"github.com/gorilla/mux"
	"net/http"
)

import (
	"ProjectMessenger/models"
	"fmt"
)

func main() {
	Router()
	users := make([]models.Person, 0)
	chats := make([]models.Chat, 0)
	fmt.Println(users, chats)
}

// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router() {
	r := mux.NewRouter()

	api := auth.NewMyHandler()
	r.HandleFunc("/checkAuth", api.CheckAuth)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)
	r.HandleFunc("/register", api.Register)

	http.ListenAndServe(":8080", r)
}
