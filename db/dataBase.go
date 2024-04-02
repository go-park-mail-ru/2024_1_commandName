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
		m, err := migrate.New(
			"file://migrations",
			connStrToDataBase,
		)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if err := m.Up(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("Database already up-to-date")
			} else {
				fmt.Println("Error applying migrations:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Migrations applied successfully")
		}
	*/
	return dataBase
}
