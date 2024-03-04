package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Request struct {
	method  string
	url     string
	payLoad []byte
	problem string
}

type testCase struct {
	name    string
	request Request
}

func TestRegisterLoginLogout(t *testing.T) {
	api := NewMyHandler()
	api.ClearUserData()

	var emptyUsernameUser = map[string]interface{}{
		"username": "",
		"password": "123",
		"email":    "admin@mail.ru",
	}

	var validUser = map[string]interface{}{
		"username": "admin",
		"password": "123",
		"email":    "admin@mail.ru",
	}

	var unvalidUserPassword = map[string]interface{}{
		"username": "admin",
		"password": "1234",
		"email":    "admin@mail.ru",
	}

	var invalidJsonUser = map[string]interface{}{
		"username": "",
		"password": "123",
		"email":    "admin@mail.ru",
	}

	var userNotFound = map[string]interface{}{
		"username": "Somebody",
		"password": "password",
		"email":    "myEMAIL@mail.ru",
	}

	const (
		ProblemRawUser       = "raw_user"
		ProblemMethodGet     = "method_get"
		ProblemNotJSON       = "not_json"
		ProblemUserNotFound  = "user_not_found"
		ProblemWrongPassword = "wrong_passord"
		ProblemUserExists    = "user_already_exists"
	)

	testCases := []testCase{
		{
			name: "VaildUserRegistration",
			request: Request{
				method:  "POST",
				url:     "/register",
				payLoad: converToJSON(validUser),
			},
		},
		{
			name: "VaildUserLogin",
			request: Request{
				method:  "POST",
				url:     "/login",
				payLoad: converToJSON(validUser),
			},
		},
		{
			name: "EmptyUsernameRegistration",
			request: Request{
				method:  "POST",
				url:     "/register",
				payLoad: converToJSON(emptyUsernameUser),
				problem: ProblemRawUser,
			},
		},
		{
			name: "MethodGet",
			request: Request{
				method:  "GET",
				url:     "/register",
				payLoad: converToJSON(validUser),
				problem: ProblemMethodGet,
			},
		},
		{
			name: "notJSON",
			request: Request{
				method:  "POST",
				url:     "/register",
				payLoad: converToJSON(invalidJsonUser),
				problem: ProblemNotJSON,
			},
		},
		{
			name: "UserNotFound",
			request: Request{
				method:  "POST",
				url:     "/login",
				payLoad: converToJSON(userNotFound),
				problem: ProblemUserNotFound,
			},
		},
		{
			name: "WrongPassword",
			request: Request{
				method:  "POST",
				url:     "/login",
				payLoad: converToJSON(unvalidUserPassword),
				problem: ProblemWrongPassword,
			},
		},
		{
			name: "UserAlreadyExists",
			request: Request{
				method:  "POST",
				url:     "/register",
				payLoad: converToJSON(validUser),
				problem: ProblemUserExists,
			},
		},
		{
			name: "checkAuth",
			request: Request{
				method:  "POST",
				url:     "/checkAuth",
				payLoad: converToJSON(validUser),
				problem: ProblemUserExists,
			},
		},
		{
			name: "Logout",
			request: Request{
				method:  "POST",
				url:     "/logout",
				payLoad: nil,
			},
		},
	}

	sessionID := ""
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.request.method, tc.request.url, bytes.NewReader(tc.request.payLoad))
			if err != nil {
				t.Fatal(err)
			}
			if tc.request.problem == ProblemNotJSON {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			handler := http.HandlerFunc(api.Login)
			rr := httptest.NewRecorder()
			if tc.request.url == "/login" {
				handler = http.HandlerFunc(api.Login)
			} else if tc.request.url == "/register" {
				handler = http.HandlerFunc(api.Register)
			} else if tc.request.url == "/logout" {
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: sessionID,
				}
				req.AddCookie(cookie)
				handler = http.HandlerFunc(api.Logout)
			} else if tc.request.url == "/checkAuth" {
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: sessionID,
				}
				req.AddCookie(cookie)
				handler = http.HandlerFunc(api.CheckAuth)
			}
			handler.ServeHTTP(rr, req)

			status := rr.Code
			if tc.request.url == "/register" {
				cookies := rr.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						sessionID = cookie.Value
						break
					}
				}
			}
			responseBodyText := rr.Body.String()
			if status != http.StatusOK {
				if tc.request.problem == ProblemRawUser && rr.Body.String() == "{\"status\":400,\"body\":{\"error\":\"required field is empty\"}}" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for raw user data")
					return
				}
				if tc.request.problem == ProblemMethodGet && rr.Body.String() == "{\"status\":405,\"body\":{\"error\":\"use POST\"}}" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for another request method")
					return
				}
				if tc.request.problem == ProblemNotJSON && rr.Body.String() == "Content-Type header is not application/json\n" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for not JSON type")
					return
				}
				if tc.request.problem == ProblemUserNotFound && rr.Body.String() == "{\"status\":400,\"body\":{\"error\":\"user not found\"}}" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for not JSON type")
					return
				}
				if tc.request.problem == ProblemWrongPassword && rr.Body.String() == "{\"status\":400,\"body\":{\"error\":\"wrong password\"}}" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for not JSON type")
					return
				}
				if tc.request.problem == ProblemUserExists && rr.Body.String() == "{\"status\":400,\"body\":{\"error\":\"user already exists\"}}" {
					fmt.Println(tc.name, ": ------------- STATUS: OK")
					t.Skip("Expected error for not JSON type")

					return
				}
				t.Errorf("Login handler returned wrong status code: got %v want %v. Body: %v", status, http.StatusOK, responseBodyText)
			} else {
				fmt.Println(tc.name, ": ------------- STATUS: OK")
			}
		})
	}

}

func converToJSON(userData map[string]interface{}) []byte {
	body, err := json.Marshal(userData)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
