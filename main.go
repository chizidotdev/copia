package main

import (
	"github.com/chizidotdev/copia/app"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/service"
	"github.com/chizidotdev/copia/util"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	util.LoadConfig()

	conn, err := gorm.Open(postgres.Open(util.EnvVars.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := repository.NewStore(conn)
	newService := service.NewService(store)
	server := app.NewServer(newService)

	port := util.EnvVars.PORT
	if port == "" {
		port = "8080"
	}
	serverAddr := "0.0.0.0:" + port

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
