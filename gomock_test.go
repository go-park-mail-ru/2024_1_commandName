package _024_1_commandName

import (
	"database/sql"
	"testing"
	"time"

	"ProjectMessenger/internal/auth/delivery"
	"github.com/DATA-DOG/go-sqlmock"
)

/*
func TestGetByUserID(t *testing.T) {
	// Создание нового контроллера gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока UserStore
	mockUserStore := mock.NewMockUserStore(ctrl)

	// Устанавливаем ожидание вызова метода GetByUserID и его возвращаемое значение
	mockUserStore.EXPECT().GetByUserID(gomock.Any(), uint(1)).Return(domain.Person{}, true)

	// Создание объекта, использующего мок UserStore
	service := mypackage.NewService(mockUserStore)

	// Вызов метода, который использует GetByUserID
	user, found := service.GetUserByID(context.Background(), uint(1))

	// Проверка результата
	assert.True(t, found)
	assert.Equal(t, domain.Person{}, user)
}*/

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
	/*

		mock.
			ExpectQuery("SELECT uid FROM users WHERE phone =").
			WithArgs("81111111111").
			WillReturnError(sql.ErrNoRows)

		mock.
			ExpectExec("UPDATE users SET").
			WithArgs("81111111111", "kate@mail.ru", "Kate", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userId := models.User{
			Uid:  1,
			Name: "Daria",
		}

		user := models.UserData{
			Phone:  "81111111111",
			Email:  "kate@mail.ru",
			Name:   "Kate",
			Avatar: "http://127.0.0.1:5000/default/avatar/stas.jpg",
		}
		user.ID = 1

		c := context.Background()
		ctx := context.WithValue(c, "User", userId)

		err = userRepo.UpdateData(ctx, user)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}*/
}
