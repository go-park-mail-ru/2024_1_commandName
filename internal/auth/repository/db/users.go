package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
)

type Users struct {
	db           *sql.DB
	countOfUsers uint
	pathToAvatar string
}

func NewUserStorage(db *sql.DB, pathToAvatar string) *Users {
	slog.Info("created user storage")
	return &Users{
		db:           CreateFakeUsers(6, db),
		countOfUsers: 6,
		pathToAvatar: pathToAvatar,
	}
}

func NewRawUserStorage(db *sql.DB, pathToAvatar string) *Users {
	slog.Info("created user storage")
	return &Users{
		db:           db,
		countOfUsers: 6,
		pathToAvatar: pathToAvatar,
	}
}

func (u *Users) GetAllUserIDs(ctx context.Context) (userIDs []uint) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	userIDs = make([]uint, 0)
	rows, err := u.db.QueryContext(ctx, "SELECT id FROM auth.person")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetAllUserIDs: no IDs")
			return userIDs
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetAllUserIDs, users.go",
		}
		logger.Error(customErr.Error())
		return nil
	}

	for rows.Next() {
		var id uint
		err = rows.Scan(&id)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method GetAllUserIDs, users.go",
			}
			logger.Error(customErr.Error())
			return nil
		}
		userIDs = append(userIDs, id)
	}
	logger.Debug("GetAllUserIDs: found contacts", "numOfIDs", len(userIDs))
	return userIDs
}

func (u *Users) GetByUsername(ctx context.Context, username string) (user domain.Person, found bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("GetByUsername", "username", username)
	err := u.db.QueryRowContext(ctx, "SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.AvatarPath, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetByUsername didn't found user", "username", username)
			return user, false
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetByUsername, users.go",
		}
		logger.Error(customErr.Error())
		return user, false
	}
	logger.Debug("GetByUsername found user", "username", username)
	return user, true
}

func (u *Users) CreateUser(ctx context.Context, user domain.Person) (userID uint, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	if user.AvatarPath == "" {
		user.AvatarPath = "avatars/avatar.jpg"
	}

	err = u.db.QueryRowContext(ctx, "INSERT INTO auth.person (username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.AvatarPath, user.PasswordSalt).
		Scan(&userID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateUser, users.go",
		}
		logger.Error(customErr.Error())
		return 0, err
	}

	logger.Info("created user", "userID", user)
	u.countOfUsers++
	return userID, nil
}

func (u *Users) GetByUserID(ctx context.Context, userID uint) (user domain.Person, found bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	err := u.db.QueryRowContext(ctx, "SELECT id, username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt FROM auth.person WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.About, &user.Password, &user.CreateDate, &user.LastSeenDate, &user.AvatarPath, &user.PasswordSalt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetByUserID didn't found user", "userID", userID)
			return user, false
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetByUserID, users.go",
		}
		logger.Error(customErr.Error())
		return user, false
	}
	logger.Debug("GetByUserID found user", "userID", userID)
	return user, true
}

func (u *Users) UpdateUser(ctx context.Context, userUpdated domain.Person) (ok bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	oldUser, found := u.GetByUserID(ctx, userUpdated.ID)
	if !found {
		logger.Debug("UpdateUser didn't found user(via GetByUserID)", "userID", userUpdated.ID)
		return false
	}

	_, err := u.db.ExecContext(ctx, "UPDATE auth.person SET username = $1, email = $2, name = $3, surname = $4, about = $5, password_hash = $6, created_at = $7, lastseen_at = $8, avatar_path = $9, password_salt = $10 WHERE id = $11",
		userUpdated.Username, userUpdated.Email, userUpdated.Name, userUpdated.Surname, userUpdated.About, userUpdated.Password, userUpdated.CreateDate, userUpdated.LastSeenDate, userUpdated.AvatarPath, userUpdated.PasswordSalt, oldUser.ID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateUser, users.go",
		}
		logger.Error(customErr.Error())
		return false
	}
	logger.Debug("UpdateUser success", "userID", userUpdated.ID)
	return true
}

func (u *Users) StoreAvatar(ctx context.Context, multipartFile multipart.File, fileHandler *multipart.FileHeader) (name string, err error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	originalName := fileHandler.Filename
	fileNameSlice := strings.Split(originalName, ".")
	if len(fileNameSlice) < 2 {
		logger.Info("StoreAvatar filename without extension")
		return "", fmt.Errorf("Файл не имеет расширения")
	}
	extension := fileNameSlice[len(fileNameSlice)-1]
	if extension != "png" && extension != "jpg" && extension != "jpeg" && extension != "webp" && extension != "pjpeg" {
		logger.Info("StoreAvatar file isn't an image")
		return "", fmt.Errorf("Файл не является изображением")
	}

	filename := misc.RandStringRunes(16)
	filePath := u.pathToAvatar + "avatars/" + filename + "." + extension

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("StoreAvatar failed to open a file", "path", filePath)
		return "", fmt.Errorf("internal error")
	}
	defer f.Close()

	_, err = io.Copy(f, multipartFile)
	if err != nil {
		logger.Error("StoreAvatar failed to copy file", "path", filePath)
		return "", fmt.Errorf("internal error")
	}
	logger.Debug("StoreAvatar success", "path", filePath)
	return filename + "." + extension, nil
}

