package main

import (
	"github.com/chizidotdev/copia/api/http"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/adapters"
	"github.com/chizidotdev/copia/internal/app/usecases"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	config.LoadConfig()

	conn, err := gorm.Open(postgres.Open(config.EnvVars.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	userRepo := adapters.NewUserRepository(conn)
	userService := usecases.NewUserService(userRepo)
	server := http.NewHTTPServer(userService)

	port := config.EnvVars.PORT
	if port == "" {
		port = "8080"
	}
	serverAddr := "0.0.0.0:" + port

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
