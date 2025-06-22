package db

import (
	"database/sql"
	"log"
)

func GetConnection() *sql.DB {
	connection, err := sql.Open("sqlite3", "")

	if err != nil {
		log.Fatal(err)
	}

	return connection
}
