package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Sessions struct {
	db *sql.DB
}

func (s *Sessions) GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	err := s.db.QueryRowContext(ctx, "SELECT userid FROM auth.session WHERE sessionid = $1", sessionID).Scan(&userID)
	logger.Debug("GetUserIDbySessionID", "userID", userID, "sessionID", sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetUserIDbySessionID, sessions.go",
		}
		//fmt.Println(customErr.Error())
		logger.Error(customErr.Error())
		return 0, false
	}
	return userID, true
}

func (s *Sessions) CreateSession(ctx context.Context, userID uint) (sessionID string) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("CreateSession", "userID", userID)
	//fmt.Println("create session for user", userID)
	sessionID = misc.RandStringRunes(32)
	_, err := s.db.ExecContext(ctx, "INSERT INTO auth.session (sessionid, userid) VALUES ($1, $2)", sessionID, userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateSession, sessions.go",
		}
		//fmt.Println(customErr.Error())
		logger.Error(customErr.Error())
		return ""
	}
	return sessionID
}

func (s *Sessions) DeleteSession(ctx context.Context, sessionID string) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("DeleteSession", "sessionID", sessionID)
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
	slog.Info("created session storage")
	return &Sessions{db: db}
}
