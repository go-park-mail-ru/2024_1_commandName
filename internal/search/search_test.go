package search

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"ProjectMessenger/internal/chats/usecase"
	database "ProjectMessenger/internal/search/repository/db"
	tl "ProjectMessenger/internal/translate/usecase"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"

	"github.com/gorilla/websocket"
)

type Search struct {
	db          *sql.DB
	Connections map[uint]*websocket.Conn
	mu          sync.RWMutex
	Chats       usecase.ChatStore
	WebSocket   *MockWebsocket
	Translate   tl.TranslateStore
}

func TestSearchChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	searchRepo := database.NewSearchStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT c.id, c.type_id, c.name, c.description, c.avatar_path, c.created_at, c.edited_at, c.creator_id\n    FROM chat.chat c\n    JOIN chat.chat_user cu ON").
		WithArgs("name", "name", "name", "name", 1, "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type_id", "name", "description", "avatar_path", "created_at", "edited_at", "creator_id"}).AddRow(1, "1", "name", "desc", "", fixedTime, fixedTime, 1))

	ctx := context.Background()
	foundChat := searchRepo.SearchChats(ctx, "name", 1, "1")
	if len(foundChat.Chats) == 0 {
		t.Error("len is 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	searchRepo := database.NewSearchStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT m.id, m.user_id, m.chat_id, m.message, m.edited, m.created_at FROM chat.message m WHERE").
		WithArgs("new", "new", "new", "new", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "chat_id", "message", "edited", "created_at"}).AddRow(1, 1, 1, "new message", false, fixedTime))

	ctx := context.Background()
	foundMessage := searchRepo.SearchMessages(ctx, "new", 1)
	if len(foundMessage.Messages) == 0 {
		t.Error("len is 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchContacts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	searchRepo := database.NewSearchStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON").
		WithArgs("new", "new", "new", "new", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", fixedTime, ""))

	ctx := context.Background()
	foundMessage := searchRepo.SearchContacts(ctx, "new", 1)
	if len(foundMessage.Contacts) == 0 {
		t.Error("len is 0")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSendMatchedChatsSearchResponse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	searchRepo := database.NewSearchStorage(db)
	fixedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON").
		WithArgs("new", "new", "new", "new", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", fixedTime, ""))

	ctx := context.Background()
	foundMessage := searchRepo.SearchContacts(ctx, "new", 1)
	if len(foundMessage.Contacts) == 0 {
		t.Error("len is 0")
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

func (m *MockWebsocket) AddConnection(ctx context.Context, connection *websocket.Conn, userID uint) error {
	args := m.Called(ctx, connection, userID)
	return args.Error(0)
}

func TestAddConnection(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
	}))
	defer s.Close()

	dialer := websocket.Dialer{}
	mockConn, _, _ := dialer.Dial("ws"+s.URL[4:], nil)
	defer mockConn.Close()

	ctx := context.Background()
	searchStore := database.NewSearchStorage(db)
	searchStore.AddConnection(ctx, mockConn, uint(1))
	searchStore.DeleteConnection(1)
}

/*
func TestSendMatchedChatsSearchResponse1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebSocket := NewMockWebSocket(ctrl)
	mockWebSocket.EXPECT().SendMessageToUser(uint(1), gomock.Any()).Return(nil).Times(1)

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("не удалось создать mock: %s", err)
	}
	defer db.Close()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
	}))

	searchRepo := database.NewSearchStorage(db)
	searchRepo.WebSocket = mockWebSocket

	response := domain.ChatSearchResponse{}
	userID := uint(1)

	dialer := websocket.Dialer{}
	mockConn, _, _ := dialer.Dial("ws"+s.URL[4:], nil)
	defer mockConn.Close()

	ctx := context.Background()
	searchRepo.AddConnection(ctx, mockConn, userID)
	searchRepo.SendMatchedChatsSearchResponse(response, userID)

}*/
