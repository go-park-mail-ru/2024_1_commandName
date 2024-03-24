package usecase

import (
	"context"
	"fmt"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type SessionStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool)
	CreateSession(userID uint) (sessionID string)
	DeleteSession(ctx context.Context, sessionID string)
}

type UserStore interface {
	GetByUsername(ctx context.Context, username string) (user domain.Person, found bool)
	CreateUser(ctx context.Context, user domain.Person) (userID uint, err error)
}

func CheckAuthorized(ctx context.Context, sessionID string, storage SessionStore) (authorized bool, userID uint) {
	userID, authorized = storage.GetUserIDbySessionID(ctx, sessionID)
	return authorized, userID
}

func createSession(user domain.Person, sessionStorage SessionStore) string {
	sessionID := sessionStorage.CreateSession(user.ID)
	return sessionID
}

func RegisterAndLoginUser(ctx context.Context, user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" || user.Password == "" {
		return "", fmt.Errorf("required field is empty")
	}
	_, userFound := userStorage.GetByUsername(ctx, user.Username)
	if userFound {
		return "", fmt.Errorf("Пользователь с таким именем уже существет")
	}

	passwordHash, passwordSalt := misc.GenerateHashAndSalt(user.Password)
	user.Password = passwordHash
	user.PasswordSalt = passwordSalt

	var userID uint
	userID, err = userStorage.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	user.ID = userID
	sessionID = createSession(user, sessionStorage)

	return sessionID, nil
}

func LoginUser(ctx context.Context, user domain.Person,
	userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" {
		return "", fmt.Errorf("wrong json structure")
	}
	userFromStorage, userFound := userStorage.GetByUsername(ctx, user.Username)
	if !userFound {
		return "", fmt.Errorf("Пользователь не найден")
	}
	passwordProvided := user.Password
	passwordProvidedHash := misc.GenerateHash(passwordProvided, userFromStorage.PasswordSalt)
	if userFromStorage.Password != passwordProvidedHash {
		return "", fmt.Errorf("Неверный пароль")
	}
	sessionID = createSession(userFromStorage, sessionStorage)
	return sessionID, nil
}

func LogoutUser(ctx context.Context, sessionID string, sessionStorage SessionStore) {
	sessionStorage.DeleteSession(ctx, sessionID)
	return
}
