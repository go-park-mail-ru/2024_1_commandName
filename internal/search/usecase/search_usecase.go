package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type SessionStore interface {
	GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool)
	CreateSession(ctx context.Context, userID uint) (sessionID string)
	DeleteSession(ctx context.Context, sessionID string)
}

type UserStore interface {
	GetByUserID(ctx context.Context, userID uint) (user domain.Person, found bool)
	UpdateUser(ctx context.Context, userUpdated domain.Person) (ok bool)
	StoreAvatar(ctx context.Context, multipartFile multipart.File, fileHandler *multipart.FileHeader) (path string, err error)
	GetByUsername(ctx context.Context, username string) (user domain.Person, found bool)
	CreateUser(ctx context.Context, user domain.Person) (userID uint, err error)
	GetContacts(ctx context.Context, userID uint) []domain.Person
	AddContact(ctx context.Context, userID1, userID2 uint) (ok bool)
	GetAllUserIDs(ctx context.Context) (userIDs []uint)
	GetAvatarStoragePath() string
}

func CheckAuthorized(ctx context.Context, sessionID string, storage SessionStore) (authorized bool, userID uint) {
	userID, authorized = storage.GetUserIDbySessionID(ctx, sessionID)
	return authorized, userID
}

func createSession(ctx context.Context, user domain.Person, sessionStorage SessionStore) string {
	sessionID := sessionStorage.CreateSession(ctx, user.ID)
	return sessionID
}

func RegisterAndLoginUser(ctx context.Context, user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, userID uint, err error) {
	if user.Username == "" || user.Password == "" {
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: "Обязательное поле не заполнено",
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", 0, customErr
	}
	_, userFound := userStorage.GetByUsername(ctx, user.Username)
	if !ValidatePassword(user.Password) {
		return "", 0, fmt.Errorf("Пароль не подходит по требованиям")
	}
	if userFound {
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: "Пользователь с таким именем уже существует",
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", 0, customErr
	}

	passwordHash, passwordSalt := misc.GenerateHashAndSalt(user.Password)
	user.Password = passwordHash
	user.PasswordSalt = passwordSalt
	user.CreateDate = time.Now()
	user.LastSeenDate = user.CreateDate

	userID, err = userStorage.CreateUser(ctx, user)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "userRegistration",
			Message: err.Error(),
			Segment: "method RegisterAndLoginUser, auth_usecase.go",
		}
		fmt.Println(customErr.Error())
		return "", userID, customErr
	}
	user.ID = userID
	sessionID = createSession(ctx, user, sessionStorage)

	return sessionID, userID, nil
}

func LoginUser(ctx context.Context, user domain.Person, userStorage UserStore, sessionStorage SessionStore) (sessionID string, err error) {
	if user.Username == "" {
		return "", fmt.Errorf("wrong json structure")
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
