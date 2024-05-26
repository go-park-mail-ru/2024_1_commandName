package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"ProjectMessenger/domain"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func СreateDatabase() *sql.DB {
	connStrToDataBase := "user=chatme_user dbname=chatme password=EasyPassword( host=localhost port=8888 sslmode=disable"
	dataBase, err := sql.Open("postgres", connStrToDataBase)
	if err != nil {
		//TODO
		fmt.Println("DatBase open err:", err)
		return nil
	}

	err = dataBase.Ping()
	if err != nil {
		fmt.Println("connection to DatBase err:", err)
		return nil
	}

	/*
		driver, err := postgres.WithInstance(dataBase, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			log.Fatal(err)
		}
		//ProjectMessenger/db/migrations
		//file://migrations

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

		fmt.Println("Migration successful")*/
	numUsers := 1000
	users := generateUsers(numUsers)
	file, _ := os.Create("users.json")
	defer file.Close()
	json.NewEncoder(file).Encode(users)
	return dataBase
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateUsers(numUsers int) []domain.Person {
	users := make([]domain.Person, numUsers)
	for i := 0; i < numUsers; i++ {
		users[i] = domain.Person{
			Username: randomString(10),
			Password: randomString(10),
		}
	}
	return users
}

//GOOSE_DBSTRING=postgresql://chatme_user:EasyPassword(@127.0.0.1:8888/chatme?sslmode=disable
// goose -dir C:\Users\m2907\GolandProjects\VK_Education_Go\2024_1_commandName\db\migrations postgres "postgresql://chatme_user:EasyPassword(@127.0.0.1:8888/chatme?sslmode=disable" up
