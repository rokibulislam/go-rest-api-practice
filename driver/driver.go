package driver

import (
	"database/sql"
	"log"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")

	if err != nil {
		panic(err)
	}

	return db
	// defer db.Close()
}
