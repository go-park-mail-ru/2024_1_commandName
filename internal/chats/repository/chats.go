package repository

import (
	"ProjectMessenger/domain"
	"time"
)

type Chats struct {
	chats    map[int]domain.Chat
	chatUser []domain.ChatUser
}

func (c *Chats) GetChatsByID(userID uint) []domain.Chat {
	userChats := make(map[int]domain.Chat)
	for _, cUser := range c.chatUser {
		if cUser.UserID == userID {
			chat, ok := c.chats[cUser.ChatID]
			if ok {
				userChats[cUser.ChatID] = chat
			}
		}
	}
	var chats []domain.Chat
	for _, chat := range userChats {
		chats = append(chats, chat)
	}
	return chats
}

func (c *Chats) getChatUsersByChatID(chatID int) []*domain.ChatUser {
	usersOfChat := make([]*domain.ChatUser, 0)
	for i := range c.chatUser {
		if c.chatUser[i].ChatID == chatID {
			usersOfChat = append(usersOfChat, &c.chatUser[i])
		}
	}
	return usersOfChat
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
		&domain.Message{ID: 1, ChatID: 1, UserID: 5, Message: "Очень хороший код, ставлю 100 баллов", Edited: false, SentAt: time.Now().Add(-1 * time.Hour)},
	)

	chat1 := domain.Chat{Name: "mentors", ID: 1, Type: "group", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat1, Users: c.getChatUsersByChatID(1),
		LastMessageSentAt: messagesChat1[0].SentAt}
	c.chats[chat1.ID] = chat1

	messagesChat2 := make([]*domain.Message, 0)
	messagesChat2 = append(messagesChat2,
		&domain.Message{ID: 1, ChatID: 2, UserID: 2, Message: "Пойдём в столовку?", Edited: false, SentAt: time.Now().Add(-2 * time.Hour)},
	)
	chat2 := domain.Chat{Name: "ArtemkaChernikov", ID: 2, Type: "person", Description: "", AvatarPath: "", CreatorID: "2", Messages: messagesChat2, Users: c.getChatUsersByChatID(2),
		LastMessageSentAt: messagesChat2[0].SentAt}
	c.chats[chat2.ID] = chat2

	messagesChat3 := make([]*domain.Message, 0)
	messagesChat3 = append(messagesChat3,
		&domain.Message{ID: 1, ChatID: 3, UserID: 3, Message: "В бауманке открывают новые общаги, а Измайлово под снос", Edited: false, SentAt: time.Now().Add(-3 * time.Hour)},
	)
	chat3 := domain.Chat{Name: "Bauman News", ID: 3, Type: "channel", Description: "", AvatarPath: "", CreatorID: "3", Messages: messagesChat3, Users: c.getChatUsersByChatID(3),
		LastMessageSentAt: messagesChat3[0].SentAt}
	c.chats[chat3.ID] = chat3

	messagesChat4 := make([]*domain.Message, 0)
	messagesChat4 = append(messagesChat4,
		&domain.Message{ID: 1, ChatID: 4, UserID: 1, Message: "Ты когда тесты и авторизацию допилишь?", Edited: false, SentAt: time.Now().Add(-4 * time.Hour)},
	)
	chat4 := domain.Chat{Name: "IvanNaumov", ID: 4, Type: "person", Description: "", AvatarPath: "", CreatorID: "1", Messages: messagesChat4, Users: c.getChatUsersByChatID(4),
		LastMessageSentAt: messagesChat4[0].SentAt}
	c.chats[chat4.ID] = chat4

	messagesChat5 := make([]*domain.Message, 0)
	messagesChat5 = append(messagesChat5,
		&domain.Message{ID: 1, ChatID: 5, UserID: 4, Message: "Фронт уже готов, когда бек доделаете??", Edited: false, SentAt: time.Now().Add(-5 * time.Hour)},
	)
	chat5 := domain.Chat{Name: "AlexanderVolohov", ID: 5, Type: "person", Description: "", AvatarPath: "", CreatorID: "5", Messages: messagesChat5, Users: c.getChatUsersByChatID(5),
		LastMessageSentAt: messagesChat5[0].SentAt}
	c.chats[chat5.ID] = chat5
}

func NewChatsStorage() *Chats {
	chats := Chats{chats: make(map[int]domain.Chat), chatUser: make([]domain.ChatUser, 0)}
	chats.fillFakeChats()
	return &chats
}
