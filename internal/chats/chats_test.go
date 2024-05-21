package chats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"ProjectMessenger/domain"
	authDelivery "ProjectMessenger/internal/auth/delivery"
	user "ProjectMessenger/internal/auth/repository/db"
	chatUsecase "ProjectMessenger/internal/chats/usecase"
	chats "ProjectMessenger/internal/chats_service/proto"
	chat "ProjectMessenger/internal/chats_service/repository"
	contactsProto "ProjectMessenger/internal/contacts_service/proto"
	session "ProjectMessenger/internal/sessions_service/proto"
	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc"
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

	chatRepo := chat.NewChatsStorage(db)

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "2", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1).
			AddRow(2, "2", "name2", "desc", "avatar_path", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"lastseen_message_id"}).
			AddRow(1))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"lastseen_message_id"}).
			AddRow(0))

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

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	contacts := chatRepo.GetChatsForUser(ctx, 1)
	if len(contacts) != 0 {
		t.Error("len must be 0!")
	}

	fmt.Println(contacts)
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

	chatRepo := chat.NewRawChatsStorage(db)

	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Утверждение ожидания запроса к базе данных и возвращение результата
	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "2", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1))

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

	chatRepo := chat.NewRawChatsStorage(db)
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
	chatStr := "chat"
	descStr := "desc"

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

	chatUsecase.UpdateGroupChat(ctx, uint(1), uint(1), &chatStr, &descStr, chatsManager)

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

	chatRepo := chat.NewRawChatsStorage(db)
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

func TestUserRepo_CheckPrivateChatExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT cu1.chat_id, cu1.user_id, cu2.user_id FROM chat.chat_user cu1 INNER JOIN chat.chat_user cu2 ON cu1.chat_id = cu2.chat_id WHERE cu1.user_id = \\$1 AND cu2.user_id = \\$2 AND cu1.user_id <> cu2.user_id$").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id1", "user_id2"}).AddRow(1, 1, 2))

	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "1", "test@mail.ru", "Test", "User", fixedTime, fixedTime, 1))

	ctx := context.Background()
	exists, _, err := chatRepo.CheckPrivateChatExists(ctx, 1, 2)
	if !exists {
		t.Error("Chat does not exists, but must to be")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_CheckPrivateChatExists_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT cu1.chat_id, cu1.user_id, cu2.user_id FROM chat.chat_user cu1 INNER JOIN chat.chat_user cu2 ON cu1.chat_id = cu2.chat_id WHERE cu1.user_id = \\$1 AND cu2.user_id = \\$2 AND cu1.user_id <> cu2.user_id$").
		WithArgs(1, 2).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	exists, _, err := chatRepo.CheckPrivateChatExists(ctx, 1, 2)
	if exists {
		t.Error("Chat exists, but must not to be")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_CheckPrivateChatExists_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT cu1.chat_id, cu1.user_id, cu2.user_id FROM chat.chat_user cu1 INNER JOIN chat.chat_user cu2 ON cu1.chat_id = cu2.chat_id WHERE cu1.user_id = \\$1 AND cu2.user_id = \\$2 AND cu1.user_id <> cu2.user_id$").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id1", "user_id2"}).AddRow(1, 1, 2))

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	exists, _, err := chatRepo.CheckPrivateChatExists(ctx, 1, 2)
	if exists {
		t.Error("Chat exists, but must not to be")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetLastSeenMessageId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"lastseen_message_id"}).
			AddRow(1))

	ctx := context.Background()
	id := chatRepo.GetLastSeenMessageId(ctx, 2, 1)
	if id != 1 {
		t.Error("Chat does not exists, but must to be")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetLastSeenMessageId_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 2).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	id := chatRepo.GetLastSeenMessageId(ctx, 2, 1)
	if id != 0 {
		t.Error("Must return err")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetFirstChatMessageID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT id FROM chat.message WHERE chat_id = \\$1 ORDER BY created_at LIMIT 1$").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	ctx := context.Background()
	id := chatRepo.GetFirstChatMessageID(ctx, 1)
	if id != 1 {
		t.Error("id is 0, not 1")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetFirstChatMessageID_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)

	mock.ExpectQuery("^SELECT id FROM chat.message WHERE chat_id = \\$1 ORDER BY created_at LIMIT 1$").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	ctx := context.Background()
	id := chatRepo.GetFirstChatMessageID(ctx, 1)
	if id != 0 {
		t.Error("id is not 0")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetMessagesByChatID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	chatRepo := chat.NewRawChatsStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("^SELECT message.id, user_id, chat_id, message.message, message.created_at, message.edited, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = \\$1$").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "created_at", "edited", "username"}).
			AddRow(1, 1, 1, "desc", fixedTime, false, "artem").
			AddRow(2, 2, 2, "desc", fixedTime, false, "alex"))

	ctx := context.Background()
	messages := chatRepo.GetMessagesByChatID(ctx, 1)
	if len(messages) == 0 {
		t.Error("len is 0!")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatUsecase_GetChatByChatID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	grcpSessions, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpSessions.Close()
	sessManager := session.NewAuthCheckerClient(grcpSessions)

	grcpContacts, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpContacts.Close()
	contactsManager := contactsProto.NewContactsClient(grcpContacts)

	authHandler := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)
	//chatsHandler := chatsDelivery.NewChatsHandler(authHandler, chatsManager)

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
	chat1, err := chatUsecase.GetChatByChatID(ctx, uint(1), uint(1), authHandler.Users, chatsManager)
	if err != nil {
	}
	fmt.Println(chat1)
}

