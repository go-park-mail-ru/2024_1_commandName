package db

import (
	"database/sql"
	"errors"
	"fmt"

	"ProjectMessenger/internal/misc"
)

type Sessions struct {
	db *sql.DB
}

func (s *Sessions) GetUserIDbySessionID(sessionID string) (userID uint, sessionExists bool) {
	err := s.db.QueryRow("SELECT id FROM auth.session WHERE sessionid = $1", sessionID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false
		}
		//TODO
		fmt.Println("err in func GetUserIDbySessionID:", err)
	}
	return userID, true
}

func (s *Sessions) CreateSession(userID uint) (sessionID string) {
	fmt.Println("create session for user", userID)
	sessionID = misc.RandStringRunes(32)
	_, err := s.db.Exec("INSERT INTO auth.session (sessionid, userid) VALUES ($1, $2)", sessionID, userID)
	if err != nil {
		//TODO
		fmt.Println("err in func CreateSession:", err)
	}
	return sessionID
}

func (s *Sessions) DeleteSession(sessionID string) {
	_, err := s.db.Exec("DELETE FROM auth.session WHERE sessionID = $1", sessionID)
	if err != nil {
		//TODO
		fmt.Println("err in func DeleteSession:", err)
	}
}

func NewSessionStorage(db *sql.DB) *Sessions {
	return &Sessions{db: db}
}