func (u *Users) GetAvatarStoragePath() string {
	return u.pathToAvatar
}

func (u *Users) SetFirebaseToken(ctx context.Context, userID uint, token string) (ok bool) {
	query := "INSERT INTO auth.notification (user_id, token) VALUES ($1, $2)"
	u.db.QueryRowContext(ctx, query, userID, token)
	return true
}

func (u *Users) GetTokensForUser(ctx context.Context, userID uint) ([]string, error) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	query := "SELECT token FROM auth.notification WHERE user_id = $1"

	rows, err := u.db.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetTokensForUser: no tokens")
			return nil, nil
		} else {
			logger.Error(err.Error())
			return nil, err
		}
	}
	tokens := make([]string, 0)
	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func CreateFakeUsers(countOfUsers int, db *sql.DB) *sql.DB {
	counter := 0
	_ = db.QueryRow("SELECT count(id) FROM auth.person").Scan(&counter)
	if counter == 0 {
		_, err := db.Exec("ALTER SEQUENCE auth.person_id_seq RESTART WITH 1")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			slog.Error(customErr.Error())
		}
		_, err = db.Exec("ALTER SEQUENCE auth.session_id_seq RESTART WITH 1")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			slog.Error(customErr.Error())
		}
		_, err = db.Exec("DELETE FROM auth.person")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			slog.Error(customErr.Error())
		}

		_, err = db.Exec("DELETE FROM auth.session")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			slog.Error(customErr.Error())
		}

		for i := 0; i < countOfUsers; i++ {
			query := `INSERT INTO auth.person (username, email, name, surname, about, password_hash, created_at, lastseen_at, avatar_path, password_salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			user := getFakeUser(i + 1)
			_, err := db.Exec(query, user.Username, user.Email, user.Name, user.Surname, user.About, user.Password, user.CreateDate, user.LastSeenDate, user.AvatarPath, user.PasswordSalt)
			if err != nil {
				customErr := &domain.CustomError{
					Type:    "database",
					Message: err.Error(),
					Segment: "method CreateFakeUsers, users.go",
				}
				slog.Error(customErr.Error())
			}
		}
	}

	counter = 0
	_ = db.QueryRow("SELECT count(id) FROM chat.contacts").Scan(&counter)
	if counter == 0 {
		_, err := db.Exec("ALTER SEQUENCE chat.contacts_id_seq RESTART WITH 1")
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers, users.go",
			}
			slog.Error(customErr.Error())
		}

		query := `INSERT INTO chat.contact_state (id, name) VALUES ($1, $2)`
		_, err = db.Exec(query, 1, "request_from")
		_, err = db.Exec(query, 2, "request_to")
		_, err = db.Exec(query, 3, "friend")

		query = `INSERT INTO chat.contacts (user1_id, user2_id, state_id) VALUES ($1, $2, $3)`
		_, err = db.Exec(query, 1, 2, 3) // Naumov to Chernikov -- friends
		_, err = db.Exec(query, 2, 3, 3) // Chernikov to Zhuk -- friends
		_, err = db.Exec(query, 6, 5, 3) // mentor to TestUser -- no answer
		_, err = db.Exec(query, 4, 6, 3) // Volohov to TestUser -- friends
		_, err = db.Exec(query, 2, 6, 3) // Chernikov to TestUser -- friends
		_, err = db.Exec(query, 6, 1, 3) // Naumov to TestUser -- no answer
		_, err = db.Exec(query, 6, 3, 3) // TestUser to Zhuk -- no answer
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method CreateFakeUsers -> Create fake contacts, users.go",
			}
			slog.Error(customErr.Error())
		}
	}
	slog.Info("created fake users")
	return db
}

func getFakeUser(number int) domain.Person {
	testUserHash, testUserSalt := misc.GenerateHashAndSalt("Demouser123!")
	users := map[int]domain.Person{
		1: {ID: 1, Username: "IvanNaumov", Email: "ivan@mail.ru", Name: "Ivan", Surname: "Naumov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
		2: {ID: 2, Username: "ArtemkaChernikov", Email: "artem@mail.ru", Name: "Artem", Surname: "Chernikov",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
		3: {ID: 3, Username: "ArtemZhuk", Email: "artemZhuk@mail.ru", Name: "Artem", Surname: "Zhuk",
			About: "Backend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
		4: {ID: 4, Username: "AlexanderVolohov", Email: "Volohov@mail.ru", Name: "Alexander", Surname: "Volohov",
			About: "Frontend Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
		5: {ID: 5, Username: "mentor", Email: "mentor@mail.ru", Name: "Mentor", Surname: "Mentor",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
		6: {ID: 6, Username: "TestUser", Email: "test@mail.ru", Name: "Test", Surname: "User",
			About: "Developer", CreateDate: time.Now(), LastSeenDate: time.Now(), AvatarPath: "",
			PasswordSalt: testUserSalt, Password: testUserHash},
	}
	return users[number]
}
