package repository

import (
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Users struct {
	users     map[string]domain.Person
	currentID uint
}

func (u *Users) GetByUsername(username string) (user domain.Person, found bool) {
	user, found = u.users[username]
	return user, found
}

func (u *Users) CreateUser(user domain.Person) (userID uint, err error) {
	user.ID = u.currentID
	u.currentID++
	u.users[user.Username] = user
	return user.ID, nil
}

func fillFakeUsers() map[string]domain.Person {
	usersHash, usersSalt := misc.GenerateHashAndSalt("Admin123.")
	testUserHash, testUserSalt := misc.GenerateHashAndSalt("Demouser123!")
	return map[string]domain.Person{
		"IvanNaumov": {ID: 1, Username: "IvanNaumov", Email: "ivan@mail.ru", Name: "Ivan", Surname: "Naumov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"ArtemkaChernikov": {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"ArtemZhuk": {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"AlexanderVolohov": {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"mentor": {ID: 5, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		"testUser": {ID: 6, Username: "TestUser", Email: "test@mail.ru", Name: "Test", Surname: "User",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: testUserSalt, Password: testUserHash},
	}
}

func NewUserStorage() *Users {
	return &Users{
		users:     fillFakeUsers(),
		currentID: 7,
	}
}
