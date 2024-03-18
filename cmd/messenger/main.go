package main

import (
	"ProjectMessenger/internal/auth/usecase"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"ProjectMessenger/internal/middleware"
)

var DEBUG = false

func main() {
	Router()
}

// Router
// @Title Messenger authorization API
// @Version 1.0
// @schemes http
// @host localhost:8080
// @BasePath  /
func Router() {
	r := mux.NewRouter()

	//api := auth.NewMyHandler(DEBUG)
	newApi := usecase.NewAuthHandler()

	r.HandleFunc("/checkAuth", newApi.CheckAuth)
	r.HandleFunc("/login", newApi.Login)
	r.HandleFunc("/logout", newApi.Logout)
	r.HandleFunc("/register", newApi.Register)
	r.HandleFunc("/getChats", newApi.GetChats)

	// middleware
	if DEBUG {
		r.Use(middleware.CORS)
	}

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("err")
		log.Fatal(err)
		return
	}
}
