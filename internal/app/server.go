package app

import (
	"github.com/chizidotdev/copia/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	*service.Service
}

// NewServer creates a new HTTP server and setup routing
func NewServer(service *service.Service) *Server {
	router := gin.Default()

	server := &Server{
		router:  router,
		Service: service,
	}

	corsConfig(server)
	createRoutes(server)

	//if err := server.Store.ConnDB.Ping(); err != nil {
	//	log.Fatal("Cannot connect to db:", err)
	//}

	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// corsConfig sets up the CORS configuration
func corsConfig(server *Server) {
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://copia.aidmedium.com"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}
