package usecase

import (
	"fmt"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type SessionStore interface {
	GetUserIDbySessionID(sessionID string) (userID uint, sessionExists bool)
	CreateSession(userID uint) (sessionID string)
	DeleteSession(sessionID string)
}

type UserStore interface {
	GetByUsername(username string) (user domain.Person, found bool)
	CreateUser(user domain.Person) (userID uint, err error)
}

func CheckAuthorized(sessionID string, storage SessionStore) (authorized bool, userID uint) {
	userID, authorized = storage.GetUserIDbySessionID(sessionID)
	return authorized, userID
}

func createSession(user domain.Person, sessionStorage SessionStore) string {
	sessionID := sessionStorage.CreateSession(user.ID)
	return sessionID
}

func RegisterAndLoginUser(user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" || user.Password == "" {
		return "", fmt.Errorf("required field is empty")
	}
	_, userFound := userStorage.GetByUsername(user.Username)
	if userFound {
		return "", fmt.Errorf("Пользователь с таким именем уже существет")
	}

	passwordHash, passwordSalt := misc.GenerateHashAndSalt(user.Password)
	user.Password = passwordHash
	user.PasswordSalt = passwordSalt

	_, err = userStorage.CreateUser(user)
	if err != nil {
		return "", err
	}
	sessionID = createSession(user, sessionStorage)

	return sessionID, nil
}

func LoginUser(user domain.Person,
	userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" {
		return "", fmt.Errorf("wrong json structure")
	}
	userFromStorage, userFound := userStorage.GetByUsername(user.Username)
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

func LogoutUser(sessionID string, sessionStorage SessionStore) {
	sessionStorage.DeleteSession(sessionID)
	return
}
