package profile

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"ProjectMessenger/domain"
	authDelivery "ProjectMessenger/internal/auth/delivery"
	contactsProto "ProjectMessenger/internal/contacts_service/proto"
	"ProjectMessenger/internal/profile/usecase"
	session "ProjectMessenger/internal/sessions_service/proto"
	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc"
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

	mock.ExpectQuery("SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "TestUser", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("MerBzjWU8qUKGaStn4iFNoDOk7SUyW9w").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "MerBzjWU8qUKGaStn4iFNoDOk7SUyW9w"))

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
	ctx := context.Background()
	usecase.GetProfileInfo(ctx, uint(1), auth.Users)

}

func TestUpdateProfileInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("new", "test_user@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

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
	//profile := delivery.NewProfileHandler(auth, contactsManager)
	person := domain.Person{Username: "new", Email: "test_user@example.com"}

	ctx := context.Background()

	usecase.UpdateProfileInfo(ctx, person, 2, uint(1), auth.Users)
}

func TestUpdateProfileInfo_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "artem", "test@mail.ru", "Test", "User", "Developer", "5baae85b9413d75de29d9e54b0550eae8ea8eaabb80b0cea8974bb5ee844b82fd9c45d188938bbc57716a495a3766b1728bdffb04f256a67ad545b62d9e69ac7", time.Now(), time.Now(), "", "gxYdyp8Z"))

	mock.ExpectExec(`UPDATE auth\.person SET username = \$1, email = \$2, name = \$3, surname = \$4, about = \$5, password_hash = \$6, created_at = \$7, lastseen_at = \$8, avatar_path = \$9, password_salt = \$10 WHERE id = \$11`).
		WithArgs("new", "test_user@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

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
	//profile := delivery.NewProfileHandler(auth, contactsManager)
	person := domain.Person{Username: "new", Email: "test_user@example.com", Name: "Artem", Surname: "Chernikov", About: "about"}

	ctx := context.Background()

	usecase.UpdateProfileInfo(ctx, person, 4, uint(1), auth.Users)

}

func TestChangePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

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
	//profile := delivery.NewProfileHandler(auth, contactsManager)
	ctx := context.Background()
	usecase.ChangePassword(ctx, "Demouser123!", "newPass123!", uint(1), auth.Users)
}

func TestGetContacts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT userid, sessionid FROM auth.session WHERE sessionid = \\$1$").
		WithArgs("yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo").
		WillReturnRows(sqlmock.NewRows([]string{"userid", "sessionid"}).AddRow(1, "yOQGFWqFFEkWwigIT29cP8N9HMtkGwDo"))

	mock.ExpectQuery("^SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, ap.lastseen_at, ap.avatar_path FROM chat.contacts cc JOIN auth.person ap ON ").
		WithArgs(uint(1), 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "lastseen_at", "avatar_path"}).
			AddRow(1, "Artem", "artem@mail.ru", "Artem", "Chernikov", "Developer", time.Now(), ""))

	grcpSessions, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpSessions.Close()
	//sessManager := session.NewAuthCheckerClient(grcpSessions)

	grcpContacts, err := grpc.Dial(
		"127.0.0.1:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpContacts.Close()
	contactsManager := contactsProto.NewContactsClient(grcpContacts)

	//auth := authDelivery.NewRawAuthHandler(db, sessManager, "", contactsManager)
	ctx := context.Background()
	usecase.GetContacts(ctx, uint(1), contactsManager)

}

func TestAddContact(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

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
	ctx := context.Background()
	usecase.AddContactByUsername(ctx, uint(1), "Friend", auth.Users, contactsManager)
}

func TestAddContact_Error1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = ?`).
		WithArgs("Friend").
		WillReturnError(errors.New("some err"))

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
	ctx := context.Background()
	usecase.AddContactByUsername(ctx, uint(1), "Friend", auth.Users, contactsManager)
}

func TestAddToAllContacts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id FROM auth.person`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2).AddRow(1))

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
	ctx := context.Background()
	usecase.AddToAllContacts(ctx, uint(1), auth.Users, contactsManager)
}

func TestChangeAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "Friend", "friend@mail.ru", "Friend", "friend", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "avatars/image.jpg", "5t2HF7Tq"))

	mock.ExpectQuery(`SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = ?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name", "surname", "about", "password_hash", "created_at", "lastseen_at", "avatar_path", "password_salt"}).
			AddRow(1, "Friend", "friend@mail.ru", "Friend", "friend", "Developer", "59c85ba56081d0b478d5acaa0a53e8c3c8f3bfd62a3fbafe7a1b09df37ede22e8745eda7646f67b565fcc533f50a7e9802e6972c29f6816d6a7bdb2c01eda7f2", time.Now(), time.Now(), "avatars/image.jpg", "5t2HF7Tq"))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE auth.person SET username = $1, email = $2, name = $3, surname = $4, about = $5, password_hash = $6, created_at = $7, lastseen_at = $8, avatar_path = $9, password_salt = $10 WHERE id = $11`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

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
	ctx := context.Background()
	multipartFile, fileHeader, err := createFakeMultipartFile()
	usecase.ChangeAvatar(ctx, multipartFile, fileHeader, uint(1), auth.Users)
}

func createFakeMultipartFile() (*os.File, *multipart.FileHeader, error) {
	// Создаем буфер для хранения multipart данных
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Открываем изображение
	file, err := os.Open("avatars/image1.jpg")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Создаем форму файла
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, nil, err
	}

	// Копируем содержимое файла в форму
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	// Закрываем writer, чтобы завершить формирование multipart данных
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	// Создаем multipart.File и multipart.FileHeader для тестирования
	multipartFile, multipartFileHeader, err := createMultipartFileFromBuffer(body, writer)
	if err != nil {
		return nil, nil, err
	}

	return multipartFile, multipartFileHeader, nil
}

// createMultipartFileFromBuffer создает multipart.File и multipart.FileHeader из буфера
func createMultipartFileFromBuffer(buffer *bytes.Buffer, writer *multipart.Writer) (*os.File, *multipart.FileHeader, error) {
	// Создаем multipart reader
	reader := multipart.NewReader(buffer, writer.Boundary())

	// Читаем форму
	part, err := reader.NextPart()
	if err != nil {
		return nil, nil, err
	}

	// Создаем временный файл
	tempFile, err := os.CreateTemp("", "multipart-*")
	if err != nil {
		return nil, nil, err
	}

	// Копируем содержимое части в временный файл
	_, err = io.Copy(tempFile, part)
	if err != nil {
		return nil, nil, err
	}

	// Возвращаемся к началу временного файла
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, nil, err
	}

	// Создаем FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: part.FileName(),
		Header:   part.Header,
		Size:     int64(buffer.Len()),
	}

	return tempFile, fileHeader, nil
}
