package usecase

import "ProjectMessenger/domain"
import authusecase "ProjectMessenger/internal/auth/usecase"

func GetProfileInfo(userID uint, userStorage authusecase.UserStore) (user domain.Person, found bool) {
	user, found = userStorage.GetByUserID(userID)
	return user, found
}
