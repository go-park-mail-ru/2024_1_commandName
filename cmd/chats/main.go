package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	chats "ProjectMessenger/internal/chats_service/proto"
	"ProjectMessenger/internal/chats_service/repository"
	"ProjectMessenger/internal/chats_service/usecase"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
		fmt.Println("starting metrics server at :9092")
		log.Fatal(http.ListenAndServe(":9092", mux))
	}()
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	dataBase := СreateDatabase()
	chatStorage := repository.NewChatsStorage(dataBase)
	chats.RegisterChatServiceServer(server, usecase.NewChatManager(chatStorage))
	fmt.Println("starting server at :8082")
	server.Serve(lis)
}
