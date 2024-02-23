package main

import (
	"log"
	"os"

	"github.com/chizidotdev/shop/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config.LoadConfig()
	dbSource := config.EnvVars.DBSource

	m, err := migrate.New("file://repository/migrations", dbSource)

	if err != nil {
		log.Fatal("Init error", err)
	}

	arg := os.Args[1]

	if arg == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration error", err)
		}

		log.Println("Migration down successful")
		return
	}

	if arg == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration error", err)
		}

		log.Println("Migration up successful")
		return
	}
}
