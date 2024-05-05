package translate

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authDelivery "ProjectMessenger/internal/auth/delivery"
	chatsDelivery "ProjectMessenger/internal/chats/delivery"
	delivery "ProjectMessenger/internal/translate/delivery"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

	auth := authDelivery.NewRawAuthHandler(db, "")
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

	auth := authDelivery.NewRawAuthHandler(db, "")
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

	auth := authDelivery.NewRawAuthHandler(db, "")
	chats := chatsDelivery.NewRawChatsHandler(auth, db)
	ts := delivery.NewTranslateHandler(db, chats)
	ts.TranslateMessage(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}
