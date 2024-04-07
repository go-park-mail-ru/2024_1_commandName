package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"gopkg.in/yaml.v3"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	messagedelivery "ProjectMessenger/internal/messages/delivery"
	"ProjectMessenger/internal/middleware"
	profiledelivery "ProjectMessenger/internal/profile/delivery"

	database "ProjectMessenger/db"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	cfg := loadConfig()
	Router(cfg)
}

type config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	App struct {
		IsDebug    bool   `yaml:"isDebug"`
		InMemory   bool   `yaml:"inMemory"`
		AvatarPath string `yaml:"avatarPath"`
	} `yaml:"app"`
}

func loadConfig() config {
	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

// swag init -d cmd/messenger/,domain/,internal/

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router(cfg config) {
	router := mux.NewRouter()

	var authHandler *authdelivery.AuthHandler
	var chatsHandler *chatsdelivery.ChatsHandler
	var profileHandler *profiledelivery.ProfileHandler
	var messageHandler *messagedelivery.MessageHandler

	if cfg.App.InMemory {
		authHandler = authdelivery.NewAuthMemoryStorage()
		//chatsHandler = chatsdelivery.NewChatsHandlerMemory(authHandler)
	} else {
		dataBase := database.СreateDatabase()
		authHandler = authdelivery.NewAuthHandler(dataBase, cfg.App.AvatarPath)
		chatsHandler = chatsdelivery.NewChatsHandler(authHandler, dataBase)
		messageHandler = messagedelivery.NewMessagesHandler(authHandler, dataBase)
	}
	profileHandler = profiledelivery.NewProfileHandler(authHandler)

	router.HandleFunc("/checkAuth", authHandler.CheckAuth)
	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/register", authHandler.Register)

	router.HandleFunc("/getChats", chatsHandler.GetChats)
	router.HandleFunc("/getChat", chatsHandler.GetChat)
	router.HandleFunc("/createPrivateChat", chatsHandler.CreatePrivateChat)
	router.HandleFunc("/deleteChat", chatsHandler.DeleteChat)

	router.HandleFunc("/getProfileInfo", profileHandler.GetProfileInfo)
	router.HandleFunc("/updateProfileInfo", profileHandler.UpdateProfileInfo)
	router.HandleFunc("/changePassword", profileHandler.ChangePassword)
	router.HandleFunc("/uploadAvatar", profileHandler.UploadAvatar)
	router.HandleFunc("/getContacts", profileHandler.GetContacts)
	router.HandleFunc("/addContact", profileHandler.AddContact)

	router.HandleFunc("/sendMessage", messageHandler.SendMessage)
	router.HandleFunc("/getChatMessages", messageHandler.GetChatMessages)

	// middleware
	if cfg.App.IsDebug {
		router.Use(middleware.CORS)
	}
	router.Use(middleware.AccessLogMiddleware)

	slog.Info("http server starting on " + strconv.Itoa(cfg.Server.Port))
	err := http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), router)
	if err != nil {
		slog.Error("server failed with ", "error", err)
		return
	}
}
