package repository

import (
	"ProjectMessenger/domain"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type Contacts struct {
	db *sql.DB
}

func NewContactsStorage(db *sql.DB) *Contacts {
	return &Contacts{
		db: db,
	}
}

func (u *Contacts) GetContacts(ctx context.Context, userID uint) []domain.Person {
	logger := slog.With("requestID", ctx.Value("traceID"))
	contactArr := make([]domain.Person, 0)
	rows, err := u.db.QueryContext(ctx,
		`
    SELECT ap.id, ap.username, ap.email, ap.name, ap.surname, ap.about, 
             ap.lastseen_at, ap.avatar_path
    FROM chat.contacts cc
    JOIN auth.person ap ON 
      (cc.user2_id = ap.id AND cc.user1_id = $1)  -- user is user2_id
    OR (cc.user1_id = ap.id AND cc.user2_id = $1)  -- user is user1_id
    WHERE cc.state_id = $2;
  `,
		userID, 3)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("GetContacts no contacts", "userID", userID)
			return contactArr
		}

		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method getContact, users.go",
		}
		logger.Error(customErr.Error())
		return contactArr
	}

	for rows.Next() {
		var userContact domain.Person
		err = rows.Scan(&userContact.ID, &userContact.Username, &userContact.Email,
			&userContact.Name, &userContact.Surname, &userContact.About,
			&userContact.LastSeenDate, &userContact.AvatarPath)
		if err != nil {
			customErr := &domain.CustomError{
				Type:    "database",
				Message: err.Error(),
				Segment: "method getContact, users.go",
			}
			logger.Error(customErr.Error())
			empty := make([]domain.Person, 0)
			return empty
		}
		contactArr = append(contactArr, userContact)
	}
	logger.Debug("GetContacts found contacts", "userID", userID, "numOfContacts", len(contactArr))
	return contactArr
}

func (u *Contacts) AddContact(ctx context.Context, userID1, userID2 uint) (ok bool) {
	logger := slog.With("requestID", ctx.Value("traceID"))
	entryID := 0
	err := u.db.QueryRowContext(ctx, "INSERT INTO chat.contacts (user1_id, user2_id, state_id) VALUES ($1, $2, $3) RETURNING id",
		userID1, userID2, 3).Scan(&entryID)
	if err != nil {
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method AddContact, users.go",
		}
		slog.Error(customErr.Error())
		return false
	}
	logger.Info("AddContact: success", "userID1", userID1, "userID2", userID2)
	return true
}
