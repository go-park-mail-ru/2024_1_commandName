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
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	user, found := usecase.GetProfileInfo(r.Context(), userID, p.AuthHandler.Users)
	if !found {
		misc.WriteInternalErrorJson(w)
		return
	}
	user.ID = 0
	user.Password = ""
	user.PasswordSalt = ""

	misc.WriteStatusJson(w, 200, domain.User{User: user})
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
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser updateUserStruct[domain.Person]
	err := decoder.Decode(&jsonUser)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	if jsonUser.NumOfUpdatedFields <= 0 {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
		return
	}

	err = usecase.UpdateProfileInfo(r.Context(), jsonUser.User, jsonUser.NumOfUpdatedFields, userID, p.AuthHandler.Users)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		return
	}
	misc.WriteStatusJson(w, 200, nil)
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
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var passwordsJson changePasswordStruct
	err := decoder.Decode(&passwordsJson)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
		return
	}
	if passwordsJson.OldPassword == "" || passwordsJson.NewPassword == "" {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "Поля пустые"})
		return
	}

	err = usecase.ChangePassword(r.Context(), passwordsJson.OldPassword, passwordsJson.NewPassword, userID, p.AuthHandler.Users)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		return
	}
	misc.WriteStatusJson(w, 200, nil)
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
	authorized, userID := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}
	if r.Method != http.MethodPost {
		misc.WriteStatusJson(w, 405, domain.Error{Error: "use POST"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum avatar size
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "Недопустимый файл"})
		return
	}

	// Get the avatar from the request
	avatar, handler, err := r.FormFile("avatar")
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "Недопустимый файл"})
		return
	}
	defer avatar.Close()

	err = usecase.ChangeAvatar(r.Context(), avatar, handler, userID, p.AuthHandler.Users)
	if err != nil {
		if err.Error() == "internal error" {
			misc.WriteInternalErrorJson(w)
		} else {
			misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		}
		return
	}

	misc.WriteStatusJson(w, 200, nil)
}
