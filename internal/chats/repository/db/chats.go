package inMemory

import (
	"database/sql"
	"fmt"

	"ProjectMessenger/domain"
)

type Chats struct {
	db       *sql.DB
	chats    map[int]domain.Chat
	chatUser []domain.ChatUser
}

func (c *Chats) GetChatsByID(userID uint) []domain.Chat {
	fmt.Println(userID)
	chats := make([]domain.Chat, 0)
	rows, err := c.db.Query("SELECT c.* FROM chat.chat_user cu JOIN chat.chat c ON cu.chat_id = c.id WHERE cu.user_id = $1", userID)
	if err != nil {
		//TODO
		fmt.Println("err in func GetChatsByID:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var chat domain.Chat
		if err = rows.Scan(&chat.ID, &chat.Type, &chat.Name, &chat.Description, &chat.CreatorID, &chat.CreatorID); err != nil {
			//TODO
			fmt.Println("err in func GetChatsByID:", err)
			return nil
		}
		chat.Messages = c.GetMessagesByChatID(chat.ID)
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		//TODO
		fmt.Println("err in func GetChatsByID:", err)
		return nil
	}

	return chats
}

func (c *Chats) getChatUsersByChatID(chatID int) []*domain.ChatUser {
	chatUsers := make([]*domain.ChatUser, 0)
	rows, err := c.db.Query("SELECT * FROM chat.chat_user WHERE chat_id = $1", chatID)
	if err != nil {
		//TODO
		fmt.Println("err in func getChatUsersByChatID:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var chatUser domain.ChatUser
		if err = rows.Scan(&chatUser.ChatID, &chatUser.UserID); err != nil {
			//TODO
			fmt.Println("err in func getChatUsersByChatID:", err)
			return nil
		}
		chatUsers = append(chatUsers, &chatUser)
	}
	return chatUsers
}

func (c *Chats) fillFakeChats() {
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 1, UserID: 6})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 1, UserID: 5})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 2, UserID: 6})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 2, UserID: 2})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 3, UserID: 6})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 3, UserID: 3})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 4, UserID: 6})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 4, UserID: 1})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 5, UserID: 6})
	c.chatUser = append(c.chatUser, domain.ChatUser{ChatID: 5, UserID: 4})

	messagesChat1 := make([]*domain.Message, 0)
	messagesChat1 = append(messagesChat1,
		&domain.Message{ID: 1, ChatID: 1, UserID: 5, Message: "Очень хороший код, ставлю 100 баллов", Edited: false},
	)

	chat1 := domain.Chat{Name: "mentors", ID: 1, Type: "group", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1, Users: c.getChatUsersByChatID(1)}
	c.chats[chat1.ID] = chat1

	messagesChat2 := make([]*domain.Message, 0)
	messagesChat2 = append(messagesChat2,
		&domain.Message{ID: 1, ChatID: 2, UserID: 2, Message: "Пойдём в столовку?", Edited: false},
	)
	chat2 := domain.Chat{Name: "ArtemkaChernikov", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "2", Messages: messagesChat2, Users: c.getChatUsersByChatID(2)}
	c.chats[chat2.ID] = chat2

	messagesChat3 := make([]*domain.Message, 0)
	messagesChat3 = append(messagesChat3,
		&domain.Message{ID: 1, ChatID: 3, UserID: 3, Message: "В Бауманке открывают новые общаги, а Измайлово под снос", Edited: false},
	)
	chat3 := domain.Chat{Name: "Bauman News", ID: 3, Type: "channel", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat3, Users: c.getChatUsersByChatID(3)}
	c.chats[chat3.ID] = chat3

	messagesChat4 := make([]*domain.Message, 0)
	messagesChat4 = append(messagesChat4,
		&domain.Message{ID: 1, ChatID: 4, UserID: 1, Message: "Ты когда базу данных уже допилишь? Docker запустился??", Edited: false},
	)
	chat4 := domain.Chat{Name: "IvanNaumov", ID: 4, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat4, Users: c.getChatUsersByChatID(4)}
	c.chats[chat4.ID] = chat4

	messagesChat5 := make([]*domain.Message, 0)
	messagesChat5 = append(messagesChat5,
		&domain.Message{ID: 1, ChatID: 5, UserID: 4, Message: "Фронт уже готов, когда бек доделаете??", Edited: false},
	)
	chat5 := domain.Chat{Name: "AlexanderVolohov", ID: 5, Type: "person", Description: "", AvatarPath: "", CreatorID: "5", Messages: messagesChat5, Users: c.getChatUsersByChatID(5)}
	c.chats[chat5.ID] = chat5
}

