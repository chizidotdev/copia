package adapters

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgresClient(connString string) *sql.DB {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Cannot connect to postgres database: ", err)
	}

	return conn
}
