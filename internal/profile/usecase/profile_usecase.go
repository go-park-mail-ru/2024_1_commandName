package usecase

import (
	"context"
	"fmt"
	"io"
	//"io"
	"mime/multipart"
	"net/http"
	"os"

	"ProjectMessenger/microservices/contacts_service/proto"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"

	//"net/http"
	//"os"
	authusecase "ProjectMessenger/internal/auth/usecase"
)

func convertToNormalUser(person *chats.Person) domain.Person {
	return domain.Person{
		ID:           uint(person.GetID()),
		Username:     person.GetUsername(),
		Email:        person.GetEmail(),
		Name:         person.GetName(),
		Surname:      person.GetSurname(),
		About:        person.GetAbout(),
		Password:     person.GetPassword(),
		CreateDate:   person.CreateTime.AsTime(),
		LastSeenDate: person.LastSeenDate.AsTime(),
		AvatarPath:   person.GetAvatarPath(),
		PasswordSalt: person.GetPasswordSalt(),
	}
}

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
	if user.AvatarPath != "" && user.AvatarPath != "avatars/avatar.jpg" {
		oldAvatarFilename = user.AvatarPath
	}

	fileName, err := userStorage.StoreAvatar(ctx, multipartFile, fileHandler)
	if err != nil {
		return err
	}
	user.AvatarPath = "avatars/" + fileName
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

func GetContacts(ctx context.Context, userID uint, contactGRPC chats.ContactsClient) []domain.Person {
	contactsResp, err := contactGRPC.GetContacts(ctx, &chats.UserIDContacts{UserID: uint64(userID)})
	if err != nil {
		return nil
	}
	res := make([]domain.Person, 0)
	for i := range contactsResp.Persons {
		res = append(res, convertToNormalUser(contactsResp.Persons[i]))
	}
	return res
}

func AddContactByUsername(ctx context.Context, userAddingID uint, usernameToAdd string, userStorage authusecase.UserStore, contactGRPC chats.ContactsClient) (err error) {
	userToAdd, found := userStorage.GetByUsername(ctx, usernameToAdd)
	if !found {
		customErr := &domain.CustomError{
			Type:    "non type",
			Message: "Такого имени пользователя не существует",
			Segment: "method AddContactByUsername, profile_usecase.go",
		}
		return customErr
	}

	_, err = contactGRPC.AddContactByUsername(ctx, &chats.AddByUsernameReq{
		UserAddingID:  uint64(userAddingID),
		UsernameToAdd: usernameToAdd,
		UserToAddID:   uint64(userToAdd.ID),
	})
	if err != nil {
		return err
	}
	return nil
}

func AddToAllContacts(ctx context.Context, userAddingID uint, userStorage authusecase.UserStore, contactGRPC chats.ContactsClient) (ok bool) {
	userIDs := userStorage.GetAllUserIDs(ctx)
	if userIDs == nil {
		return false
	}
	usersIDsReq := make([]*chats.UserIDContacts, 0)
	for i := range userIDs {
		usersIDsReq = append(usersIDsReq, &chats.UserIDContacts{UserID: uint64(userIDs[i])})
	}

	_, err := contactGRPC.AddToAllContacts(ctx, &chats.AddToAllReq{
		Users:        &chats.UserIDArray{Users: usersIDsReq},
		UserAddingID: uint64(userAddingID),
	})
	if err != nil {
		return false
	}
	return true
}

func SetFirebaseToken(ctx context.Context, userID uint, token string, userStorage authusecase.UserStore) (err error) {
	ok := userStorage.SetFirebaseToken(ctx, userID, token)
	if !ok {
		return fmt.Errorf("internal error")
	}
	return nil
}
