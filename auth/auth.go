func (m *Messenger) fillDB(users map[string]*models.Person) {
	m.chats = make(map[int]*models.Chat)
	for username, person := range users {
		fmt.Printf("Username: %s, ID: %d\n", username, person.ID)
	}
	if len(users) > 3 {
		messagesChat1 := make([]*models.Message, 0)
		messagesChat1 = append(messagesChat1,
			&models.Message{ID: 1, ChatID: 1, UserID: users["admin1"].ID, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
			&models.Message{ID: 2, ChatID: 1, UserID: users["admin2"].ID, Message: "Балдёж балдёж", Edited: false},
		)
		chat1 := models.Chat{Name: "noName", ID: 1, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1}
		m.chats[chat1.ID] = &chat1
		
		messagesChat2 := make([]*models.Message, 0)
		messagesChat2 = append(messagesChat2,
			&models.Message{ID: 1, ChatID: 2, UserID: users["admin3"].ID, Message: "Пойдём в столовку?", Edited: false},
			&models.Message{ID: 2, ChatID: 2, UserID: users["admin"].ID, Message: "Уже бегу", Edited: false},
		)
		chat2 := models.Chat{Name: "noName", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat2}
		m.chats[chat2.ID] = &chat2
		fmt.Println("MESSAGES Chat 1:")
		for _, message := range messagesChat1 {

			creatorUsername := users[findUser(message.UserID, users)]
			fmt.Printf("Message ID: %d\n", message.ID)
			fmt.Printf("Creator: %s (ID: %d)\n", creatorUsername, message.UserID)
			fmt.Printf("Message: %s\n", message.Message)
			fmt.Println("---------------------")
		}

		fmt.Println("MESSAGES Chat 2:")
		for _, message := range messagesChat2 {
			creatorUsername := users[findUser(message.UserID, users)]
			fmt.Printf("Message ID: %d\n", message.ID)
			fmt.Printf("Creator: %s (ID: %d)\n", creatorUsername, message.UserID)
			fmt.Printf("Message: %s\n", message.Message)
			fmt.Println("---------------------")
		}
	}
	fmt.Println("CHATS:", m.chats)

}

func (m *Messenger) getChats() {
	chats := m.chats
	
}

func NewMessenger() *Messenger {
	return &Messenger{
		chats: map[int]*models.Chat{},
	}
}

func findUser(ID uint, users map[string]*models.Person) string {
	for _, user := range users {
		if user.ID == ID {
			return user.Username
		}
	}
	return ""
}
