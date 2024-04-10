package _024_1_commandName

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	//"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/delivery"
	database "ProjectMessenger/internal/auth/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewUserRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	defer db.Close()
	userRepo := delivery.NewAuthHandler(db, "")
	if userRepo != nil {
		return
	}
}

func TestUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	//userStorage := database.NewUserStorage(db, "")

	mock.
		ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE id = ?").
		WithArgs(111).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("INSERT INTO auth.person (username, email, name, surname, about, password_hash, create_date, lastseen_datetime, avatar, password_salt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
		WithArgs("test_username", "test_email@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1)) // Возвращаем результат, что одна строка была изменена

	mock.
		ExpectExec("UPDATE auth.person SET username = ?, email = ?, name = ?, surname = ?, about = ?, password_hash = ?, create_date = ?, lastseen_datetime = ?, avatar = ?, password_salt = ?").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1)) // Возвращаем результат, что одна строка была изменена

	mock.
		ExpectExec("INSERT INTO chat.contacts (user1_id, user2_id, state) VALUES ($1, $2, $3)").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.
		ExpectExec("INSERT INTO chat.contacts (user1_id, user2_id, state) VALUES ($1, $2, $3)").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnError(sql.ErrTxDone)
}

func TestUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewUserStorage(db, "")

	//test #1
	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE id = ?").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "create_date", "lastseen_datetime", "avatar", "password_salt"}).
			AddRow(6, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	ctx := context.Background()
	user, found := userRepo.GetByUserID(ctx, 6)
	if !found {
		t.Error("expected user to be found, but it was not found")
	}
	if user.ID != 6 {
		t.Errorf("expected user ID to be 6, got %d", user.ID)
	}
	if user.Username != "TestUser" {
		t.Errorf("expected username to be test_username, got %s", user.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//test #2
	/*
		mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE id = ?").
			WithArgs(6).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "create_date", "lastseen_datetime", "avatar", "password_salt"}).
				AddRow(6, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

		ctx = context.Background()
		user = &domain.Person{
			ID:       8,
			Username: "New username",
		}
		found := userRepo.UpdateUser(ctx, user)
		if !found {
			t.Error("expected user to be found, but it was not found")
		}
		if user.ID != 6 {
			t.Errorf("expected user ID to be 6, got %d", user.ID)
		}
		if user.Username != "TestUser" {
			t.Errorf("expected username to be test_username, got %s", user.Username)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}*/
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewUserStorage(db, "")

	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, create_date, lastseen_datetime, avatar, password_salt FROM auth.person WHERE id = ?").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "create_date", "lastseen_datetime", "avatar", "password_salt"}).
			AddRow(6, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	ctx := context.Background()
	user, found := userRepo.GetByUserID(ctx, 6)
	fmt.Println(found)
	if !found {
		t.Error("expected user to be found, but it was not found")
	}

	if user.ID != 6 {
		t.Errorf("expected user ID to be 6, got %d", user.ID)
	}
	if user.Username != "TestUser" {
		t.Errorf("expected username to be test_username, got %s", user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
