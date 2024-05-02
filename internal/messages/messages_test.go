package messages

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"ProjectMessenger/domain"
	database "ProjectMessenger/internal/messages/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestSetMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewMessageStorage(db)
	mock.ExpectQuery(`INSERT INTO chat\.message \(user_id, chat_id, message, edited_at, created_at\) VALUES(.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	ctx := context.Background()
	var message domain.Message
	message.ID = 1
	message.UserID = 1
	messageSaved := userRepo.SetMessage(ctx, message)
	fmt.Println(messageSaved)
	if messageSaved.ID == 0 {
		t.Error("id is 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSetMessage_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewMessageStorage(db)
	mock.ExpectQuery(`INSERT INTO chat\.message \(user_id, chat_id, message, edited_at, created_at\) VALUES(.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	var message domain.Message
	message.ID = 1
	message.UserID = 1
	messageSaved := userRepo.SetMessage(ctx, message)
	fmt.Println(messageSaved)
	if messageSaved.ID != 0 {
		t.Error("id must be 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetChatMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewMessageStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery("^SELECT message.id, user_id, chat_id, message.message, created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = \\$1$").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited_at", "username"}).AddRow(1, 1, 2, "message", fixedTime, fixedTime, "artem"))

	ctx := context.Background()

	messages := userRepo.GetChatMessages(ctx, 1, 1)
	if len(messages) == 0 {
		t.Error("len is 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetChatMessages_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := database.NewMessageStorage(db)

	mock.ExpectQuery("^SELECT message.id, user_id, chat_id, message.message, created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = \\$1$").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()

	messages := userRepo.GetChatMessages(ctx, 1, 1)
	if len(messages) != 0 {
		t.Error("len is not 0, but it`s test with error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
