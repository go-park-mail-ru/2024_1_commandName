package chats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"ProjectMessenger/domain"
	chat "ProjectMessenger/internal/chats/repository/db"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestChatRepo_GetChatByChatID_Succes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "1", "test@mail.ru", "Test", "User", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

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
	chatRepo := chat.NewRawChatsStorage(db)

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
	chatRepo := chat.NewRawChatsStorage(db)

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
	chatRepo := chat.NewRawChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery(`INSERT INTO chat\.chat \(type_id, name, description, avatar_path, created_at,edited_at, creator_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
	fmt.Println("END OF TEST SUCCESSFULL")
}

func TestChatRepo_CreateChat_LengthErr(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewRawChatsStorage(db)

	ctx := context.Background()
	arrOfUserIDs := make([]uint, 0)
	arrOfUserIDs = append(arrOfUserIDs, uint(1))
	_, err = chatRepo.CreateChat(ctx, "chat1", "desc", arrOfUserIDs...)
	if err == nil {
		t.Error("err == nil, must be not nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_CreateChat_Group(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewRawChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery(`INSERT INTO chat\.chat \(type_id, name, description, avatar_path, created_at,edited_at, creator_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	fmt.Println("here")
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	ctx := context.Background()
	arrOfUserIDs := make([]uint, 0)
	arrOfUserIDs = append(arrOfUserIDs, uint(1))
	arrOfUserIDs = append(arrOfUserIDs, uint(2))
	arrOfUserIDs = append(arrOfUserIDs, uint(3))
	_, err = chatRepo.CreateChat(ctx, "chat1", "desc", arrOfUserIDs...)
	if err != nil {
		t.Error("err:", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_CreateChat_CustomErrorByFirst(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewRawChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery(`INSERT INTO chat\.chat \(type_id, name, description, avatar_path, created_at,edited_at, creator_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some err"))

	// Создание контекста и вызов функции
	ctx := context.Background()
	arrOfUserIDs := []uint{1, 2, 3} // Используем сразу инициализированный срез
	_, err = chatRepo.CreateChat(ctx, "chat1", "desc", arrOfUserIDs...)

	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected an error but got nil")
	} else {
		// Тут можно добавить дополнительные проверки, если нужно
	}

	// Проверка выполнения ожиданий mock объекта
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_CreateChat_CustomErrorSecond(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// Создание userRepo с mock базы данных
	chatRepo := chat.NewRawChatsStorage(db)

	// Утверждение ожидания запроса к базе данных
	mock.ExpectQuery(`INSERT INTO chat\.chat \(type_id, name, description, avatar_path, created_at,edited_at, creator_id\) VALUES (.+) RETURNING id`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	fmt.Println("here")

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	arrOfUserIDs := []uint{1, 2, 3} // Используем сразу инициализированный срез
	_, err = chatRepo.CreateChat(ctx, "chat1", "desc", arrOfUserIDs...)

	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected an error but got nil")
	} else {
		// Тут можно добавить дополнительные проверки, если нужно
	}

	// Проверка выполнения ожиданий mock объекта
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_DeleteChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectExec("DELETE FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM chat.message WHERE chat_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	ok, err := chatRepo.DeleteChat(ctx, 1)

	// Проверка наличия ошибки
	if err != nil {
		t.Error("Expected nil, but got nil")
	}
	if !ok {
		t.Error("Expected ok, got false")
	}

	// Проверка выполнения ожиданий mock объекта
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_DeleteChat_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectExec("DELETE FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	ok, err := chatRepo.DeleteChat(ctx, 1)

	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected an error but got nil")
	}
	if ok {
		t.Error("Expected !ok, got true")
	}

	// Проверка выполнения ожиданий mock объекта
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepo_DeleteChat_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectExec("DELETE FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM chat.message WHERE chat_id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	ok, err := chatRepo.DeleteChat(ctx, 1)

	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected nil, but got nil")
	}
	if err == nil {
		t.Error("Expected an error but got nil")
	}
	if ok {
		t.Error("Expected !ok, got true")
	}
}

func TestChatRepo_DeleteChat_Error3(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectExec("DELETE FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM chat.message WHERE chat_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	ok, err := chatRepo.DeleteChat(ctx, 1)

	// Проверка наличия ошибки
	if err == nil {
		t.Error("Expected nil, but got nil")
	}
	if err == nil {
		t.Error("Expected an error but got nil")
	}
	if ok {
		t.Error("Expected !ok, got true")
	}
}

func TestUserRepo_GetChatsForUser(t *testing.T) {
	// Создание mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "2", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1).
			AddRow(2, "2", "name2", "desc", "avatar_path", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited", "username"}).
			AddRow(1, 1, 1, "desc", fixedTime, false, "artem").
			AddRow(2, 2, 2, "desc", fixedTime, false, "alex"))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited", "username"}).
			AddRow(1, 1, 1, "desc", fixedTime, false, "artem").
			AddRow(2, 2, 2, "desc", fixedTime, false, "alex"))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

	ctx := context.Background()
	contacts := chatRepo.GetChatsForUser(ctx, 1)
	if len(contacts) == 0 {
		t.Error("lem must be not 0!")
	}

	fmt.Println(contacts)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetChatsForUser_CustomError1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	contacts := chatRepo.GetChatsForUser(ctx, 1)
	if len(contacts) != 0 {
		t.Error("lem must be 0!")
	}

	fmt.Println(contacts)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetChatsForUser_CustomError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)

	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "2", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1).
			AddRow(2, "2", "name2", "desc", "avatar_path", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited", "username"}).
			AddRow(1, 1, 1, "desc", fixedTime, false, "artem").
			AddRow(2, 2, 2, "desc", fixedTime, false, "alex"))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

	ctx := context.Background()
	contacts := chatRepo.GetChatsForUser(ctx, 1)
	if contacts == nil {
		t.Error("len must be 0!")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetChatsForUser_CustomError3(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)

	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "2", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited", "username"}).
			AddRow(1, 1, 1, "desc", fixedTime, false, "artem"))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	contacts := chatRepo.GetChatsForUser(ctx, 1)
	if contacts == nil {
		t.Error("len must be 0!")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateGroupChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)
	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectExec(`UPDATE chat\.chat SET name=\$1, description=\$2 WHERE id=\$3`).
		WithArgs("newChat", "newDesc", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedChat := domain.Chat{
		ID:          1,
		Name:        "newChat",
		Description: "newDesc",
	}

	ctx := context.Background()
	ok := chatRepo.UpdateGroupChat(ctx, updatedChat)
	if !ok {
		t.Error("err: ok is false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateGroupChat_CustomError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewChatsStorage(db)
	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectExec(`UPDATE chat\.chat SET name=\$1, description=\$2 WHERE id=\$3`).
		WithArgs("newChat", "newDesc", 1).
		WillReturnError(errors.New("some error"))

	updatedChat := domain.Chat{
		ID:          1,
		Name:        "newChat",
		Description: "newDesc",
	}

	ctx := context.Background()
	ok := chatRepo.UpdateGroupChat(ctx, updatedChat)
	if ok {
		t.Error("err: ok is true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}