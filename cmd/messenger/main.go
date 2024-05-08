package main

import (
	chats "ProjectMessenger/internal/chats_service/proto"
	"ProjectMessenger/internal/sessions_service/proto"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"ProjectMessenger/domain"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/echo-swagger/example/docs"
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

// swag init -d cmd/messenger/,domain/,internal/

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router(cfg domain.Config) {
	router := mux.NewRouter()

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()
	sessManager := session.NewAuthCheckerClient(grcpConn)

	grcpConn1, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn1.Close()

	chatsManager := chats.NewChatServiceClient(grcpConn1)

	var authHandler *authdelivery.AuthHandler
	var chatsHandler *chatsdelivery.ChatsHandler
	var profileHandler *profiledelivery.ProfileHandler
	var messageHandler *messagedelivery.MessageHandler
	var searchHandler *searchdelivery.SearchHandler
	var translateHandler *translatedelivery.TranslateHandler

	dataBase := database.СreateDatabase()
	authHandler = authdelivery.NewAuthHandler(dataBase, sessManager, cfg.App.AvatarPath)
	chatsHandler = chatsdelivery.NewChatsHandler(authHandler, chatsManager)
	messageHandler = messagedelivery.NewMessagesHandler(chatsHandler, dataBase)
	profileHandler = profiledelivery.NewProfileHandler(authHandler)
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