func addFakeChatUsers(db *sql.DB) {
	_, err := db.Exec("DELETE FROM chat.chat_user")
	_, err = db.Exec("DELETE FROM chat.chat")
	_, err = db.Exec("DELETE FROM chat.message")
	_, err = db.Exec("ALTER SEQUENCE chat.chat_id_seq RESTART WITH 1")
	_, err = db.Exec("ALTER SEQUENCE chat.message_id_seq RESTART WITH 1")

	if err != nil {
		//TODO
		fmt.Println("Error in fakeData Chats -> addFakeChatUsers -> delete from table chat_user:", err)
	}
	query := `INSERT INTO chat.chat_user (chat_id, user_id) VALUES
		              (1, 6), 
		              (1, 5),
		              (1, 1),
		              (1, 2),
		              (1, 3),
		              (1, 4),
		              
		              (2, 6),
		              (2, 2),
		              
		              (3, 6),
		              (3, 3),
		              (3, 1),
		              (3, 2),
		              (3, 4),
		              
		              (4, 6),
		              (4, 1),
		              
		              (5, 6),
		              (5, 4)`

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("err in addFakeChatUsers:", err)
	}
}

func fillTablesMessageAndChatWithFakeData(db *sql.DB) *sql.DB {
	/*
		fillTableChatWithFakeData("2", "mentor", "no desc", "avatar_path", 1, db) // type - group
		fillTableChatWithFakeData("1", "ArtemkaChernikov", "no desc", "avatar_path", 2, db)
		fillTableChatWithFakeData("3", "ArtemZhuk", "no desc", "avatar_path", 3, db) // type - channel
		fillTableChatWithFakeData("1", "IvanNaumov", "no desc", "avatar_path", 4, db)
		fillTableChatWithFakeData("1", "AlexanderVolohov", "no desc", "avatar_path", 5, db)

		addFakeChatUsers(db)

		addFakeMessage(5, 1, "Очень хороший код, ставлю 100 баллов!", false, db)                   // mentor to group
		addFakeMessage(2, 2, "Погнали в столовку? Там солянка сейчас", false, db)                  // Chernikov to TestUser
		addFakeMessage(3, 3, "В Бауманке открывают новые общаги, а Измайлово под снос", false, db) // Zhuk to channel
		addFakeMessage(1, 4, "Ты когда базу данных уже допилишь? Docker запустился??", false, db)  // Naumov to TestUser
		addFakeMessage(4, 5, "Фронт уже готов, когда бек доделаете?", false, db)                   // Volohov to TestUser
	*/
	return db
}

func addFakeMessage(user_id, chat_id int, message string, edited bool, db *sql.DB) {
	query := `INSERT INTO chat.message (user_id, chat_id, message, edited, create_datetime) VALUES ($1, $2, $3, $4, NOW())`
	_, err := db.Exec(query, user_id, chat_id, message, edited)
	if err != nil {
		//TODO
		fmt.Println("Error in addFakeMessage:", err)
	}
}

func fillTableChatWithFakeData(chatType, name, description, avatar_path string, creatorID int, db *sql.DB) {
	query := `INSERT INTO chat.chat (type, name, description, avatar_path, creator_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query, chatType, name, description, avatar_path, creatorID)
	if err != nil {
		//TODO
		fmt.Println("Error in fakeData fillChatDB:", err)
	}
}

func (c *Chats) GetMessagesByChatID(chatID int) []*domain.Message {
	chatMessagesArr := make([]*domain.Message, 0)
	rows, err := c.db.Query("SELECT id, user_id, chat_id, message.message, edited FROM chat.message WHERE chat_id = $1", chatID)
	if err != nil {
		//TODO
		fmt.Println("err in func GetChatsByID:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var mess domain.Message
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.Edited); err != nil {
			//TODO
			fmt.Println("err in func GetMessagesByChatID:", err)
			return nil
		}
		chatMessagesArr = append(chatMessagesArr, &mess)
	}
	if err = rows.Err(); err != nil {
		//TODO
		fmt.Println("err in func GetMessagesByChatID:", err)
		return nil
	}
	return chatMessagesArr
}

func NewChatsStorage(db *sql.DB) *Chats {
	return &Chats{
		db: fillTablesMessageAndChatWithFakeData(db),
	}
}
