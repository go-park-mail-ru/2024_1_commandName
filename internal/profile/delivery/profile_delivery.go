package delivery

import (
	"ProjectMessenger/domain"
	authdelivery "ProjectMessenger/internal/auth/delivery"
	"ProjectMessenger/internal/misc"
	"ProjectMessenger/internal/profile/usecase"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProfileHandler struct {
	AuthHandler *authdelivery.AuthHandler
}

type updateUserStruct struct {
	User               domain.Person `json:"user"`
	NumOfUpdatedFields int           `json:"numUpdated"`
}

func NewProfileHandler(authHandler *authdelivery.AuthHandler) *ProfileHandler {
	return &ProfileHandler{AuthHandler: authHandler}
}

// GetProfileInfo gets profile info
//
// @Summary gets profile info
// @ID GetProfileInfo
// @Produce json
// @Success 200 {object}  domain.Response[domain.User]
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
	user.Password = ""
	user.PasswordSalt = ""

	misc.WriteStatusJson(w, 200, domain.User{User: user})
}

func (p *ProfileHandler) UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	authorized, _ := p.AuthHandler.CheckAuthNonAPI(w, r)
	if !authorized {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser updateUserStruct
	err := decoder.Decode(&jsonUser)
	if err != nil {
		misc.WriteStatusJson(w, 400, domain.Error{Error: "wrong json structure"})
	}

	if jsonUser.NumOfUpdatedFields > 0 {

	}

	fmt.Println(jsonUser)
	misc.WriteStatusJson(w, 200, nil)
}
