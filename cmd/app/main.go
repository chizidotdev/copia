package main

import (
	"log"

	"github.com/chizidotdev/shop/api"
	"github.com/chizidotdev/shop/config"
)

func main() {
	config.LoadConfig()

	server := api.NewHTTPServer()

	port := config.EnvVars.PORT
	if port == "" {
		port = "5000"
	}
	serverAddr := "0.0.0.0:" + port
	err := server.Start(serverAddr)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
