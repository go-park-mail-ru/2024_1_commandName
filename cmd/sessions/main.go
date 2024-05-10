package main

import (
	session "ProjectMessenger/internal/sessions_service/proto"
	"ProjectMessenger/internal/sessions_service/repository"
	"ProjectMessenger/internal/sessions_service/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
)

func СreateDatabase() *sql.DB {
	connStrToDataBase := "user=chatme_user dbname=chatme password=EasyPassword( host=localhost port=8888 sslmode=disable"
	dataBase, err := sql.Open("postgres", connStrToDataBase)
	if err != nil {
		fmt.Println("DatBase open err:", err)
		return nil
	}

	err = dataBase.Ping()
	if err != nil {
		fmt.Println("connection to DatBase err:", err)
		return nil
	}
	return dataBase
}

func main() {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		fmt.Println("starting metrics server at :9091")
		log.Fatal(http.ListenAndServe(":9091", mux))
	}()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()
	dataBase := СreateDatabase()
	sessionStorage := repository.NewSessionStorage(dataBase)
	session.RegisterAuthCheckerServer(server, usecase.NewSessionManager(sessionStorage))

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
