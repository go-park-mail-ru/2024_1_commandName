package main

import "ProjectMessenger/auth"

import (
	"ProjectMessenger/models"
	"fmt"
)

// @Title Messenger authorization API
// @Version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	auth.Start()
	users := make([]models.User, 0)
	chats := make([]models.Chat, 0)
	fmt.Println(users, chats)
}
