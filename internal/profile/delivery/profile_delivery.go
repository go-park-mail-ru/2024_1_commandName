package delivery

import (
	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/misc"
	"ProjectMessenger/internal/profile/usecase"
	"encoding/json"
	"net/http"
	"time"
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
	user, found := usecase.GetProfileInfo(userID, p.AuthHandler.Users)
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

	err = usecase.UpdateProfileInfo(jsonUser.User, jsonUser.NumOfUpdatedFields, userID, p.AuthHandler.Users)
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
		misc.WriteStatusJson(w, 400, domain.Error{Error: "passwords are empty"})
		return
	}

	err = usecase.ChangePassword(passwordsJson.OldPassword, passwordsJson.NewPassword, userID, p.AuthHandler.Users)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: err.Error()})
		return
	}
	misc.WriteStatusJson(w, 200, nil)
}
