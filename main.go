package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hex-aragon/go-backend-boilerplate/api"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
	"github.com/hex-aragon/go-backend-boilerplate/util"
	_ "github.com/lib/pq"
)


func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db",err)
	}
	fmt.Println("conn",conn)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	fmt.Println("store",store)
	fmt.Println("server",server)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}