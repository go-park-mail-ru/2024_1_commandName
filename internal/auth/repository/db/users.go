package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Users struct {
	db           *sql.DB
	countOfUsers uint
}

func (u *Users) GetByUsername(ctx context.Context, username string) (user domain.Person, found bool) {
	err := u.db.QueryRowContext(ctx, "SELECT * FROM auth.person WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.Avatar, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false
		}
		//TODO
		fmt.Println("Err in func GetByUsername", err)
		return user, false
	}
	fmt.Println("ID of user: ", user.ID)
	return user, true
}

func (u *Users) CreateUser(ctx context.Context, user domain.Person) (userID uint, err error) {
	err = u.db.QueryRowContext(ctx, "INSERT INTO auth.person (username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id",
		user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.Avatar, user.PasswordSalt).Scan(&userID)
	if err != nil {
		return 0, err
	}
	u.countOfUsers++
	fmt.Println(userID)

	return userID, nil
}

func CreateFakeUsers(countOfUsers int, db *sql.DB) *sql.DB {
	fmt.Println("In create fake users.")
	counter := 0
	_ = db.QueryRow("SELECT count(id) FROM auth.person").Scan(&counter)
	if counter == 0 {
		_, err := db.Exec("ALTER SEQUENCE auth.person_id_seq RESTART WITH 1")
		if err != nil {
			//TODO
			fmt.Println("Error in fakeData -> createUsers -> ALTER SEQUENCE auth.person_id_seq:", err)
		}
		_, err = db.Exec("ALTER SEQUENCE auth.session_id_seq RESTART WITH 1")
		if err != nil {
			//TODO
			fmt.Println("Error in fakeData -> createUsers -> ALTER SEQUENCE auth.session_id_seq:", err)
		}
		_, err = db.Exec("DELETE FROM auth.person")
		if err != nil {
			//TODO
			fmt.Println("Error in fakeData -> createUsers -> delete from table person:", err)
		}

		_, err = db.Exec("DELETE FROM auth.session")
		if err != nil {
			//TODO
			fmt.Println("Error in fakeData -> createUsers -> delete from table session:", err)
		}

		for i := 0; i < countOfUsers; i++ {
			query := `INSERT INTO auth.person (username, email, name, surname, aboat, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			user := getFakeUser(i + 1)
			_, err := db.Exec(query, user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.Avatar, user.PasswordSalt)
			if err != nil {
				//TODO
				fmt.Println("Error in fakeData createUsers:", err)
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

func NewUserStorage(db *sql.DB) *Users {
	return &Users{
		db:           CreateFakeUsers(6, db),
		countOfUsers: 6,
	}
}
