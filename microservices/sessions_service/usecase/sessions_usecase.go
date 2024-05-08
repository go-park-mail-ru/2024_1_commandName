package usecase

import (
	"ProjectMessenger/microservices/sessions_service/proto"
	"context"
)

type SessionStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool)
	CreateSession(ctx context.Context, userID uint) (sessionID string)
	DeleteSession(ctx context.Context, sessionID string)
}

type SessionManager struct {
	sessions.UnimplementedAuthCheckerServer
	storage SessionStore
}

func NewSessionManager(storage SessionStore) *SessionManager {
	return &SessionManager{storage: storage}
}

func (sm *SessionManager) CheckAuthorizedRPC(ctx context.Context, in *sessions.Session) (*sessions.UserFound, error) {
	sessionID := in.GetID()
	userID, authorized := sm.storage.GetUserIDbySessionID(ctx, sessionID)
	res := &sessions.UserFound{
		User:       &sessions.User{ID: uint64(userID)},
		Authorized: authorized,
	}
	return res, nil
}

func (sm *SessionManager) CreateSessionRPC(ctx context.Context, in *sessions.User) (*sessions.Session, error) {
	userID := in.GetID()
	sessionID := sm.storage.CreateSession(ctx, uint(userID))
	return &sessions.Session{ID: sessionID}, nil
}

func (sm *SessionManager) LogoutUserRPC(ctx context.Context, in *sessions.Session) (*sessions.ResultBool, error) {
	sessionID := in.GetID()
	sm.storage.DeleteSession(ctx, sessionID)
	return &sessions.ResultBool{Result: true}, nil
}
