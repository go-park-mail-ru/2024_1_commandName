package auth

import (
	"time"

	"ProjectMessenger/models"
)

func (api *MyHandler) fillDB() {
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 1, UserID: 6})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 1, UserID: 5})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 2, UserID: 6})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 2, UserID: 2})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 3, UserID: 6})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 3, UserID: 3})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 4, UserID: 6})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 4, UserID: 1})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 5, UserID: 6})
	api.chatUser = append(api.chatUser, &models.ChatUser{ChatID: 5, UserID: 4})
	/////////////////////////////////////////////////////////

	messagesChat1 := make([]*models.Message, 0)
	messagesChat1 = append(messagesChat1,
		&models.Message{ID: 1, ChatID: 1, UserID: api.users["mentor"].ID, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
	)

	chat1 := models.Chat{Name: "mentor", ID: 1, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1, Users: api.getChatUsersByChatID(1)}
	api.chats[chat1.ID] = &chat1
	/////////////////////////////////////////////////////////////

	messagesChat2 := make([]*models.Message, 0)
	messagesChat2 = append(messagesChat2,
		&models.Message{ID: 1, ChatID: 2, UserID: api.users["ArtemkaChernikov"].ID, Message: "Пойдём в столовку?", Edited: false},
	)
	chat2 := models.Chat{Name: "ArtemkaChernikov", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "2", Messages: messagesChat2, Users: api.getChatUsersByChatID(2)}
	api.chats[chat2.ID] = &chat2
	////////////////////////////////////////////////////////////

	messagesChat3 := make([]*models.Message, 0)
	messagesChat3 = append(messagesChat3,
		&models.Message{ID: 1, ChatID: 3, UserID: api.users["ArtemZhuk"].ID, Message: "Ты пр уже создал? А то пора уже с мейном мерджить", Edited: false},
	)
	chat3 := models.Chat{Name: "ArtemZhuk", ID: 3, Type: "person", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat3, Users: api.getChatUsersByChatID(3)}
	api.chats[chat3.ID] = &chat3
	////////////////////////////////////////////////////////////

	messagesChat4 := make([]*models.Message, 0)
	messagesChat4 = append(messagesChat4,
		&models.Message{ID: 1, ChatID: 4, UserID: api.users["IvanNaumov"].ID, Message: "Ты когда тесты и авторизацию допилишь?", Edited: false},
	)
	chat4 := models.Chat{Name: "IvanNaumov", ID: 4, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat4, Users: api.getChatUsersByChatID(4)}
	api.chats[chat4.ID] = &chat4
	//////////////////////////////////////////////////////////////

	messagesChat5 := make([]*models.Message, 0)
	messagesChat5 = append(messagesChat5,
		&models.Message{ID: 1, ChatID: 5, UserID: api.users["AlexanderVolohov"].ID, Message: "Фронт уже готов, когда бек доделаете??", Edited: false},
	)
	chat5 := models.Chat{Name: "AlexanderVolohov", ID: 5, Type: "person", Description: "", AvatarPath: "", CreatorID: "5", Messages: messagesChat5, Users: api.getChatUsersByChatID(5)}
	api.chats[chat5.ID] = &chat5
	///////////////////////////////////////////////////////////////

}

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
}
