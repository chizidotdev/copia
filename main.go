package main

import (
	"github.com/chizidotdev/copia/internal/app"
	"github.com/chizidotdev/copia/internal/repository"
	"github.com/chizidotdev/copia/internal/service"
	"github.com/chizidotdev/copia/pkg/utils"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	utils.LoadConfig()

	conn, err := gorm.Open(postgres.Open(utils.EnvVars.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := repository.NewStore(conn)
	newService := service.NewService(store)
	server := app.NewServer(newService)

	port := utils.EnvVars.PORT
	if port == "" {
		port = "8080"
	}
	serverAddr := "0.0.0.0:" + port

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
