package usecase

import "ProjectMessenger/domain"
import authusecase "ProjectMessenger/internal/auth/usecase"

func GetProfileInfo(userID uint, userStorage authusecase.UserStore) (user domain.Person, found bool) {
	user, found = userStorage.GetByUserID(userID)
	return user, found
}

func UpdateProfileInfo(userToUpdate domain.Person, numOfUpdatedFields int, userID uint, userStorage authusecase.UserStore) (err error) {
	updatedAlready := 0
	if userToUpdate.Username != "" {

	}

	return nil
}
