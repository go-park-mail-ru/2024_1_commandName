package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"ProjectMessenger/domain"
)

type Chats struct {
	db *sql.DB
}

func NewChatsStorage(db *sql.DB) *Chats {
	return &Chats{
		db: fillTablesMessageAndChatWithFakeData(db),
	}
}

func NewRawChatsStorage(db *sql.DB) *Chats {
	return &Chats{
		db: db,
	}
}

func (c *Chats) GetChatByChatID(ctx context.Context, chatID uint) (domain.Chat, error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	chat := domain.Chat{}
	err := c.db.QueryRowContext(ctx, "SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = $1", chatID).Scan(&chat.ID, &chat.Type, &chat.Name, &chat.Description, &chat.AvatarPath, &chat.CreatedAt, &chat.LastActionDateTime, &chat.CreatorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetChat didn't found chat", "chatID", chatID)
			return chat, fmt.Errorf("Чат не найден")
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetChatsForUser, chats.go",
		}
		logger.Error(err.Error(), "segment", customErr.Segment)
		fmt.Println(customErr.Error())
		return domain.Chat{}, fmt.Errorf("internal error")
	}
	//chat.Messages = c.GetMessagesByChatID(ctx, chat.ID)
	chat.Users = c.GetChatUsersByChatID(ctx, chat.ID)
	logger.Debug("GetChat: found chat", "chatID", chatID)
	return chat, nil
}

func (c *Chats) CheckPrivateChatExists(ctx context.Context, userID1, userID2 uint) (exists bool, chatID uint, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	rows, err := c.db.QueryContext(ctx, "SELECT cu1.chat_id, cu1.user_id, cu2.user_id FROM chat.chat_user cu1 INNER JOIN chat.chat_user cu2 ON cu1.chat_id = cu2.chat_id WHERE cu1.user_id = $1 AND cu2.user_id = $2 AND cu1.user_id <> cu2.user_id", userID1, userID2)
	if err != nil {
		fmt.Println("ERR:", err)
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CheckPrivateChatExists, chats.go",
		}
		fmt.Println(customErr.Error())
		return false, 0, fmt.Errorf("internal error")
	}
	for rows.Next() {
		var chatID uint
		var user1_ID uint
		var user2_ID uint
		if err = rows.Scan(&chatID, &user1_ID, &user2_ID); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CheckPrivateChatExists, chats.go",
			}
			logger.Error(customErr.Error())
			return false, 0, fmt.Errorf("internal error")
		}
		chat, err := c.GetChatByChatID(ctx, chatID)
		fmt.Println("CHAT: ", chat)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method1 CheckPrivateChatExists, chats.go",
			}
			logger.Error(customErr.Error())
			return false, 0, fmt.Errorf("internal error")
		}
		if chat.Type == "1" {
			return true, chat.ID, nil
		}
	}
	return false, 0, nil
}

