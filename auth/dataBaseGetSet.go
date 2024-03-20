package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"ProjectMessenger/models"
)

func (api *ChatMe) getUserByValue(sessionValue string) (models.Person, error) {
	var user models.Person
	id := 0
	sid := ""
	userID := 0
	err := api.db.QueryRow("SELECT * FROM auth.session WHERE sessionid = $1", sessionValue).Scan(&id, &sid, &userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user`s session not found")
		}
		return user, err
	}
	user.ID = uint(userID)
	err = api.db.QueryRow("SELECT * FROM auth.person WHERE ID = $1", user.ID).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (api *ChatMe) getSessioByCookieValue(sessionValue string) (models.Person, bool, error) {
	fmt.Println("Value = ", sessionValue)
	var user models.Person
	var userID uint
	id := 0
	sid := ""
	err := api.db.QueryRow("SELECT * FROM auth.session WHERE sessionid = $1", sessionValue).Scan(&id, &sid, &userID)
	fmt.Println(id, sid, userID)
	fmt.Println("error = ", err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, errors.New("Session not found")
		}
		return user, false, err
	}
	user.ID = userID
	fmt.Println("params: ", id, sid, userID)
	err = api.db.QueryRow("SELECT * FROM auth.person WHERE ID = $1", user.ID).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, errors.New("user not found")
		}
		return user, false, err
	}
	return user, false, nil
}

func (api *ChatMe) isSessionExistByValue(value string) (bool, error) {
	sessionCounter := 0
	err := api.db.QueryRow("SELECT COUNT(*) FROM auth.session WHERE sessionid = $1", value).Scan(&sessionCounter)
	if err != nil {
		return true, err
	}
	if sessionCounter > 0 {
		return true, nil
	}
	return false, nil
}

func (api *ChatMe) getUserByUsername(username string) (models.Person, bool, error) {
	fmt.Println("In get user")
	var user models.Person
	counter := 0
	err := api.db.QueryRow("SELECT COUNT(*) FROM auth.person WHERE username = $1", username).Scan(&counter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, err
		}
		fmt.Println(err)
		return user, false, err
	}
	if counter == 0 {
		return user, false, nil
	}
	err = api.db.QueryRow("SELECT * FROM auth.person WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, err
		}
		fmt.Println(err)
		return user, true, err
	}

	return user, true, nil
}

func (api *ChatMe) deleteSessionByCookieValue(value string) error {
	_, err := api.db.Exec("DELETE FROM auth.session WHERE sessionID = $1", value)
	fmt.Println("in delete", value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (api *ChatMe) getCountOfUsers() (int, error) {
	count := 0
	err := api.db.QueryRow("SELECT COUNT(*) FROM auth.person").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (api *ChatMe) getChatsByID(userID uint) ([]*models.Chat, error) {
	/*
		userChats := make([]*models.Chat, 0)
		err := api.db.QueryRow("SELECT * FROM chat.chat_user WHERE user_id = $1", userID).Scan(&userChats)
		if err != nil {
			return nil, err
		}
		return userChats, nil*/
	rows, err := api.db.Query("SELECT chat_id FROM chat.chat_user WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIDs []int
	for rows.Next() {
		var chatID int
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var chats []*models.Chat
	for _, chatID := range chatIDs {
		chat, err := api.getChatByChatID(chatID)
		if err == nil {
			chats = append(chats, chat)
		}
	}
	return chats, nil
}

func (api *ChatMe) getChatByChatID(chatID int) (*models.Chat, error) {
	var chat models.Chat
	err := api.db.QueryRow("SELECT * FROM chat.chat_user WHERE chat_id = $1", chatID).Scan(&chat)
	if errors.Is(err, sql.ErrNoRows) {
		return &chat, errors.New("Session not found")
	}
	return &chat, nil
}

func (api *ChatMe) getChatUsersByChatID(chatID int) ([]*models.ChatUser, error) {
	rows, err := api.db.Query("SELECT * FROM chat.chat_user WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersOfChat []*models.ChatUser
	for rows.Next() {
		var chatUser models.ChatUser
		if err = rows.Scan(&chatUser.UserID, &chatUser.ChatID); err != nil {
			return nil, err
		}
		usersOfChat = append(usersOfChat, &chatUser)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usersOfChat, nil
}

func (api *ChatMe) setSessionBySessionID(sessionID string, user models.Person) error {
	fmt.Println("in set session")
	fmt.Println(sessionID, user.ID)
	_, err := api.db.Exec("INSERT INTO auth.session (sessionid, userid) VALUES ($1, $2)", sessionID, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (api *ChatMe) setUserByUsername(user models.Person) error {
	fmt.Println("User in setUser = ", user)
	_, err := api.db.Exec("INSERT INTO auth.person (username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.Avatar, user.PasswordSalt)
	if err != nil {
		return err
	}
	return nil
}

func (api *ChatMe) checkForExistence(request string) (int, error) {
	elemCounter := 0
	err := api.db.QueryRow(request).Scan(&elemCounter)
	if err != nil {
		return 0, err
	}
	return elemCounter, err
}
