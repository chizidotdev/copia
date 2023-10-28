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

	conn, err := gorm.Open(postgres.Open(config.EnvVars.DBSource), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	userRepo := adapters.NewUserRepository(conn)
	userService := usecases.NewUserService(userRepo)

	s3Repo := adapters.NewS3Repository()
	productRepo := adapters.NewProductRepository(conn)
	productService := usecases.NewProductService(productRepo, s3Repo)

	orderRepo := adapters.NewOrderRepository(conn)
	orderService := usecases.NewOrderService(orderRepo)

	server := http.NewHTTPServer(
		userService,
		productService,
		orderService,
	)

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