func TestChatUsecase_GetChatsForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	userRepo := user.NewUserStorage(db, "")
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	grcpSessions, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpSessions.Close()

	grcpContacts, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpContacts.Close()

	//authHandler := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

	mock.ExpectQuery("SELECT id, type_id, name, description, avatar_path, created_at, edited_at,creator_id FROM chat.chat_user cu JOIN chat.chat c ON").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).
			AddRow(1, "1", "name1", "desc", "avatar_path", fixedTime, fixedTime, 1).
			AddRow(2, "2", "name2", "desc", "avatar_path", fixedTime, fixedTime, 1))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"lastseen_message_id"}).
			AddRow(1))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(2, 2))

	mock.ExpectQuery("^SELECT lastseen_message_id FROM chat.chat_user WHERE user_id = \\$1 and chat_id = \\$2$").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"lastseen_message_id"}).
			AddRow(0))

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

	ctx := context.Background()
	chats := chatUsecase.GetChatsForUser(ctx, uint(1), chatsManager, userRepo)
	if len(chats) == 0 {
		t.Error("lem must be not 0!")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		//t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatUsecase_CreatePrivateChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)
	userRepo := user.NewUserStorage(db, "")

	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?").
		WithArgs(3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(3, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery("SELECT cu1.chat_id, cu1.user_id, cu2.user_id FROM chat.chat_user cu1 INNER JOIN chat.chat_user cu2 ON cu1.chat_id = cu2.chat_id WHERE").
		WithArgs(1, 3).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id1", "user_id2"}).
			AddRow(1, 1, 2))

	ctx := context.Background()
	chatID, isNew, err := chatUsecase.CreatePrivateChat(ctx, uint(1), uint(3), chatsManager, userRepo)
	if err != nil {
	}
	fmt.Println(chatID, isNew)
}

func TestChatUsecase_DeleteChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

	mock.ExpectQuery("SELECT chat_id, user_id FROM chat.chat_user WHERE chat_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"chat_id", "user_id"}).
			AddRow(1, 1).
			AddRow(1, 2))

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
	wasDeleted, err := chatUsecase.DeleteChat(ctx, uint(1), uint(1), chatsManager)
	if err != nil {
	}
	fmt.Println(wasDeleted)
}

func TestChatUsecase_CreateGroupChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)
	//userRepo := user.NewUserStorage(db, "")

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

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	ctx := context.Background()
	userIDs := []uint{1, 2, 3}
	chatID, err := chatUsecase.CreateGroupChat(ctx, uint(1), userIDs, "new", "desc", chatsManager)
	if err != nil {
	}
	fmt.Println(chatID)
}

func TestChatUsecase_UpdateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

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

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	mock.ExpectExec(`INSERT INTO chat\.chat_user \(chat_id, user_id, lastseen_message_id\) VALUES\(.+\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(result)

	ctx := context.Background()
	userIDs := []uint{1, 2, 3}
	chatID, err := chatUsecase.CreateGroupChat(ctx, uint(1), userIDs, "new", "desc", chatsManager)
	if err != nil {
	}
	fmt.Println(chatID)
}
