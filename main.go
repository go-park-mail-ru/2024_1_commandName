package main

import (
	"ProjectMessenger/auth"
	"github.com/gorilla/mux"
	"log"
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

// Router @Title Messenger authorization API
// @Version 1.0
// @BasePath /
func Router() {
	r := mux.NewRouter()

	api := auth.NewMyHandler()
	r.HandleFunc("/checkAuth", api.CheckAuth)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)
	r.HandleFunc("/register", api.Register)
	r.HandleFunc("/getChats", api.GetChats)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
		return
	}
}
