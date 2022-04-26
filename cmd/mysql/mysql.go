package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	conn := dbConn()

	err := conn.Ping()
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
}

func dbConn() (db *sql.DB) {
	// alias mysqltest='MYSQL_PWD=test123 mysql -h 127.0.0.1 -P 3306 -u root'
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "test123"
	dbName := "zapps"

	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return db
}
