package adapters

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewPostgresClient(connString string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Cannot connect to postgres database: ", err)
	}

	return conn
}
