package _024_1_commandName

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/delivery"
	database "ProjectMessenger/internal/auth/repository/db"
	chat "ProjectMessenger/internal/chats/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

//chat "ProjectMessenger/internal/chats/repository/db"

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
		WithArgs(1, 2, 3).
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

func TestUserRepo_GetByUserID_ErrNoRows(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	_, found := userRepo.GetByUserID(ctx, 6)
	if found {
		t.Error("expected user to be not found, but it was found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetByUserID_CustomError(t *testing.T) {
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
		WillReturnError(errors.New("some database error"))

	ctx := context.Background()
	_, found := userRepo.GetByUserID(ctx, 6)
	if found {
		t.Error("expected user to be not found, but it was found")
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

func TestUserRepo_GetByUsername_ErrNoRows(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	_, found := userRepo.GetByUsername(ctx, "TestUser")
	if found {
		t.Error("expected user to be not found, but it was found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetByUsername_CustomError(t *testing.T) {
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
		WillReturnError(errors.New("some database error"))

	ctx := context.Background()
	_, found := userRepo.GetByUsername(ctx, "TestUser")
	if found {
		t.Error("expected user to be not found, but it was found")
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

	mock.ExpectQuery(`INSERT INTO auth\.person \(username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

func TestUserRepo_CreateUser_CustomError(t *testing.T) {
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

	mock.ExpectQuery(`INSERT INTO auth\.person \(username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some database error"))

	id, err := userRepo.CreateUser(ctx, *newUser)
	if err == nil {
		t.Error("err:", err, " ", id)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateUser(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("test_username", "test@example.com", "Test1", "User", "About", "hashed_password", sqlmock.AnyArg(), sqlmock.AnyArg(), "avatar_url", "password_salt", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateUser := domain.Person{
		ID:           uint(1),
		Username:     "test_username",
		Email:        "test@example.com",
		Name:         "Test1",
		Surname:      "User",
		About:        "About",
		Password:     "hashed_password",
		CreateDate:   time.Now(),
		LastSeenDate: time.Now(),
		AvatarPath:   "avatar_url",
		PasswordSalt: "password_salt",
	}

	ctx := context.Background()
	ok := userRepo.UpdateUser(ctx, updateUser)
	if !ok {
		t.Error("update failed")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateUser_UserNotFound(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных и возвращение пустого результата
	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	someUser := domain.Person{
		ID:           uint(999),
		Username:     "test_username",
		Email:        "test@example.com",
		Name:         "Test1",
		Surname:      "User",
		About:        "About",
		Password:     "hashed_password",
		CreateDate:   time.Now(),
		LastSeenDate: time.Now(),
		AvatarPath:   "avatar_url",
		PasswordSalt: "password_salt",
	}

	ctx := context.Background()
	found := userRepo.UpdateUser(ctx, someUser)

	// Проверка, что пользователь не найден
	if found {
		t.Errorf("expected user not to be found, but it was found")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateUser_CustomError(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("test_username", "test@example.com", "Test1", "User", "About", "hashed_password", sqlmock.AnyArg(), sqlmock.AnyArg(), "avatar_url", "password_salt", 1).
		WillReturnError(errors.New("some database error"))

	updateUser := domain.Person{
		ID:           uint(1),
		Username:     "test_username",
		Email:        "test@example.com",
		Name:         "Test1",
		Surname:      "User",
		About:        "About",
		Password:     "hashed_password",
		CreateDate:   time.Now(),
		LastSeenDate: time.Now(),
		AvatarPath:   "avatar_url",
		PasswordSalt: "password_salt",
	}

	ctx := context.Background()
	ok := userRepo.UpdateUser(ctx, updateUser)
	if ok {
		t.Error("expected error, got ok")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_AddContact(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")

	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)

	mock.ExpectQuery(`INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ok := userRepo.AddContact(ctx, 1, 2)
	if !ok {
		t.Error("err:", err)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_AddContact_CustomError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")

	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)

	mock.ExpectQuery(`INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some database error"))

	ok := userRepo.AddContact(ctx, 1, 2)
	if ok {
		t.Error("expected err, got ok")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_GetUserIDbySessionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)

	ctx := context.Background()

	mock.ExpectQuery("SELECT userid, sessionid FROM auth.session WHERE sessionid = ?").
		WithArgs("abcd").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).
			AddRow(123, "abcd"))

	_, sessionExists := sessionRepo.GetUserIDbySessionID(ctx, "abcd")
	if !sessionExists {
		t.Error("err: session not exist")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_GetUserIDbySessionID_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)

	ctx := context.Background()

	mock.ExpectQuery("SELECT userid, sessionid FROM auth.session WHERE sessionid = ?").
		WithArgs("abcd").
		WillReturnError(sql.ErrNoRows)

	_, sessionExists := sessionRepo.GetUserIDbySessionID(ctx, "abcd")
	if sessionExists {
		t.Error("expected false, got true")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_GetUserIDbySessionID_CustomErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)

	ctx := context.Background()

	mock.ExpectQuery("SELECT userid, sessionid FROM auth.session WHERE sessionid = ?").
		WithArgs("abcd").
		WillReturnError(errors.New("some database error"))

	_, sessionExists := sessionRepo.GetUserIDbySessionID(ctx, "abcd")
	if sessionExists {
		t.Error("expected false, got true")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_CreateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)
	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)
	//INSERT INTO auth.\session \(sessionid, userid\) VALUES (.+)
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(`INSERT INTO auth\.session \(sessionid, userid\) VALUES \(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	sessionID := sessionRepo.CreateSession(ctx, 123)
	if sessionID == "" {
		t.Error("err:", err)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_CreateSession_CustomError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)
	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)
	//INSERT INTO auth.\session \(sessionid, userid\) VALUES (.+)
	mock.ExpectExec(`INSERT INTO auth\.session \(sessionid, userid\) VALUES \(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some database error"))

	sessionID := sessionRepo.CreateSession(ctx, 123)
	if sessionID != "" {
		t.Error("err:", err)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_DeleteSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)
	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)
	//INSERT INTO auth.\session \(sessionid, userid\) VALUES (.+)
	mock.ExpectExec("DELETE FROM auth.session WHERE sessionID = ?").
		WithArgs("abcd").
		WillReturnResult(sqlmock.NewResult(1, 1))

	sessionRepo.DeleteSession(ctx, "abcd")
	if err != nil {
		t.Error("err:", err)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepo_DeleteSession_CustomError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sessionRepo := database.NewSessionStorage(db)
	ctx := context.Background()

	//INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`)
	//INSERT INTO auth.\session \(sessionid, userid\) VALUES (.+)
	mock.ExpectExec("DELETE FROM auth.session WHERE sessionID = ?").
		WithArgs("abcd").
		WillReturnError(errors.New("some database error"))

	sessionRepo.DeleteSession(ctx, "abcd")
	if err != nil {
		t.Error("err:", err)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetAllUserIDs(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных и возвращение фиктивных результатов
	mock.ExpectQuery("SELECT id FROM auth.person").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).
			AddRow(2).
			AddRow(3))

	ctx := context.Background()
	userIDs := userRepo.GetAllUserIDs(ctx)

	// Проверка, что полученные userIDs соответствуют ожидаемым
	expectedUserIDs := []uint{1, 2, 3}
	if !reflect.DeepEqual(userIDs, expectedUserIDs) {
		t.Errorf("expected userIDs to be %v, got %v", expectedUserIDs, userIDs)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetAllUserIDs_ErrNoRows(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных и возвращение ошибки sql.ErrNoRows
	mock.ExpectQuery("SELECT id FROM auth.person").
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	userIDs := userRepo.GetAllUserIDs(ctx)

	// Проверка, что возвращается пустой список userIDs
	if len(userIDs) != 0 {
		t.Errorf("expected userIDs to be empty, got %v", userIDs)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetAllUserIDs_CustomError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")

	mock.ExpectQuery("SELECT id FROM auth.person").
		WillReturnError(errors.New("some database error"))

	ctx := context.Background()
	userIDs := userRepo.GetAllUserIDs(ctx)

	if len(userIDs) != 0 {
		t.Errorf("expected userIDs to be empty, got %v", userIDs)
	}

	// Проверка, что ожидаемый CustomError создается и передается в logger.Error
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetContacts_Success(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON").
		WithArgs(123, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(1, "username1", "email1", "name1", "surname1", "about1", fixedTime, "avatar1").
			AddRow(2, "username2", "email2", "name2", "surname2", "about2", fixedTime, "avatar2"))

	ctx := context.Background()
	contacts := userRepo.GetContacts(ctx, 123)

	// Проверка, что получены ожидаемые контакты
	expectedContacts := []domain.Person{
		{ID: 1, Username: "username1", Email: "email1", Name: "name1", Surname: "surname1", About: "about1", AvatarPath: "avatar1"},
		{ID: 2, Username: "username2", Email: "email2", Name: "name2", Surname: "surname2", About: "about2", AvatarPath: "avatar2"},
	}
	if !reflect.DeepEqual(contacts, expectedContacts) {
		t.Errorf("expected contacts to be %v, got %v", expectedContacts, contacts)
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetContacts_ErrNoRows(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON").
		WithArgs(123, 3).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	contacts := userRepo.GetContacts(ctx, 123)

	// Проверка, что получены ожидаемые контакты
	if len(contacts) != 0 {
		t.Errorf("expected len = 0")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetContacts_CustomError(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	userRepo := database.NewRawUserStorage(db, "")

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON").
		WithArgs(123, 3).
		WillReturnError(errors.New("some database error"))

	ctx := context.Background()
	contacts := userRepo.GetContacts(ctx, 123)

	// Проверка, что получены ожидаемые контакты
	if len(contacts) != 0 {
		t.Errorf("expected len = 0")
	}

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_GetChatByChatID_Succes(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "1", "test@mail.ru", "Test", "User", time.Now(), time.Now(), 1))

	ctx := context.Background()
	chat1, err := chatRepo.GetChatByChatID(ctx, 1)
	if err != nil {
		t.Error("err:", err)
	}
	fmt.Println(chat1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_GetChatByChatID_ErrNoRows(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	chat1, err := chatRepo.GetChatByChatID(ctx, 1)
	if err == nil {
		t.Error("expected err!")
	}
	fmt.Println(chat1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_GetChatByChatID_CustomError(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some database error"))

	ctx := context.Background()
	chat1, err := chatRepo.GetChatByChatID(ctx, 1)
	if err == nil {
		t.Error("expected err!")
	}
	fmt.Println(chat1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_CreateChat(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery(`INSERT INTO chat\.chat \(type_id, name, description, avatar_path, created_at,edited_at, creator_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	fmt.Println("here")
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	ctx := context.Background()
	arrOfUserIDs := make([]uint, 0)
	arrOfUserIDs = append(arrOfUserIDs, uint(1))
	arrOfUserIDs = append(arrOfUserIDs, uint(2))
	_, err = chatRepo.CreateChat(ctx, "chat1", "desc", arrOfUserIDs...)
	if err != nil {
		t.Error("err:", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
