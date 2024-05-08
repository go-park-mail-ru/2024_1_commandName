package main

import (
	contacts "ProjectMessenger/microservices/contacts_service/proto"
	"ProjectMessenger/microservices/contacts_service/repository"
	"ProjectMessenger/microservices/contacts_service/usecase"
	"database/sql"
	"fmt"
	"log"
	"net"

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
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()
	dataBase := СreateDatabase()
	contactsStorage := repository.NewContactsStorage(dataBase)
	contacts.RegisterContactsServer(server, usecase.NewContactsManager(contactsStorage))

	fmt.Println("starting server at :8083")
	server.Serve(lis)
}
