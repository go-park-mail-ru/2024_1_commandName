package translate

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authDelivery "ProjectMessenger/internal/auth/delivery"
	chatsDelivery "ProjectMessenger/internal/chats/delivery"
	contactsProto "ProjectMessenger/internal/contacts_service/proto"
	session "ProjectMessenger/internal/sessions_service/proto"
	delivery "ProjectMessenger/internal/translate/delivery"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestTranslate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/translate", strings.NewReader(`{"texts": ["привет"], "folderId": "b1gq4i9e5unl47m0kj5f", "targetLanguageCode": "en"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, 1))

	mock.ExpectQuery("^SELECT language FROM auth.person WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"language"}).AddRow("ru"))

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

	auth := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)
	chats := chatsDelivery.NewRawChatsHandler(auth, db)
	ts := delivery.NewTranslateHandler(db, chats)
	ts.TranslateMessage(w, req)
}

func TestTranslate_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/translate", strings.NewReader(`{"texts": ["привет"], "folderId": "b1gq4i9e5unl47m0kj5f", "targetLanguageCode": "en"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGw")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGw").
		WillReturnError(sql.ErrNoRows)

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

	auth := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)

	chats := chatsDelivery.NewRawChatsHandler(auth, db)
	ts := delivery.NewTranslateHandler(db, chats)
	ts.TranslateMessage(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestTranslate_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/translate", strings.NewReader(`{"texts: ["привет"], "folderId": "b1gq4i9e5unl47m0kj5f", "targetLanguageCode": "en"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGw")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGw").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, 1))

	mock.ExpectQuery("^SELECT language FROM auth.person WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"language"}).AddRow("ru"))

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

	auth := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)

	chats := chatsDelivery.NewRawChatsHandler(auth, db)
	ts := delivery.NewTranslateHandler(db, chats)
	ts.TranslateMessage(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}
