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
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: "required field is empty",
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}
	_, userFound := userStorage.GetByUsername(ctx, user.Username)
	if userFound {
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: "Пользователь с таким именем уже существует",
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}

	passwordHash, passwordSalt := misc.GenerateHashAndSalt(user.Password)
	user.Password = passwordHash
	user.PasswordSalt = passwordSalt

	var userID uint
	userID, err = userStorage.CreateUser(ctx, user)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: err.Error(),
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}
	user.ID = userID
	sessionID = createSession(user, sessionStorage)

	return sessionID, nil
}

func LoginUser(ctx context.Context, user domain.Person,
	userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" {
		customErr := &domain.CustomError{
			Type:    "userLogin",
			Message: "wrong json structure",
			Segment: "method LoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}

	userFromStorage, userFound := userStorage.GetByUsername(ctx, user.Username)
	if !userFound {
		customErr := &domain.CustomError{
			Type:    "userLogin",
			Message: "Пользователь не найден",
			Segment: "method LoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}
	passwordProvided := user.Password
	passwordProvidedHash := misc.GenerateHash(passwordProvided, userFromStorage.PasswordSalt)
	if userFromStorage.Password != passwordProvidedHash {
		customErr := &domain.CustomError{
			Type:    "userLogin",
			Message: "Неверный пароль",
			Segment: "method LoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", customErr
	}
	sessionID = createSession(userFromStorage, sessionStorage)
	return sessionID, nil
}

func LogoutUser(ctx context.Context, sessionID string, sessionStorage SessionStore) {
	sessionStorage.DeleteSession(ctx, sessionID)
	return
}
