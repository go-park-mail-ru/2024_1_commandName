package auth

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ProjectMessenger/domain"
	authDelivery "ProjectMessenger/internal/auth/delivery"
	database "ProjectMessenger/internal/auth/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewUserStorage(db, "")

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")
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

	userRepo := database.NewRawUserStorage(db, "")

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

	userRepo := database.NewRawUserStorage(db, "")

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

	userRepo := database.NewRawUserStorage(db, "")

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewRawUserStorage(db, "")
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

func TestUser_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/login", strings.NewReader(`{
  "password": "Demouser123!",
  "username": "TestUser"
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo")
	w := httptest.NewRecorder()
	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo").WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = ?`).
		WithArgs("TestUser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(2, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	mock.ExpectExec(`INSERT INTO auth\.session \(sessionid, userid\) VALUES (.+)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg())

	authHandler := authDelivery.NewRawAuthHandler(db, "")
	fmt.Println(authHandler, w, req)
	authHandler.Login(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUser_Login_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/login", strings.NewReader(`{
  "password": "Demouser123!",
  "username": "TestUser"
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo")
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo").WillReturnError(sql.ErrNoRows)

	authHandler := authDelivery.NewRawAuthHandler(db, "")
	authHandler.Login(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUser_Logout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/login", strings.NewReader(`{
  "password": "Demouser123!",
  "username": "TestUser"
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo")
	w := httptest.NewRecorder()
	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).
			AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo"))

	mock.ExpectExec(`DELETE FROM auth.session WHERE sessionID = ?`).
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo").
		WillReturnResult(sqlmock.NewResult(0, 1))

	authHandler := authDelivery.NewRawAuthHandler(db, "")
	authHandler.Logout(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUser_Register(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	authHandler := authDelivery.NewRawAuthHandler(db, "")
	reqBody := []byte(`{
        "username": "testuser",
        "password": "Testpassword123!",
        "email": "testuser@example.com"
    }`)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Error("expected a session cookie; got none")
	}
}

func TestUser_CheckAuth(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	authHandler := authDelivery.NewRawAuthHandler(db, "")
	reqBody := []byte(`{
        "username": "testuser",
        "password": "Testpassword123!",
        "email": "testuser@example.com"
    }`)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDoo")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.CheckAuth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Error("expected a session cookie; got none")
	}
}
