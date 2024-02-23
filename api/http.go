package api

import (
	"net/http"
	"net/url"

	"github.com/chizidotdev/shop/config"
	"github.com/chizidotdev/shop/repository"
	"github.com/chizidotdev/shop/repository/adapters"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	pgStore    *repository.Queries
	redisStore *adapters.RedisClient
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// NewHTTPServer creates a new HTTP server and sets up routing
func NewHTTPServer() *Server {
	router := gin.Default()

	parsedURL, err := url.Parse(config.EnvVars.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisAddress := parsedURL.Host
	redisPassword, _ := parsedURL.User.Password()
	store, err := redis.NewStore(10, "tcp", redisAddress, redisPassword, []byte(config.EnvVars.AuthSecret))
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{
		MaxAge:   86400 * 30, // 30 days
		Secure:   false,
		HttpOnly: true,
		Domain:   config.EnvVars.CookieDomain,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	router.Use(sessions.Sessions("am_store_auth", store))

	redisStore := adapters.NewRedisClient(config.EnvVars.RedisUrl)
	pgConn := adapters.NewPostgresClient(config.EnvVars.DBSource)
	pgStore := repository.New(pgConn)

	server := &Server{
		router:     router,
		pgStore:    pgStore,
		redisStore: redisStore,
	}

	corsConfig(server)
	createRoutes(server)

	return server
}

// corsConfig sets up the CORS configuration
func corsConfig(server *Server) {
	server.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://copia.aidmedium.com",
			"https://copia.up.railway.app",
		},
		AllowMethods:     []string{"PUT", "POST", "GET", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}
