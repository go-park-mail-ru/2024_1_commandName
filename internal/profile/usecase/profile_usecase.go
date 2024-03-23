package usecase

import (
	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
	"fmt"
	"mime/multipart"
	"net/http"
)
import authusecase "ProjectMessenger/internal/auth/usecase"

func GetProfileInfo(userID uint, userStorage authusecase.UserStore) (user domain.Person, found bool) {
	user, found = userStorage.GetByUserID(userID)
	return user, found
}

func UpdateProfileInfo(updatedFields domain.Person, numOfUpdatedFields int, userID uint, userStorage authusecase.UserStore) (err error) {
	userFromStorage, found := userStorage.GetByUserID(userID)
	if !found {
		return fmt.Errorf("user not found")
	}
	numOfUpdatedAlready := 0
	if updatedFields.Username != "" {
		numOfUpdatedAlready++
		userFromStorage.Username = updatedFields.Username
	}
	if updatedFields.Email != "" {
		numOfUpdatedAlready++
		userFromStorage.Email = updatedFields.Email
	}
	if updatedFields.Name != "" {
		numOfUpdatedAlready++
		userFromStorage.Name = updatedFields.Name
	}
	if updatedFields.Surname != "" {
		numOfUpdatedAlready++
		userFromStorage.Surname = updatedFields.Surname
	}
	if updatedFields.About != "" {
		numOfUpdatedAlready++
		userFromStorage.About = updatedFields.About
	}

	if numOfUpdatedFields != numOfUpdatedAlready {
		return fmt.Errorf("number of update fields is mismatched")
	}

	ok := userStorage.UpdateUser(userFromStorage)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}

func ChangePassword(oldPassword string, newPassword string, userID uint, userStorage authusecase.UserStore) (err error) {
	userFromStorage, found := userStorage.GetByUserID(userID)
	if !found {
		return fmt.Errorf("user not found")
	}
	if !authusecase.ValidatePassword(newPassword) {
		return fmt.Errorf("Новый пароль не подходит по требованиям")
	}
	storagePasswordHash := userFromStorage.Password
	oldPasswordHash := misc.GenerateHash(oldPassword, userFromStorage.PasswordSalt)
	if storagePasswordHash != oldPasswordHash {
		return fmt.Errorf("Старый пароль введён неверно")
	}

	newPasswordHash, newPasswordSalt := misc.GenerateHashAndSalt(newPassword)
	userFromStorage.Password = newPasswordHash
	userFromStorage.PasswordSalt = newPasswordSalt

	ok := userStorage.UpdateUser(userFromStorage)
	if !ok {
		return fmt.Errorf("error updating password")
	}
	return nil
}

func ChangeAvatar(multipartFile multipart.File, fileHandler *multipart.FileHeader, userID uint, userStorage authusecase.UserStore) (err error) {
	user, found := userStorage.GetByUserID(userID)
	if !found {
		return fmt.Errorf("internal error")
	}
	buff := make([]byte, 512)
	if _, err = multipartFile.Read(buff); err != nil {
		return fmt.Errorf("internal error")
	}
	mimeType := http.DetectContentType(buff)
	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/pjpeg" && mimeType != "image/webp" {
		return fmt.Errorf("Файл не является изображением")
	}

	path, err := userStorage.StoreAvatar(multipartFile, fileHandler)
	if err != nil {
		return err
	}
	user.Avatar = path
	ok := userStorage.UpdateUser(user)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}
