package messages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"regexp"
	"testing"
	"time"

	"ProjectMessenger/domain"
	authDelivery "ProjectMessenger/internal/auth/delivery"
	chats "ProjectMessenger/internal/chats_service/proto"
	contactsProto "ProjectMessenger/internal/contacts_service/proto"
	database "ProjectMessenger/internal/messages/repository/db"
	"ProjectMessenger/internal/messages/usecase"
	session "ProjectMessenger/internal/sessions_service/proto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
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

	// "SELECT message.id, user_id, chat_id, message.message, message.created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = $1"
	// "^SELECT message.id, user_id, chat_id, message.message, created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = \\$1$"
	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = ?").
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

	mock.ExpectQuery("SELECT message.id, user_id, chat_id, message.message, message.created_at, edited_at, username FROM chat.message JOIN auth.person ON message.user_id = person.id WHERE chat_id = ?").
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

type MockWebsocket struct {
	mock.Mock
}

func (m *MockWebsocket) SendMessageToUser(userID uint, message []byte) error {
	args := m.Called(userID, message)
	return args.Error(0)
}

func TestSendMessageToUser(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
	}))
	defer s.Close()

	dialer := websocket.Dialer{}
	mockConn, _, _ := dialer.Dial("ws"+s.URL[4:], nil)
	defer mockConn.Close()
	userID := uint(1)
	message := []byte("Test message")
	mockWebsocket := new(MockWebsocket)
	mockWebsocket.On("SendMessageToUser", userID, message).Return(nil)

	ws := &database.Websocket{
		Connections: map[uint]*websocket.Conn{},
	}

	ctx := context.Background()
	ws.AddConnection(ctx, mockConn, 1)
	err := ws.SendMessageToUser(userID, message)

	assert.Equal(t, nil, err)
}

func TestSendMessageToUser_Error(t *testing.T) {

	expectedErr := errors.New("No connection found for user")
	userID := uint(1)
	message := []byte("Test message")
	mockWebsocket := new(MockWebsocket)
	mockWebsocket.On("SendMessageToUser", userID, message).Return(expectedErr)

	ws := &database.Websocket{
		Connections: map[uint]*websocket.Conn{},
	}

	err := ws.SendMessageToUser(userID, message)

	assert.Equal(t, expectedErr, err)
}

func TestCreateWSStorage(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	wsStorage := database.NewWsStorage(db)
	if wsStorage == nil {
		t.Errorf("wsStorage is nil")
	}
}

func TestDeleteConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
	}))
	defer s.Close()

	dialer := websocket.Dialer{}
	mockConn, _, _ := dialer.Dial("ws"+s.URL[4:], nil)
	defer mockConn.Close()

	ws := &database.Websocket{
		Connections: map[uint]*websocket.Conn{},
	}

	ctx := context.Background()
	ws.AddConnection(ctx, mockConn, 1)
	ws.DeleteConnection(1)
}

func TestDeleteMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	ctx := context.Background()
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), created_at FROM chat.message WHERE id = $1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"})).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"}).
			AddRow(1, 1, uint(1), "message", false, fixedTime, fixedTime))

	messRepo := database.NewMessageStorage(db)
	usecase.DeleteMessage(ctx, uint(1), uint(1), messRepo)
}

func TestEditMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	ctx := context.Background()
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), created_at FROM chat.message WHERE id = $1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"})).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"}).
			AddRow(1, 1, uint(1), "message", false, fixedTime, fixedTime))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE chat.message SET message = $1, edited = $2, edited_at = $3 WHERE id = $4")).
		WithArgs("new msg", true, fixedTime, uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	messRepo := database.NewMessageStorage(db)
	err = usecase.EditMessage(ctx, uint(1), uint(1), "new text", messRepo)
	if err != nil {
		fmt.Println(err)
	}
}

func TestEditMessage_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	ctx := context.Background()
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), created_at FROM chat.message WHERE id = $1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"})).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"}).
			AddRow(1, 1, uint(1), "message", false, fixedTime, fixedTime))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE chat.message SET message = $1, edited = $2, edited_at = $3 WHERE id = $4")).
		WithArgs("new msg", true, fixedTime, uint(1)).
		WillReturnError(errors.New("some err"))

	messRepo := database.NewMessageStorage(db)
	err = usecase.EditMessage(ctx, uint(1), uint(1), "new text", messRepo)
	if err != nil {
		fmt.Println(err)
	}
}

