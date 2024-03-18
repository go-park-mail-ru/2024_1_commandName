package main

import (
	"ProjectMessenger/internal/auth/usecase"
	"fmt"
	"log"
	"net/http"

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
	authHandler := usecase.NewAuthHandler()

	// middleware
	if DEBUG {
		authHandler.Rt.Use(middleware.CORS)
	}

	err := http.ListenAndServe(":8080", authHandler.Rt)
	if err != nil {
		fmt.Println("err")
		log.Fatal(err)
		return
	}
}
