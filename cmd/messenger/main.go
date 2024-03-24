package main

import (
	"fmt"
	"log"
	"net/http"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/middleware"

	database "ProjectMessenger/db"

	"github.com/gorilla/mux"
)

var DEBUG = true

func main() {
	Router()
}

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router() {
	router := mux.NewRouter()

	dataBase := database.СreateDatabase()

	authHandler := authdelivery.NewAuthHandler(dataBase)
	chatsHandler := chatsdelivery.NewChatsHandler(authHandler, dataBase)

	router.HandleFunc("/checkAuth", authHandler.CheckAuth)
	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/register", authHandler.Register)
	router.HandleFunc("/getChats", chatsHandler.GetChats)

	// middleware
	if DEBUG {
		router.Use(middleware.CORS)
	}

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("err")
		log.Fatal(err)
		return
	}
}
