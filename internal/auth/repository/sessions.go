package repository

import (
	"context"
	"log/slog"

	"ProjectMessenger/internal/misc"
)

type Sessions struct {
	sessions map[string]uint
}

func NewSessionStorage() *Sessions {
	return &Sessions{sessions: make(map[string]uint)}
}

func (s *Sessions) GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	userID, sessionExists = s.sessions[sessionID]
	logger.Info("GetUserIDbySessionID", "userID", userID, "sessionExists", sessionExists)
	return userID, sessionExists
}

func (s *Sessions) CreateSession(ctx context.Context, userID uint) (sessionID string) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	sessionID = misc.RandStringRunes(32)
	s.sessions[sessionID] = userID
	logger.Info("CreateSession", "sessionID", sessionID)
	return sessionID
}

func (s *Sessions) DeleteSession(ctx context.Context, sessionID string) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	delete(s.sessions, sessionID)
	logger.Info("DeleteSession", "sessionID", sessionID)
}
