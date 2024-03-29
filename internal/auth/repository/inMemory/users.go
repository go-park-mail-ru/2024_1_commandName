package inMemory

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Users struct {
	users     map[uint]domain.Person
	currentID uint
}

func (u *Users) GetContacts(ctx context.Context, userID uint) []domain.Person {
	//TODO implement me
	panic("implement me")
}

func (u *Users) UpdateUser(ctx context.Context, userUpdated domain.Person) (ok bool) {
	_, found := u.users[userUpdated.ID]
	if !found {
		return false
	}
	u.users[userUpdated.ID] = userUpdated
	return true
}

func (u *Users) GetByUserID(ctx context.Context, userID uint) (user domain.Person, found bool) {
	user, found = u.users[userID]
	return user, found
}

func (u *Users) GetByUsername(ctx context.Context, username string) (user domain.Person, found bool) {
	for _, v := range u.users {
		if v.Username == username {
			return v, true
		}
	}
	return user, found
}

func (u *Users) CreateUser(ctx context.Context, user domain.Person) (userID uint, err error) {
	user.ID = u.currentID
	u.currentID++
	u.users[user.ID] = user
	return user.ID, nil
}

func (u *Users) StoreAvatar(multipartFile multipart.File, fileHandler *multipart.FileHeader) (path string, err error) {
	originalName := fileHandler.Filename
	fileNameSlice := strings.Split(originalName, ".")
	if len(fileNameSlice) < 2 {
		return "", fmt.Errorf("Файл не имеет расширения")
	}
	extension := fileNameSlice[len(fileNameSlice)-1]
	if extension != "png" && extension != "jpg" && extension != "jpeg" && extension != "webp" && extension != "pjpeg" {
		return "", fmt.Errorf("Файл не является изображением")
	}

	//fmt.Println(extension)

	filename := misc.RandStringRunes(16)
	filePath := "./uploads/" + filename + "." + extension

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("internal error")
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, multipartFile)
	if err != nil {
		return "", fmt.Errorf("internal error")
	}

	return filePath, nil
}

func fillFakeUsers() map[uint]domain.Person {
	usersHash, usersSalt := misc.GenerateHashAndSalt("Admin123.")
	testUserHash, testUserSalt := misc.GenerateHashAndSalt("Demouser123!")
	return map[uint]domain.Person{
		1: {ID: 1, Username: "IvanNaumov", Email: "ivan@mail.ru", Name: "Ivan", Surname: "Naumov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: usersSalt, Password: usersHash},
		2: {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: usersSalt, Password: usersHash},
		3: {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: usersSalt, Password: usersHash},
		4: {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: usersSalt, Password: usersHash},
		5: {ID: 5, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: usersSalt, Password: usersHash},
		6: {ID: 6, Username: "testUser", Email: "test@mail.ru", Name: "Test", Surname: "User",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
	}
}

func NewUserStorage() *Users {
	return &Users{
		users:     fillFakeUsers(),
		currentID: 7,
	}
}
