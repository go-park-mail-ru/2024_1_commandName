package usecase

import (
	"ProjectMessenger/internal/sessions_service/proto"
	"context"
)

type SessionStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool)
	CreateSession(ctx context.Context, userID uint) (sessionID string)
	DeleteSession(ctx context.Context, sessionID string)
}

type SessionManager struct {
	session.UnimplementedAuthCheckerServer
	storage SessionStore
}

func NewSessionManager(storage SessionStore) *SessionManager {
	return &SessionManager{storage: storage}
}

func (sm *SessionManager) CheckAuthorizedRPC(ctx context.Context, in *session.Session) (*session.UserFound, error) {
	sessionID := in.GetID()
	userID, authorized := sm.storage.GetUserIDbySessionID(ctx, sessionID)
	res := &session.UserFound{
		User:       &session.User{ID: uint64(userID)},
		Authorized: authorized,
	}
	return res, nil
}

func (sm *SessionManager) CreateSessionRPC(ctx context.Context, in *session.User) (*session.Session, error) {
	userID := in.GetID()
	sessionID := sm.storage.CreateSession(ctx, uint(userID))
	return &session.Session{ID: sessionID}, nil
}

func (sm *SessionManager) LogoutUserRPC(ctx context.Context, in *session.Session) (*session.ResultBool, error) {
	sessionID := in.GetID()
	sm.storage.DeleteSession(ctx, sessionID)
	return &session.ResultBool{Result: true}, nil
}
