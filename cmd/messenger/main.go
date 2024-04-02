package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/echo-swagger/example/docs"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	messagedelivery "ProjectMessenger/internal/messages/delivery"
	"ProjectMessenger/internal/middleware"
	profiledelivery "ProjectMessenger/internal/profile/delivery"

	database "ProjectMessenger/db"
)

var DEBUG = false
var INMEMORY = false

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
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

	var authHandler *authdelivery.AuthHandler
	var chatsHandler *chatsdelivery.ChatsHandler
	var profileHandler *profiledelivery.ProfileHandler
	var messageHandler *messagedelivery.MessageHandler

	if INMEMORY {
		authHandler = authdelivery.NewAuthMemoryStorage()
		chatsHandler = chatsdelivery.NewChatsHandlerMemory(authHandler)
		messageHandler = messagedelivery.NewMessagesHandlerMemory(authHandler)
	} else {
		dataBase := database.СreateDatabase()
		fmt.Println(dataBase)
		authHandler = authdelivery.NewAuthHandler(dataBase)
		chatsHandler = chatsdelivery.NewChatsHandler(authHandler, dataBase)
		messageHandler = messagedelivery.NewMessagesHandler(authHandler, dataBase)
	}
	profileHandler = profiledelivery.NewProfileHandler(authHandler)

	router.HandleFunc("/checkAuth", authHandler.CheckAuth)
	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/register", authHandler.Register)
	router.HandleFunc("/getChats", chatsHandler.GetChats)
	router.HandleFunc("/getProfileInfo", profileHandler.GetProfileInfo)
	router.HandleFunc("/updateProfileInfo", profileHandler.UpdateProfileInfo)
	router.HandleFunc("/changePassword", profileHandler.ChangePassword)
	router.HandleFunc("/uploadAvatar", profileHandler.UploadAvatar)
	router.HandleFunc("/getContacts", profileHandler.GetContacts)
	router.HandleFunc("/sendMessage", messageHandler.SendMessage)
	router.HandleFunc("/getChatMessages", messageHandler.GetChatMessages)

	// middleware
	if DEBUG {
		router.Use(middleware.CORS)
	}
	router.Use(middleware.AccessLogMiddleware)

	slog.Info("http server starting on 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		slog.Error("server failed with ", "error", err)
		return
	}
}
