package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=artem_chernikov password=Artem557 dbname=Messenger sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	/////////////////////// CREATE PERSON
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS auth.person (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    username TEXT CHECK(length(username) <= 20),
    email TEXT CHECK(length(email) <= 30),
    name TEXT CHECK(length(name) <= 30),
    surname TEXT CHECK(length(surname) <= 20),
    aboat TEXT CHECK(length(aboat) <= 50),
    password_hash TEXT,
    create_date TIMESTAMP,
    lastseen_datetime TIMESTAMP,
    avatar TEXT          
)`)
	if err != nil {
		log.Fatal(err)
	}

	///////////////////////// CREATE CHAT
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chat.chat(
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    type varchar(1),
    name TEXT CHECK(length(name) <= 20),
    description TEXT CHECK(length(description) <= 70),
    avatar_path TEXT,
    creator_id INT REFERENCES auth.person(id)
)`)
	if err != nil {
		log.Fatal(err)
	}

	///////////////////////// CREATE CHAT_USER
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chat.chat_user (
	chat_id INT REFERENCES chat.chat(id),
	user_id INT REFERENCES auth.person(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////// CREATE MESSAGE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chat.message(
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id INT REFERENCES auth.person(id),
    chat_id INT REFERENCES chat.chat(id),
    message TEXT CHECK(length(message) <= 1000),
    edited BOOLEAN,
    create_datetime TIMESTAMP
)`)
	if err != nil {
		log.Fatal(err)
	}
}
