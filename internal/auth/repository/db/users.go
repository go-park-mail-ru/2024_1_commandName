package db

import (
	"context"
	"database/sql"
	"errors"
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
	db           *sql.DB
	countOfUsers uint
}

func (u *Users) GetByUsername(ctx context.Context, username string) (user domain.Person, found bool) {
	fmt.Println("get by username")
	err := u.db.QueryRowContext(ctx, "SELECT id, username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetByUsername, users.go",
		}
		fmt.Println(customErr.Error())

		return user, false
	}
	return user, true
}

func (u *Users) CreateUser(ctx context.Context, user domain.Person) (userID uint, err error) {
	err = u.db.QueryRowContext(ctx, "INSERT INTO auth.person (username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id",
		user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.Avatar, user.PasswordSalt).Scan(&userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateUser, users.go",
		}
		fmt.Println(customErr.Error())

		return 0, err
	}

	u.countOfUsers++
	return userID, nil
}

func CreateFakeUsers(countOfUsers int, db *sql.DB) *sql.DB {
	fmt.Println("In create fake users.")
	counter := 0
	_ = db.QueryRow("SELECT count(id) FROM auth.person").Scan(&counter)
	if counter == 0 {
		_, err := db.Exec("ALTER SEQUENCE auth.person_id_seq RESTART WITH 1")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			fmt.Println(customErr.Error())
		}
		_, err = db.Exec("ALTER SEQUENCE auth.session_id_seq RESTART WITH 1")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			fmt.Println(customErr.Error())
		}
		_, err = db.Exec("DELETE FROM auth.person")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			fmt.Println(customErr.Error())
		}

		_, err = db.Exec("DELETE FROM auth.session")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			fmt.Println(customErr.Error())
		}

		for i := 0; i < countOfUsers; i++ {
			query := `INSERT INTO auth.person (username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			user := getFakeUser(i + 1)
			_, err := db.Exec(query, user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.Avatar, user.PasswordSalt)
			if err != nil {
				customErr := &domain.CustomError{
					Type:    "database",
					Message: err.Error(),
					Segment: "method CreateFakeUsers, users.go",
				}
				fmt.Println(customErr.Error())
			}
		}
	}
	return db
}

func getFakeUser(number int) domain.Person {
	usersHash, usersSalt := misc.GenerateHashAndSalt("testPassword!")
	testUserHash, testUserSalt := misc.GenerateHashAndSalt("Demouser123!")
	users := map[int]domain.Person{
		1: {ID: 1, Username: "IvanNaumov", Email: "ivan@mail.ru", Name: "Ivan", Surname: "Naumov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		2: {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		3: {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		4: {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		5: {ID: 5, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: usersSalt, Password: usersHash},
		6: {ID: 6, Username: "TestUser", Email: "test@mail.ru", Name: "Test", Surname: "User",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), Avatar: "avatarPath",
			PasswordSalt: testUserSalt, Password: testUserHash},
	}
	return users[number]
}

func (u *Users) GetByUserID(ctx context.Context, userID uint) (user domain.Person, found bool) {
	err := u.db.QueryRowContext(ctx, "SELECT id, username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetByUserID, users.go",
		}
		fmt.Println(customErr.Error())
		return user, false
	}
	return user, true
}

func (u *Users) UpdateUser(ctx context.Context, userUpdated domain.Person) (ok bool) {
	oldUser, found := u.GetByUserID(ctx, userUpdated.ID)
	if !found {
		return false
	}

	_, err := u.db.ExecContext(ctx, "UPDATE auth.person SET username = $1, email = $2, name = $3, surname = $4, aboat = $5, password_hash = $6, create_date = $7, lastseen_datetime = $8, avatar = $9, password_salt = $10 where id = $11",
		userUpdated.Username, userUpdated.Email, userUpdated.Name, userUpdated.Surname, userUpdated.About, userUpdated.Password, userUpdated.CreateDate, userUpdated.LastSeenDate, userUpdated.Avatar, userUpdated.PasswordSalt, oldUser.ID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateUser, users.go",
		}
		fmt.Println(customErr.Error())

		return false
	}
	return true
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

func (u *Users) GetContacts(ctx context.Context, userID uint) []domain.Person {
	contactArr := make([]domain.Person, 0)
	rows, err := u.db.QueryContext(ctx, "SELECT id, username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM chat.contacts cc JOIN auth.person ap ON cc.user1_id = ap.id WHERE cc.state = $1 and (cc.user1_id = $2 or cc.user2_id = $2)", 3, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return contactArr
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method getContact, users.go",
		}
		fmt.Println(customErr.Error())
		return contactArr
	}

	for rows.Next() {
		var userContact *domain.Person
		err = rows.Scan(&userContact.ID, &userContact.Username, &userContact.Email, &userContact.Name, &userContact.Surname, &userContact.About, &userContact.Password, &userContact.CreateDate, &userContact.LastSeenDate, &userContact.Avatar, &userContact.PasswordSalt)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method getContact, users.go",
			}
			fmt.Println(customErr.Error())
			empty := make([]domain.Person, 0)
			return empty
		}
		contactArr = append(contactArr, *userContact)
	}
	return contactArr
}

func NewUserStorage(db *sql.DB) *Users {
	return &Users{
		db:           CreateFakeUsers(6, db),
		countOfUsers: 6,
	}
}
