package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type SessionStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool)
	CreateSession(ctx context.Context, userID uint) (sessionID string)
	DeleteSession(ctx context.Context, sessionID string)
}

type UserStore interface {
	GetByUsername(username string) (user domain.Person, found bool)
	GetByUserID(userID uint) (user domain.Person, found bool)
	CreateUser(user domain.Person) (userID uint, err error)
	UpdateUser(userUpdated domain.Person) (ok bool)
	StoreAvatar(multipartFile multipart.File, fileHandler *multipart.FileHeader) (path string, err error)
}

func CheckAuthorized(ctx context.Context, sessionID string, storage SessionStore) (authorized bool, userID uint) {
	userID, authorized = storage.GetUserIDbySessionID(ctx, sessionID)
	return authorized, userID
}

func createSession(ctx context.Context, user domain.Person, sessionStorage SessionStore) string {
	sessionID := sessionStorage.CreateSession(ctx, user.ID)
	return sessionID
}

func RegisterAndLoginUser(ctx context.Context, user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" || user.Password == "" {
		return "", fmt.Errorf("Обязательное поле не заполнено")
	}
	if !ValidatePassword(user.Password) {
		return "", fmt.Errorf("Пароль не подходит по требованиям")
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
	sessionID = createSession(ctx, user, sessionStorage)

	return sessionID, nil
}

func LoginUser(ctx context.Context, user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
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
	sessionID = createSession(ctx, userFromStorage, sessionStorage)
	return sessionID, nil
}

func LogoutUser(ctx context.Context, sessionID string, sessionStorage SessionStore) {
	sessionStorage.DeleteSession(ctx, sessionID)
	return
}

func ValidatePassword(password string) (ok bool) {
	if len([]rune(password)) < 8 {
		return false
	}

	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	digitRegex := regexp.MustCompile(`[0-9]`)
	specialCharsRegex := regexp.MustCompile(`[~!@#$%^&*_+()[\]{}></\\|"'.,:;-]`)
	allowedCharsRegex := regexp.MustCompile(`^[a-zA-Z0-9~!@#$%^&*_+()[\]{}></\\|"'.,:;-]+$`)

	if !uppercaseRegex.MatchString(password) || !lowercaseRegex.MatchString(password) ||
		!digitRegex.MatchString(password) || !specialCharsRegex.MatchString(password) ||
		!allowedCharsRegex.MatchString(password) {
		return false
	}

	return true
}
