package _024_1_commandName

import (
	"bytes"
	"context"
	"database/sql"
	"testing"
	"time"

	"ProjectMessenger/domain"
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
		ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?").
		WithArgs(111).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("INSERT INTO auth.person (username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").
		WithArgs("test_username", "test_email@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1)) // Возвращаем результат, что одна строка была изменена

	mock.
		ExpectExec("UPDATE auth.person SET username = ?, email = ?, name = ?, surname = ?, about = ?, password_hash = ?, created_at = ?, lastseen_at = ?, avatar_path = ?, password_salt = ?").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1)) // Возвращаем результат, что одна строка была изменена

	mock.
		ExpectExec("INSERT INTO chat.contacts (user1_id, user2_id, state_id) VALUES ($1, $2, $3)").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.
		ExpectExec("INSERT INTO chat.contacts (user1_id, user2_id, state_id) VALUES ($1, $2, $3)").
		WithArgs("test_username1", "test_email1@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnError(sql.ErrTxDone)
}

func TestUserRepo_GetByUserID(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
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
}

func TestUserRepo_GetByUsername(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = ?").
		WithArgs("TestUser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(6, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	ctx := context.Background()
	user, found := userRepo.GetByUsername(ctx, "TestUser")
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

func TestUserRepo_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")

	ctx := context.Background()

	newUser := &domain.Person{
		Username:     "test_username",
		Email:        "test@example.com",
		Name:         "Test",
		Surname:      "User",
		About:        "About",
		Password:     "hashed_password",
		CreateDate:   time.Now(),
		LastSeenDate: time.Now(),
		AvatarPath:   "avatar_url",
		PasswordSalt: "password_salt",
	}

	mock.ExpectQuery("INSERT INTO auth.person (username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id").
		WithArgs("test_username", "test@example.com", "Test", "User", "About", "hashed_password", time.Now(), time.Now(), "avatar_url", "password_salt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := userRepo.CreateUser(ctx, *newUser)
	if err != nil {
		t.Error("err:", err, " ", id)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

type ReadSeekCloser struct {
	*bytes.Reader
}

func (r *ReadSeekCloser) Close() error { return nil }
