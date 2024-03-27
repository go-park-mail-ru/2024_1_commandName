package InMemory

import "ProjectMessenger/internal/misc"

type Sessions struct {
	sessions map[string]uint
}

func NewSessionStorage() *Sessions {
	return &Sessions{sessions: make(map[string]uint)}
}

func (s *Sessions) GetUserIDbySessionID(sessionID string) (userID uint, sessionExists bool) {
	userID, sessionExists = s.sessions[sessionID]
	return userID, sessionExists
}

func (s *Sessions) CreateSession(userID uint) (sessionID string) {
	sessionID = misc.RandStringRunes(32)
	s.sessions[sessionID] = userID
	return sessionID
}

func (s *Sessions) DeleteSession(sessionID string) {
	delete(s.sessions, sessionID)
}
