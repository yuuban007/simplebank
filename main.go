package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yuuban007/simplebank/api"
	db "github.com/yuuban007/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:gcc123456@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start("0.0.0.0:8888")
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
