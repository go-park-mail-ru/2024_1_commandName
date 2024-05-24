package main

import (
	"ProjectMessenger/microservices/chats_service/proto"
	contacts "ProjectMessenger/microservices/contacts_service/proto"
	session "ProjectMessenger/microservices/sessions_service/proto"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"ProjectMessenger/domain"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	authdelivery "ProjectMessenger/internal/auth/delivery"
	chatsdelivery "ProjectMessenger/internal/chats/delivery"
	messagedelivery "ProjectMessenger/internal/messages/delivery"
	"ProjectMessenger/internal/middleware"
	profiledelivery "ProjectMessenger/internal/profile/delivery"
	searchdelivery "ProjectMessenger/internal/search/delivery"
	translatedelivery "ProjectMessenger/internal/translate/delivery"

	database "ProjectMessenger/db"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	cfg := loadConfig()
	refreshIAM()
	Router(cfg)
}

func loadConfig() domain.Config {
	envPath := os.Getenv("GOCHATME_HOME")
	slog.Debug("env home =" + envPath)
	f, err := os.Open(envPath + "config.yml")
	slog.Debug("trying to open " + envPath + "config.yml")
	if err != nil {
		slog.Error("load config failed", "err", err)
		panic(err)
	}
	defer f.Close()

	var cfg domain.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func refreshIAM() {
	cmd := exec.Command("/bin/bash", "translate_key_refresh.sh")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Ошибка при выполнении скрипта:", err)
		return
	}

	fmt.Println("Bash-скрипт запушен в фоновом режиме")
}

func initializeNotifications() *firebase.App {
	envPath := os.Getenv("GOCHATME_HOME")
	opt := option.WithCredentialsFile(envPath + "chatme-45ce9-firebase-adminsdk-eo0po-6987003c9b.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		slog.Error("Error initializing firebase server app")
		return nil
	}
	return app
}

// swag init -d cmd/messenger/,domain/,internal/,microservices/chats_service/cmd,microservices/sessions_service/cmd,microservices/contacts_service/cmd

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router(cfg domain.Config) {
	router := mux.NewRouter()

	grcpSessions, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpSessions.Close()
	sessManager := session.NewAuthCheckerClient(grcpSessions)

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

	grcpContacts, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpContacts.Close()
	contactsManager := contacts.NewContactsClient(grcpContacts)

	var authHandler *authdelivery.AuthHandler
	var chatsHandler *chatsdelivery.ChatsHandler
	var profileHandler *profiledelivery.ProfileHandler
	var messageHandler *messagedelivery.MessageHandler
	var searchHandler *searchdelivery.SearchHandler
	var translateHandler *translatedelivery.TranslateHandler

	firebaseApp := initializeNotifications()

	dataBase := database.СreateDatabase()
	authHandler = authdelivery.NewAuthHandler(dataBase, sessManager, cfg.App.AvatarPath, contactsManager, firebaseApp)
	chatsHandler = chatsdelivery.NewChatsHandler(authHandler, chatsManager)
	messageHandler = messagedelivery.NewMessagesHandler(chatsHandler, dataBase)
	profileHandler = profiledelivery.NewProfileHandler(authHandler, contactsManager)
	searchHandler = searchdelivery.NewSearchHandler(chatsHandler, dataBase)
	translateHandler = translatedelivery.NewTranslateHandler(dataBase, chatsHandler)

	router.HandleFunc("/checkAuth", authHandler.CheckAuth)
	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/register", authHandler.Register)

	router.HandleFunc("/getChats", chatsHandler.GetChats)
	router.HandleFunc("/getMessages", chatsHandler.GetMessages)
	router.HandleFunc("/getChat", chatsHandler.GetChat)
	router.HandleFunc("/createPrivateChat", chatsHandler.CreatePrivateChat)
	router.HandleFunc("/createGroupChat", chatsHandler.CreateGroupChat)
	router.HandleFunc("/updateGroupChat", chatsHandler.UpdateGroupChat)
	router.HandleFunc("/deleteChat", chatsHandler.DeleteChat)
	router.HandleFunc("/getPopularChannels", chatsHandler.GetPopularChannels)
	router.HandleFunc("/createChannel", chatsHandler.CreateChannel)
	router.HandleFunc("/joinChannel", chatsHandler.JoinChannel)
	router.HandleFunc("/leaveChannel", chatsHandler.LeaveChannel)

	router.HandleFunc("/getProfileInfo", profileHandler.GetProfileInfo)
	router.HandleFunc("/updateProfileInfo", profileHandler.UpdateProfileInfo)
	router.HandleFunc("/changePassword", profileHandler.ChangePassword)

	router.HandleFunc("/uploadAvatar", profileHandler.UploadAvatar)
	router.HandleFunc("/getContacts", profileHandler.GetContacts)
	router.HandleFunc("/addContact", profileHandler.AddContact)
	router.HandleFunc("/setFirebaseToken", profileHandler.SetFirebaseToken)

	router.HandleFunc("/sendMessage", messageHandler.SendMessage)
	router.HandleFunc("/getChatMessages", messageHandler.GetChatMessages)
	router.HandleFunc("/editMessage", messageHandler.EditMessage)
	router.HandleFunc("/deleteMessage", messageHandler.DeleteMessage)

	router.HandleFunc("/search", searchHandler.SearchObjects)
	router.HandleFunc("/translate", translateHandler.TranslateMessage)

	if cfg.App.IsDebug {
		router.Use(middleware.CORS)
	}
	router.Use(middleware.AccessLogMiddleware)

	slog.Info("http server starting on " + strconv.Itoa(cfg.Server.Port))
	err = http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), router)
	if err != nil {
		slog.Error("server failed with ", "error", err)
		return
	}

}