func (c *Chats) CreateChat(ctx context.Context, name, description string, userIDs ...uint) (chatID uint, err error) {
	fmt.Println("Get name:", name)
	logger := slog.With("requestID", ctx.Value("traceID"))
	chatType := ""
	chatName := ""
	chatDesc := ""
	if len(userIDs) < 1 {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: "len < 1!",
			Segment: "method CreateСhat, chats.go",
		}
		logger.Error(customErr.Error())
		return 0, fmt.Errorf("internal error")
	} else if len(userIDs) == 1 {
		chatType = "3"
		chatName = name
		chatDesc = description
	} else if len(userIDs) == 2 {
		chatType = "1"
	} else if len(userIDs) > 2 {
		chatType = "2"
		chatName = name
		chatDesc = description
	}
	err = c.db.QueryRowContext(ctx, `INSERT INTO chat.chat (type_id, name, description, avatar_path, created_at,edited_at, creator_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		chatType, chatName, chatDesc, "", time.Now().UTC(), time.Now().UTC(), userIDs[0]).Scan(&chatID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateChat, chats.go",
		}
		logger.Error(customErr.Error())
		return 0, fmt.Errorf("internal error")
	}

	//firstMessageInChat := c.GetFirstChatMessageID(ctx, chatID)
	firstMessageInChat := 0

	query := `INSERT INTO chat.chat_user (chat_id, user_id, lastseen_message_id) VALUES($1, $2, $3)`
	for i := range userIDs {
		_, err = c.db.Exec(query, chatID, userIDs[i], firstMessageInChat)
		if err != nil {
			fmt.Println("here")
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateChat, chats.go",
			}
			logger.Error(customErr.Error())
			return 0, fmt.Errorf("internal error")
		}
	}
	return chatID, nil
}

func (c *Chats) DeleteChat(ctx context.Context, chatID uint) (wasDeleted bool, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	query := `DELETE FROM chat.chat_user WHERE chat_id = $1`
	res, err := c.db.Exec(query, chatID)
	if err != nil {
		logger.Error("DeleteChat: 1", "err", err.Error())
		return false, fmt.Errorf("internal error")
	}

	query = `DELETE FROM chat.message WHERE chat_id = $1`
	res, err = c.db.Exec(query, chatID)
	if err != nil {
		logger.Error("DeleteChat: 2", "err", err.Error())
		return false, fmt.Errorf("internal error")
	}

	query = `DELETE FROM chat.chat WHERE id = $1`
	res, err = c.db.Exec(query, chatID)
	if err != nil {
		logger.Error("DeleteChat: 3", "err", err.Error())
		return false, fmt.Errorf("internal error")
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("DeleteChat: 4", "err", err.Error())
		return false, fmt.Errorf("internal error")
	}

	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	} else {
		logger.Error("DeleteChat: more than one chat was deleted")
		return true, fmt.Errorf("internal error")
	}
}

func (c *Chats) GetChatsForUser(ctx context.Context, userID uint) []domain.Chat {
	fmt.Println("here")
	chats := make([]domain.Chat, 0)
	rows, err := c.db.QueryContext(ctx, "SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON cu.chat_id = c.id WHERE cu.user_id = $1", userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetChatsForUser, chats.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println("in range")
		var chat domain.Chat
		if err = rows.Scan(&chat.ID, &chat.Type, &chat.Name, &chat.Description, &chat.AvatarPath, &chat.CreatedAt, &chat.LastActionDateTime, &chat.CreatorID); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetChatsForUser, chats.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		chat.Users = c.GetChatUsersByChatID(ctx, chat.ID)
		fmt.Println("chat.ID: ", chat.ID)
		lastSeenMessageId := c.GetLastSeenMessageId(ctx, chat.ID, userID)
		chat.LastSeenMessageID = lastSeenMessageId
		if chat.Users != nil {
			chats = append(chats, chat)
		}
	}
	return chats
}

func (c *Chats) GetLastSeenMessageId(ctx context.Context, chatID uint, userID uint) (lastSeenMessageID int) {
	fmt.Println("Do for ", chatID, userID)
	err := c.db.QueryRowContext(ctx, "SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = $1 and chat_id = $2", userID, chatID).Scan(&lastSeenMessageID)
	fmt.Println("LSID:", lastSeenMessageID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetLastSeenMessageId, chats.go",
		}
		fmt.Println(customErr.Error())
	}
	return lastSeenMessageID
}

func (c *Chats) GetFirstChatMessageID(ctx context.Context, chatID uint) (firstMessageID int) {
	fmt.Println("call func with chat_id = ", chatID)
	err := c.db.QueryRowContext(ctx, "SELECT id FROM chat.message WHERE chat_id = $1 ORDER BY created_at LIMIT 1", chatID).Scan(&firstMessageID)
	fmt.Println("FMI = ", firstMessageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetFirstChatMessageID, chats.go",
		}
		fmt.Println(customErr.Error())
	}
	return firstMessageID
}

func (c *Chats) GetChatUsersByChatID(ctx context.Context, chatID uint) []*domain.ChatUser {
	chatUsers := make([]*domain.ChatUser, 0)
	rows, err := c.db.QueryContext(ctx, "SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = $1", chatID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method getChatUsersByChatID, chats.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var chatUser domain.ChatUser
		if err = rows.Scan(&chatUser.ChatID, &chatUser.UserID); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method getChatUsersByChatID, chats.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		chatUsers = append(chatUsers, &chatUser)
	}
	return chatUsers
}

func (c *Chats) GetMessagesByChatID(ctx context.Context, chatID uint) []domain.Message {
	chatMessagesArr := make([]domain.Message, 0)
	fmt.Println("in GetMessagesByChatID")
	rows, err := c.db.QueryContext(ctx, "SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = $1", chatID)
	if err != nil {
		fmt.Println("ERROR IN GetMessagesByChatID", err)
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, chats.go",
		}
		fmt.Println("error =", customErr.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var mess domain.Message
		if err = rows.Scan(&mess.ID, &mess.UserID, &mess.ChatID, &mess.Message, &mess.CreatedAt, &mess.Edited, &mess.SenderUsername); err != nil {
			fmt.Println("ERROR: ", err)
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetMessagesByChatID, chats.go",
			}
			fmt.Println(customErr.Error())
			return nil
		}
		chatMessagesArr = append(chatMessagesArr, mess)
	}
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, chats.go",
		}
		fmt.Println(customErr.Error())
		return nil
	}
	sort.Slice(chatMessagesArr, func(i, j int) bool {
		return chatMessagesArr[i].CreatedAt.Before(chatMessagesArr[j].CreatedAt)
	})
	return chatMessagesArr
}

func (c *Chats) UpdateGroupChat(ctx context.Context, updatedChat domain.Chat) (ok bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	name, desc, chatID := updatedChat.Name, updatedChat.Description, updatedChat.ID
	query := `UPDATE chat.chat SET name=$1, description=$2 WHERE id=$3`
	_, err := c.db.ExecContext(ctx, query, name, desc, chatID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, chats.go",
		}
		logger.Error(customErr.Error())
		return false
	}
	return true
}

func (c *Chats) GetNPopularChannels(ctx context.Context, userID uint, n int) ([]domain.ChannelWithCounter, error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	query := "SELECT id, name, description, creator_id, avatar_path, count(id), max(CASE WHEN user_id = $1 THEN 1 ELSE 0 END) as Aexists FROM chat.chat JOIN chat.chat_user cu on chat.id = cu.chat_id WHERE type_id = '3' GROUP BY id ORDER BY count(id) DESC LIMIT $2"
	rows, err := c.db.QueryContext(ctx, query, userID, n)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetNPopularChannels, chats.go",
		}
		logger.Error(customErr.Error())
		return []domain.ChannelWithCounter{}, fmt.Errorf("internal error")
	}
	channels := make([]domain.ChannelWithCounter, 0, n)
	defer rows.Close()
	for rows.Next() {
		var chat domain.ChannelWithCounter
		belongsToChat := 0
		if err = rows.Scan(&chat.ID, &chat.Name, &chat.Description, &chat.CreatorID, &chat.Avatar, &chat.NumOfUsers, &belongsToChat); err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetNPopularChannels, chats.go",
			}
			logger.Error(customErr.Error())
			return []domain.ChannelWithCounter{}, fmt.Errorf("internal error")
		}
		if belongsToChat == 0 {
			chat.IsMember = false
		} else {
			chat.IsMember = true
		}
		channels = append(channels, chat)
	}
	if err = rows.Err(); err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetMessagesByChatID, chats.go",
		}
		logger.Error(customErr.Error())
		return []domain.ChannelWithCounter{}, fmt.Errorf("internal error")
	}
	return channels, nil
}

func (c *Chats) AddUserToChat(ctx context.Context, userID uint, chatID uint) (err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	query := "INSERT INTO chat.chat_user (chat_id, user_id) VALUES ($1, $2)"
	_, err = c.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddUserToChat, chats.go",
		}
		logger.Error(customErr.Error())
		return fmt.Errorf("internal error")
	}
	return nil
}

func (c *Chats) RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) (err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	query := "DELETE FROM chat.chat_user WHERE chat_id = $1 AND user_id = $2"
	_, err = c.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method RemoveUserFromChat, chats.go",
		}
		logger.Error(customErr.Error())
		return fmt.Errorf("internal error")
	}
	return nil
}

func addFakeChatUsers(db *sql.DB) {
	_, err := db.Exec("DELETE FROM chat.chat_user")
	_, err = db.Exec("DELETE FROM chat.message")
	//_, err = db.Exec("ALTER SEQUENCE chat.chat_id_seq RESTART WITH 1")
	//_, err = db.Exec("ALTER SEQUENCE chat.message_id_seq RESTART WITH 1")

	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method addFakeChatUsers, chats.go",
		}
		fmt.Println(customErr.Error())
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
	counterOfRows := 0
	_ = db.QueryRow("SELECT count(id) FROM chat.chat").Scan(&counterOfRows)
	_, err = db.Exec(query)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method addFakeChatUsers, chats.go",
		}
		fmt.Println(customErr.Error())
	}
}

func fillTablesMessageAndChatWithFakeData(db *sql.DB) *sql.DB {
	fmt.Println("in db.chats")
	counterOfRows := 0
	_ = db.QueryRow("SELECT count(id) FROM chat.chat").Scan(&counterOfRows)
	if counterOfRows == 0 {
		fmt.Println("adding chats...")

		fillTableChatType(db)

		fillTableChatWithFakeData("2", "some group", "no desc", "", 1, db) // type - group
		fillTableChatWithFakeData("1", "", "no desc", "", 2, db)
		fillTableChatWithFakeData("3", "some channel", "no desc", "", 3, db) // type - channel
		fillTableChatWithFakeData("1", "", "no desc", "", 4, db)
		fillTableChatWithFakeData("1", "", "no desc", "", 5, db)

		addFakeChatUsers(db)

		addFakeMessage(5, 1, "Очень хороший код, ставлю 100 баллов!", false, db)                   // mentor to group
		addFakeMessage(2, 2, "Погнали в столовку? Там солянка сейчас", false, db)                  // Chernikov to TestUser
		addFakeMessage(3, 3, "В Бауманке открывают новые общаги, а Измайлово под снос", false, db) // Zhuk to channel
		addFakeMessage(1, 4, "Ты когда базу данных уже допилишь? Docker запустился??", false, db)  // Naumov to TestUser
		addFakeMessage(4, 5, "Фронт уже готов, когда бек доделаете?", false, db)                   // Volohov to TestUser
	}
	return db
}

func fillTableChatType(db *sql.DB) {
	_, err := db.Exec("INSERT INTO chat.chat_type (id, name) VALUES ('1', 'private'), ('2', 'group'), ('3', 'channel');")
	if err != nil {
		fmt.Println(err)
	}
}

func addFakeMessage(user_id, chat_id int, message string, edited bool, db *sql.DB) {
	query := `INSERT INTO chat.message (user_id, chat_id, message, edited, created_at) VALUES ($1, $2, $3, $4, NOW())`
	_, err := db.Exec(query, user_id, chat_id, message, edited)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method addFakeMessage, chats.go",
		}
		fmt.Println(customErr.Error())
	}
}

func fillTableChatWithFakeData(chatType, name, description, avatar_path string, creatorID int, db *sql.DB) {
	query := `INSERT INTO chat.chat (type_id, name, description, avatar_path, edited_at, creator_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, chatType, name, description, avatar_path, time.Now().UTC(), creatorID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method fillTableChatWithFakeData, chats.go",
		}
		fmt.Println(customErr.Error())
	}
}
