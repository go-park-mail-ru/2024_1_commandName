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

func TestRegisterLoginLogout(t *testing.T) {
	var user1 = map[string]interface{}{
		"username": "admin",
		"password": "123",
		"email":    "admin@mail.ru",
	}

	var user2 = map[string]interface{}{
		"username": "",
		"password": "123",
		"email":    "admin@mail.ru",
	}

	var user3 = map[string]interface{}{
		"username": "admin3",
		"password": "12345",
		"email":    "admin3@mail.ru",
	}
	fmt.Println(user3)

	api := NewMyHandler()
	api.ClearUserData()
	testRequests := make([]Request, 0)

	bodyUser1, err := json.Marshal(user1)
	if err != nil {
		log.Fatal(err)
		return
	}

	bodyUser2, err := json.Marshal(user2)
	if err != nil {
		log.Fatal(err)
		return
	}

	bodyUser3, err := json.Marshal(user2)
	if err != nil {
		log.Fatal(err)
		return
	}

	testRequests = append(testRequests, Request{method: "POST", url: "/register", payLoad: bodyUser1})
	testRequests = append(testRequests, Request{method: "POST", url: "/login", payLoad: bodyUser1})
	testRequests = append(testRequests, Request{method: "POST", url: "/logout", payLoad: nil})
	testRequests = append(testRequests, Request{method: "POST", url: "/register", payLoad: bodyUser2, problem: "raw_user"})
	testRequests = append(testRequests, Request{method: "GET", url: "/register", payLoad: bodyUser2, problem: "method_get"})
	testRequests = append(testRequests, Request{method: "POST", url: "/register", payLoad: bodyUser3, problem: "not_json"})
	sessionID := ""
	for i := range testRequests {
		t.Run("SuccessfulRegInOut", func(t *testing.T) {
			req, err := http.NewRequest(testRequests[i].method, testRequests[i].url, bytes.NewReader(testRequests[i].payLoad))
			if err != nil {
				t.Fatal(err)
			}
			if testRequests[i].problem == "not_json" {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			handler := http.HandlerFunc(api.Login)
			rr := httptest.NewRecorder()
			if testRequests[i].url == "/login" {
				handler = http.HandlerFunc(api.Login)
			} else if testRequests[i].url == "/register" {
				handler = http.HandlerFunc(api.Register)
			} else if testRequests[i].url == "/logout" {
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: sessionID,
				}
				req.AddCookie(cookie)
				handler = http.HandlerFunc(api.Logout)
			}
			handler.ServeHTTP(rr, req)

			status := rr.Code
			if testRequests[i].url == "/register" {
				cookies := rr.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						sessionID = cookie.Value
						break
					}
				}
			}

			text := rr.Body.String()

			if status != http.StatusOK {
				if testRequests[i].problem == "raw_user" && rr.Body.String() == "{\"status\":400,\"body\":{\"error\":\"required field is empty\"}}" {
					t.Skip("Expected error for raw user data")
					return
				}
				if testRequests[i].problem == "method_get" && rr.Body.String() == "{\"status\":405,\"body\":{\"error\":\"use POST\"}}" {
					t.Skip("Expected error for another request method")
					return
				}
				if testRequests[i].problem == "not_json" && rr.Body.String() == "Content-Type header is not application/json\n" {
					t.Skip("Expected error for not JSON type")
					return
				}
				t.Errorf("Login handler returned wrong status code: got %v want %v. Body: %v", status, http.StatusOK, text)
			}
		})
	}
}
