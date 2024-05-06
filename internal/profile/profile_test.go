package profile

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authDelivery "ProjectMessenger/internal/auth/delivery"
	delivery "ProjectMessenger/internal/profile/delivery"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"texts: ["привет"], "folderId": "b1gq4i9e5unl47m0kj5f", "targetLanguageCode": "en"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)
	profile.GetProfileInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}
