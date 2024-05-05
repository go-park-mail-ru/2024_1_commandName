package auth

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"ProjectMessenger/domain"
	database "ProjectMessenger/internal/auth/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

//go test -coverpkg=./... -coverprofile=cover ./... && cat cover | grep -v "mock" | grep -v  "easyjson" | grep -v "proto" > cover.out && go tool cover -func=cover.out
//go tool cover -html=cover.out

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
