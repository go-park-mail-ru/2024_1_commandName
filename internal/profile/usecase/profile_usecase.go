package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	//"io"
	"mime/multipart"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"

	//"net/http"
	//"os"
	authusecase "ProjectMessenger/internal/auth/usecase"
)

func GetProfileInfo(ctx context.Context, userID uint, userStorage authusecase.UserStore) (user domain.Person, found bool) {
	user, found = userStorage.GetByUserID(ctx, userID)
	return user, found
}

func UpdateProfileInfo(ctx context.Context, updatedFields domain.Person, numOfUpdatedFields int, userID uint, userStorage authusecase.UserStore) (err error) {
	userFromStorage, found := userStorage.GetByUserID(ctx, userID)
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

	ok := userStorage.UpdateUser(ctx, userFromStorage)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}

func ChangePassword(ctx context.Context, oldPassword string, newPassword string, userID uint, userStorage authusecase.UserStore) (err error) {
	userFromStorage, found := userStorage.GetByUserID(ctx, userID)
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

	ok := userStorage.UpdateUser(ctx, userFromStorage)
	if !ok {
		return fmt.Errorf("error updating password")
	}
	return nil
}

func ChangeAvatar(ctx context.Context, multipartFile multipart.File, fileHandler *multipart.FileHeader, userID uint, userStorage authusecase.UserStore) (err error) {
	buff := make([]byte, 512)
	if _, err = multipartFile.Read(buff); err != nil {
		return fmt.Errorf("internal error")
	}
	seek, err := multipartFile.Seek(0, io.SeekStart)
	if err != nil || seek != 0 {
		return fmt.Errorf("internal error")
	}
	mimeType := http.DetectContentType(buff)
	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/pjpeg" && mimeType != "image/webp" {
		return fmt.Errorf("Файл не является изображением")
	}

	user, found := userStorage.GetByUserID(ctx, userID)
	if !found {
		return fmt.Errorf("internal error")
	}
	oldAvatarFilename := ""
	if user.Avatar != "" {
		oldAvatarFilename = user.Avatar
	}

	fileName, err := userStorage.StoreAvatar(ctx, multipartFile, fileHandler)
	if err != nil {
		return err
	}
	user.Avatar = "avatars/" + fileName
	ok := userStorage.UpdateUser(ctx, user)
	if !ok {
		return fmt.Errorf("internal error")
	}

	if oldAvatarFilename != "" {
		err = os.Remove(userStorage.GetAvatarStoragePath() + oldAvatarFilename)
		if err != nil {
			return fmt.Errorf("internal error")
		}
	}
	return nil
}

func GetContacts(ctx context.Context, userID uint, userStorage authusecase.UserStore) []domain.Person {
	contacts := userStorage.GetContacts(ctx, userID)
	return contacts
}

func AddContactByUsername(ctx context.Context, userAddingID uint, usernameToAdd string, userStorage authusecase.UserStore) (err error) {
	userToAdd, found := userStorage.GetByUsername(ctx, usernameToAdd)
	if !found {
		return fmt.Errorf("Такого имени пользователя не существует")
	}
	contacts := GetContacts(ctx, userAddingID, userStorage)
	for i := range contacts {
		if contacts[i].Username == usernameToAdd {
			return fmt.Errorf("Такой контакт уже существует")
		}
	}

	ok := userStorage.AddContact(ctx, userAddingID, userToAdd.ID)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}

func AddToAllContacts(ctx context.Context, userAddingID uint, userStorage authusecase.UserStore) (ok bool) {
	userIDs := userStorage.GetAllUserIDs(ctx)
	if userIDs == nil {
		return false
	}
	for i := range userIDs {
		if userIDs[i] == userAddingID {
			continue
		}
		if !userStorage.AddContact(ctx, userAddingID, userIDs[i]) {
			return false
		}
	}
	return true
}
