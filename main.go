package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hex-aragon/go-backend-boilerplate/api"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:boiler@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db",err)
	}
	fmt.Println("conn",conn)

	store := db.NewStore(conn)
	server := api.NewServer(store)

	fmt.Println("store",store)
	fmt.Println("server",server)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}