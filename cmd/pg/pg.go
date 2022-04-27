package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pdk/rec/pipe"
)

func main() {

	conn := dbConn()

	err := conn.Ping()
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	pipe.ReadSQL(conn, "select * from x").Print()
}

func dbConn() (db *sql.DB) {
	dbDriver := "pgx"
	dbUser := "pdk"
	dbPass := ""
	dbName := "pdk"
	dbHostname := "localhost"
	dbPort := "5432"

	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHostname, dbPort, dbName)
	db, err := sql.Open(dbDriver, connUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return db
}
