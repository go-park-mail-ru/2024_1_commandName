package main

import (
	"fmt"
	"log"
	"net/http"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	"ProjectMessenger/internal/middleware"
	profiledelivery "ProjectMessenger/internal/profile/delivery"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

var DEBUG = false

func main() {
	Router()
}

// swag init -d cmd/messenger/,domain/,internal/

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router() {
	router := mux.NewRouter()
	authHandler := authdelivery.NewAuthHandler()
	chatsHandler := chatsdelivery.NewChatsHandler(authHandler)
	profileHandler := profiledelivery.NewProfileHandler(authHandler)

	router.HandleFunc("/checkAuth", authHandler.CheckAuth)
	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/register", authHandler.Register)
	router.HandleFunc("/getChats", chatsHandler.GetChats)
	router.HandleFunc("/getProfileInfo", profileHandler.GetProfileInfo)
	router.HandleFunc("/updateProfileInfo", profileHandler.UpdateProfileInfo)

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
