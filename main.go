package main

import "ProjectMessenger/auth"

import (
	"ProjectMessenger/models"
	"fmt"
)

func main() {
	auth.Start()
	users := make([]models.Person, 0)
	chats := make([]models.Chat, 0)
	fmt.Println(users, chats)
}
