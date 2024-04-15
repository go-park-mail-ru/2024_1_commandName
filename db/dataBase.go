package db

import (
	"database/sql"
	"fmt"

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
	return dataBase
}

//GOOSE_DBSTRING=postgresql://chatme_user:EasyPassword(@127.0.0.1:8888/chatme?sslmode=disable
// goose -dir C:\Users\m2907\GolandProjects\VK_Education_Go\2024_1_commandName\db\migrations postgres "postgresql://chatme_user:EasyPassword(@127.0.0.1:8888/chatme?sslmode=disable" up
