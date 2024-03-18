package auth

import (
	"fmt"
	"strconv"
	"time"
)

func (api *ChatMe) fillDB() error {
	fmt.Println("In fillDB")

	/*
			api.createUsers(5)

			api.fillChatDB("1", "mentor", "no desc", "avatar_path", 1)
			api.fillChatDB("1", "ArtemkaChernikov", "no desc", "avatar_path", 2)
			api.fillChatDB("1", "ArtemZhuk", "no desc", "avatar_path", 3)
			api.fillChatDB("1", "IvanNaumov", "no desc", "avatar_path", 4)
			api.fillChatDB("1", "AlexanderVolohov", "no desc", "avatar_path", 5)

			query := `INSERT INTO chat.chat_user (chat_id, user_id) VALUES
		              (1, 6),
		              (1, 5),
		              (2, 6),
		              (2, 2),
		              (3, 6),
		              (3, 3),
		              (4, 6),
		              (4, 1),
		              (5, 6),
		              (5, 4)`

			_, err := api.db.Exec(query)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("for createMess Tables")

			api.createMessagesTables(5)*/

	//api.fillMessagesDB(1, 1, "Очень хороший код, ставлю 100 баллов", false)
	api.updatePasswordForExistingUsers()
	/*

		messagesChat1 := make([]*models.Message, 0)
		messagesChat1 = append(messagesChat1,
			&models.Message{ID: 1, ChatID: 1, UserID: api.users["mentor"].ID, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
		)

		chat1 := models.Chat{Name: "mentor", ID: 1, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1, Users: api.getChatUsersByChatID(1)}
		api.chats[chat1.ID] = &chat1

		messagesChat2 := make([]*models.Message, 0)
		messagesChat2 = append(messagesChat2,
			&models.Message{ID: 1, ChatID: 2, UserID: api.users["ArtemkaChernikov"].ID, Message: "Пойдём в столовку?", Edited: false},
		)
		chat2 := models.Chat{Name: "ArtemkaChernikov", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "2", Messages: messagesChat2, Users: api.getChatUsersByChatID(2)}
		api.chats[chat2.ID] = &chat2

		messagesChat3 := make([]*models.Message, 0)
		messagesChat3 = append(messagesChat3,
			&models.Message{ID: 1, ChatID: 3, UserID: api.users["ArtemZhuk"].ID, Message: "Ты пр уже создал? А то пора уже с мейном мерджить", Edited: false},
		)
		chat3 := models.Chat{Name: "ArtemZhuk", ID: 3, Type: "person", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat3, Users: api.getChatUsersByChatID(3)}
		api.chats[chat3.ID] = &chat3

		messagesChat4 := make([]*models.Message, 0)
		messagesChat4 = append(messagesChat4,
			&models.Message{ID: 1, ChatID: 4, UserID: api.users["IvanNaumov"].ID, Message: "Ты когда тесты и авторизацию допилишь?", Edited: false},
		)
		chat4 := models.Chat{Name: "IvanNaumov", ID: 4, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat4, Users: api.getChatUsersByChatID(4)}
		api.chats[chat4.ID] = &chat4

		messagesChat5 := make([]*models.Message, 0)
		messagesChat5 = append(messagesChat5,
			&models.Message{ID: 1, ChatID: 5, UserID: api.users["AlexanderVolohov"].ID, Message: "Фронт уже готов, когда бек доделаете??", Edited: false},
		)
		chat5 := models.Chat{Name: "AlexanderVolohov", ID: 5, Type: "person", Description: "", AvatarPath: "", CreatorID: "5", Messages: messagesChat5, Users: api.getChatUsersByChatID(5)}
		api.chats[chat5.ID] = &chat5*/
	return nil
}

/*
func (api *MyHandler) fillUsers() map[string]*models.Person {
	usersHash, usersSalt := generateHashAndSalt("Admin123.")
	testUserHash, testUserSalt := generateHashAndSalt("Demouser123!")
	return map[string]*models.Person{
		"IvanNaumov": {ID: 1, Username: "IvanNaumov", Email: "ivan@mail.ru", Name: "Ivan", Surname: "Naumov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"ArtemkaChernikov": {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"ArtemZhuk": {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"AlexanderVolohov": {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"mentor": {ID: 5, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"testUser": {ID: 6, Username: "TestUser", Email: "test@mail.ru", Name: "Test", Surname: "User",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: testUserSalt, Password: testUserHash},
	}
}*/

func (api *ChatMe) createMessagesTables(count int) {
	for i := 1; i < count+1; i++ {
		newTableName := "chat.messages_chat_" + strconv.Itoa(i)
		query := `CREATE TABLE IF NOT EXISTS ` + newTableName + ` AS TABLE chat.message WITH NO DATA`

		_, err := api.db.Exec(query)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (api *ChatMe) fillMessagesDB(user_id, chat_id int, message string, edited bool) {
	fmt.Println("in fillMessageDB")
	tableName := " chat.messages_chat_" + strconv.Itoa(chat_id)
	query := `INSERT INTO` + tableName + ` (user_id, chat_id, message, edited, create_datetime) VALUES ($1, $2, $3, $4, NOW())`
	//VALUES ($1, 1, 'Очень хороший код, ставлю 100 баллов', false, NOW())`

	_, err := api.db.Exec(query, user_id, chat_id, message, edited)
	if err != nil {
		fmt.Println("Error in fakeData fillMessagesDB:", err)
	}
}

func (api *ChatMe) fillChatDB(chatType, name, description, avatar_path string, creatorID int) {
	query := `INSERT INTO chat.chat (type, name, description, avatar_path, creator_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := api.db.Exec(query, chatType, name, description, avatar_path, creatorID)
	if err != nil {
		fmt.Println("Error in fakeData fillChatDB:", err)
	}
}

func (api *ChatMe) createUsers(countOfUsers int) {
	for i := 0; i < countOfUsers; i++ {
		query := `INSERT INTO auth.person (username, email) VALUES ($1, $2)`
		username := "user_" + strconv.Itoa(i+1)
		email := "email_" + strconv.Itoa(i+1)

		_, err := api.db.Exec(query, username, email)
		if err != nil {
			fmt.Println("Error in fakeData createUsers:", err)
		}

	}
}

func (api *ChatMe) updatePasswordForExistingUsers() {
	password := "testUser123!"
	query := `UPDATE auth.person SET password_hash = $1, password_salt = $2, name = $4, surname = $5, aboat = $6, create_date = $7, lastseen_datetime = $8, avatar = $9 WHERE id = $3`

	existingUsers := []struct {
		ID int
	}{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
	}

	for _, user := range existingUsers {
		passwordHash, passwordSalt := generateHashAndSalt(password)
		_, err := api.db.Exec(query, passwordHash, passwordSalt, user.ID, "artem", "chernikov", "", time.Now(), time.Now(), "")
		if err != nil {
			fmt.Println("Error in updating password:", err)
		}
	}
}
