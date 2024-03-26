package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Sessions struct {
	db *sql.DB
}

func (s *Sessions) GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool) {
	err := s.db.QueryRowContext(ctx, "SELECT id FROM auth.session WHERE sessionid = $1", sessionID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetUserIDbySessionID, sessions.go",
		}
		fmt.Println(customErr.Error())
		return 0, false
	}
	return userID, true
}

func (s *Sessions) CreateSession(userID uint) (sessionID string) {
	fmt.Println("create session for user", userID)
	sessionID = misc.RandStringRunes(32)
	_, err := s.db.Exec("INSERT INTO auth.session (sessionid, userid) VALUES ($1, $2)", sessionID, userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateSession, sessions.go",
		}
		fmt.Println(customErr.Error())
		return ""
	}
	return sessionID
}

func (s *Sessions) DeleteSession(ctx context.Context, sessionID string) {
	_, err := s.db.ExecContext(ctx, "DELETE FROM auth.session WHERE sessionID = $1", sessionID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method DeleteSession, sessions.go",
		}
		fmt.Println(customErr.Error())
	}
}

func NewSessionStorage(db *sql.DB) *Sessions {
	return &Sessions{db: db}
}
