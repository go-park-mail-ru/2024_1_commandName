package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/misc"
	"ProjectMessenger/internal/profile/usecase"
)

type ProfileHandler struct {
	AuthHandler *authdelivery.AuthHandler
}

// Response[T]
type updateUserStruct[T any] struct {
	User               T   `json:"user"`
	NumOfUpdatedFields int `json:"numOfUpdatedFields"`
}

type changePasswordStruct struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type addContactStruct struct {
	UsernameOfUserToAdd string `json:"username_of_user_to_add"`
}

func NewProfileHandler(authHandler *authdelivery.AuthHandler) *ProfileHandler {
	return &ProfileHandler{AuthHandler: authHandler}
}

type docsUserForGetProfile struct {
	ID           uint      `json:"id" `
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	About        string    `json:"about"`
	CreateDate   time.Time `json:"create_date"`
	LastSeenDate time.Time `json:"last_seen_date"`
	Avatar       string    `json:"avatar"`
}

type docsContacts struct {
	Contacts []docsUserForGetProfile `json:"contacts"`
}

// GetProfileInfo gets profile info
//
// @Summary gets profile info
// @ID GetProfileInfo
// @Produce json
// @Success 200 {object}  domain.Response[docsUserForGetProfile]
// @Failure 400 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getProfileInfo [get]
func (p *ProfileHandler) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	user, found := usecase.GetProfileInfo(r.Context(), userID, p.AuthHandler.Users)
	if !found {
		misc.WriteInternalErrorJson(ctx, w)
		return
	}
	user.Password = ""
	user.PasswordSalt = ""

	misc.WriteStatusJson(ctx, w, 200, domain.User{User: user})
}

// UpdateProfileInfo updates profile info
//
// @Summary updates profile info
// @ID UpdateProfileInfo
// @Accept json
// @Produce json
// @Param userAndNumOfUpdatedFields body  updateUserStruct[docsUserForGetProfile] true "Send only the updated fields, and number of them"
// @Success 200 {object}  domain.Response[int]
// @Failure 401 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /updateProfileInfo [post]
func (p *ProfileHandler) UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser updateUserStruct[domain.Person]
	err := decoder.Decode(&jsonUser)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	if jsonUser.NumOfUpdatedFields <= 0 {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.UpdateProfileInfo(ctx, jsonUser.User, jsonUser.NumOfUpdatedFields, userID, p.AuthHandler.Users)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}

// ChangePassword changes profile password
//
// @Summary changes profile password
// @ID ChangePassword
// @Accept json
// @Produce json
// @Param Password body  changePasswordStruct true "Old and new passwords"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "passwords are empty"
// @Failure 401 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /changePassword [post]
func (p *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var passwordsJson changePasswordStruct
	err := decoder.Decode(&passwordsJson)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	if passwordsJson.OldPassword == "" || passwordsJson.NewPassword == "" {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "Поля пустые"})
		return
	}

	err = usecase.ChangePassword(r.Context(), passwordsJson.OldPassword, passwordsJson.NewPassword, userID, p.AuthHandler.Users)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		return
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}

// UploadAvatar uploads or changes avatar
//
// @Summary uploads or changes avatar
// @ID UploadAvatar
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "avatar image"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Описание ошибки"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /uploadAvatar [post]
func (p *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum avatar size
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "Недопустимый файл"})
		return
	}

	// Get the avatar from the request
	avatar, handler, err := r.FormFile("avatar")
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "Недопустимый файл"})
		return
	}
	defer avatar.Close()

	err = usecase.ChangeAvatar(r.Context(), avatar, handler, userID, p.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
		} else {
			misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		}
		return
	}

	misc.WriteStatusJson(ctx, w, 200, nil)
}

// GetContacts returns contacts of user
//
// @Summary returns contacts of user
// @ID GetContacts
// @Produce json
// @Success 200 {object}  domain.Response[docsContacts]
// @Failure 400 {object}  domain.Response[domain.Error] "Описание ошибки"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /getContacts [get]
func (p *ProfileHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	contacts := usecase.GetContacts(ctx, userID, p.AuthHandler.Users)
	misc.WriteStatusJson(ctx, w, 200, domain.Contacts{Contacts: contacts})
}

// AddContact adds contact for user
//
// @Summary adds contact for user
// @ID AddContact
// @Accept json
// @Produce json
// @Param usernameToAdd body addContactStruct true "username of user to add to contacts"
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "Описание ошибки"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /addContact [post]
func (p *ProfileHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var contact addContactStruct
	err := decoder.Decode(&contact)
	if err != nil {
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	err = usecase.AddContactByUsername(ctx, userID, contact.UsernameOfUserToAdd, p.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(ctx, w)
		} else {
			misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
			return
		}
	}
	misc.WriteStatusJson(ctx, w, 200, nil)
}
