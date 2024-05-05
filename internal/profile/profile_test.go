package profile

import (
	"database/sql"
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

type updateUserStruct[T any] struct {
	User               T   `json:"user"`
	NumOfUpdatedFields int `json:"numOfUpdatedFields"`
}

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

func TestUpdateProfileInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{
  "user": {
    "Username": "TestUser",
    "Email": "test_user@example.com"
  },
  "numOfUpdatedFields": 2
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("TestUser", "test_user@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.UpdateProfileInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestUpdateProfileInfo_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("GET", "/getProfileInfo", strings.NewReader(`{
  "user": {
    "Username": "TestUser",
    "Email": "test_user@example.com"
  },
  "numOfUpdatedFields": 2
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.UpdateProfileInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}
func TestUpdateProfileInfo_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{
  "user": {
    "Username": "TestUser",
    "Email "test_user@example.com"
  },
  "numOfUpdatedFields": 2
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.UpdateProfileInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestUpdateProfileInfo_Error3(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{
  "user": {
    "Username": "TestUser",
    "Email": "test_user@example.com"
  },
  "numOfUpdatedFields": 0
}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.UpdateProfileInfo(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestChangePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"oldPassword": "Demouser123!",
   "newPassword":"newPass123!"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("TestUser", "test@mail.ru", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.ChangePassword(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestChangePassword_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("GET", "/getProfileInfo", strings.NewReader(`{"oldPassword": "Demouser123!",
   "newPassword":"newPass123!"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.ChangePassword(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestChangePassword_Error2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"oldPassword: "Demouser123!",
   "newPassword":"newPass123!"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.ChangePassword(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestChangePassword_Error3(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"oldPassword": "Demouser123!",
   "newPassword":""}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.ChangePassword(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestGetContacts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"oldPassword": "Demouser123!",
   "newPassword":"newPass123!"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery("^SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON ").
		WithArgs(uint(1), 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(1, "Artem", "artem@mail.ru", "Artem", "Chernikov", "Developer", time.Now(), ""))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.GetContacts(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestAddContact(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"username_of_user_to_add": "Friend"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = ?`).
		WithArgs("Friend").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(2, "Friend", "friend@mail.ru", "Friend", "friend", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "", "5t2HF7Tq"))

	mock.ExpectQuery("^SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON ").
		WithArgs(uint(1), 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(2, "artem", "artem@mail.ru", "Artem", "Chernikov", "Developer", time.Now(), ""))

	mock.ExpectQuery(`INSERT INTO chat\.contacts \(user1_id, user2_id, state_id\) VALUES (.+) RETURNING id`).
		WithArgs(uint(1), uint(2), 3).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.AddContact(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}

func TestAddContact_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	req := httptest.NewRequest("POST", "/getProfileInfo", strings.NewReader(`{"username_of_user_to_add": "Friend"}`))
	req.Header.Set("Cookie", "session_id=yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo")
	w := httptest.NewRecorder()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = ?`).
		WithArgs("Friend").
		WillReturnError(sql.ErrNoRows)

	auth := authDelivery.NewRawAuthHandler(db, "")
	profile := delivery.NewProfileHandler(auth)

	profile.AddContact(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code должен быть 200")
}