func TestEditMessage_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	ctx := context.Background()
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, chat_id, message.message, edited, COALESCE(edited_at, '2000-01-01 00:00:00'), created_at FROM chat.message WHERE id = $1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "edited_at", "created_at"})).
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE chat.message SET message = $1, edited = $2, edited_at = $3 WHERE id = $4")).
		WithArgs("new msg", true, fixedTime, uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	messRepo := database.NewMessageStorage(db)
	err = usecase.EditMessage(ctx, uint(1), uint(1), "new text", messRepo)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSendMessageToOtherUsers(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()
	message := domain.Message{
		Message: "Hello",
		ChatID:  uint(1),
	}
	wsStorage := database.NewWsStorage(db)

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	chatsManager := chats.NewChatServiceClient(grcpChats)

	usecase.SendMessageToOtherUsers(ctx, message, uint(1), wsStorage, chatsManager)
}

func TestGetFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()
	//wsStorage := database.NewWsStorage(db)

	mock.ExpectQuery("SELECT file_path FROM chat.file WHERE message_id =?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"file_path"}).AddRow("file_path"))

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	//chatsManager := chats.NewChatServiceClient(grcpChats)
	messRepo := database.NewMessageStorage(db)
	usecase.GetFile(ctx, messRepo, uint(1))
}

func TestGetFile_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()
	//wsStorage := database.NewWsStorage(db)

	mock.ExpectQuery("SELECT file_path FROM chat.file WHERE message_id =?").
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	//chatsManager := chats.NewChatServiceClient(grcpChats)
	messRepo := database.NewMessageStorage(db)
	usecase.GetFile(ctx, messRepo, uint(1))
}

func TestSetFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()

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
	//wsStorage := database.NewWsStorage(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO chat.file (user_id, message_id, file_path) VALUES($1, $2, $3)")).
		WithArgs(uint(1), uint(1), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"file_path"}).AddRow("file_path"))

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	//chatsManager := chats.NewChatServiceClient(grcpChats)
	messRepo := database.NewMessageStorage(db)

	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после теста

	content := []byte("This is a test file content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	file := tempFile

	// Создаем фейковый multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: tempFile.Name(),
		Header:   textproto.MIMEHeader{"Content-Type": []string{"text/plain"}},
		Size:     int64(len(content)),
	}

	usecase.SetFile(messRepo, ctx, file, uint(1), uint(1), authHandler.Users, fileHeader)
}

func TestSetFile_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()

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
	//wsStorage := database.NewWsStorage(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = $1")).
		WithArgs(1).
		WillReturnError(errors.New("some err"))

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	//chatsManager := chats.NewChatServiceClient(grcpChats)
	messRepo := database.NewMessageStorage(db)

	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после теста

	content := []byte("This is a test file content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	file := tempFile

	// Создаем фейковый multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: tempFile.Name(),
		Header:   textproto.MIMEHeader{"Content-Type": []string{"text/plain"}},
		Size:     int64(len(content)),
	}

	usecase.SetFile(messRepo, ctx, file, uint(1), uint(1), authHandler.Users, fileHeader)
}

func TestSetFile_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	ctx := context.Background()

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
	//wsStorage := database.NewWsStorage(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO chat.file (user_id, message_id, file_path) VALUES($1, $2, $3)")).
		WithArgs(uint(1), uint(1), sqlmock.AnyArg()).
		WillReturnError(errors.New("some err"))

	grcpChats, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	defer grcpChats.Close()
	//chatsManager := chats.NewChatServiceClient(grcpChats)
	messRepo := database.NewMessageStorage(db)

	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после теста

	content := []byte("This is a test file content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	file := tempFile

	// Создаем фейковый multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: tempFile.Name(),
		Header:   textproto.MIMEHeader{"Content-Type": []string{"text/plain"}},
		Size:     int64(len(content)),
	}

	usecase.SetFile(messRepo, ctx, file, uint(1), uint(1), authHandler.Users, fileHeader)
}
