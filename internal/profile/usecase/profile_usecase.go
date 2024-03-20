package usecase

import (
	"ProjectMessenger/domain"
	"fmt"
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
