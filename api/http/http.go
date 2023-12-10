package http

import (
	"github.com/chizidotdev/copia/api/http/middleware"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router *gin.Engine
	*usecases.UserService
	*usecases.ProductService
	*usecases.OrderService
}

// NewHTTPServer creates a new HTTP server and sets up routing
func NewHTTPServer(
	userService *usecases.UserService,
	productService *usecases.ProductService,
	orderService *usecases.OrderService,
) *Server {
	router := gin.Default()
	store := cookie.NewStore([]byte(config.EnvVars.AuthSecret))
	store.Options(sessions.Options{
		MaxAge:   86400 * 30, // 30 days
		Secure:   false,
		HttpOnly: true,
		Domain:   config.EnvVars.CookieDomain,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	router.Use(sessions.Sessions("copia_auth", store))

	server := &Server{
		router:         router,
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}

	corsConfig(server)
	createRoutes(server)

	return server
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// corsConfig sets up the CORS configuration
func corsConfig(server *Server) {
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:3000",
			"https://copia.aidmedium.com",
			"https://copia.up.railway.app",
		},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

// createRoutes creates all the routes for the server
func createRoutes(server *Server) {
	// Create root ping route
	server.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Copia!",
		})
	})

	// User routes
	userHandler := NewUserHandler(server.UserService)
	server.router.GET("/user", middleware.IsAuthenticated, userHandler.getUser)
	server.router.POST("/register", userHandler.createUser)
	server.router.POST("/login", userHandler.login)
	server.router.GET("/login/google", userHandler.loginWithSSO)
	server.router.GET("/callback", userHandler.ssoCallback)
	server.router.POST("/logout", userHandler.logout)

	server.router.POST("/send-verification-email", userHandler.sendVerificationEmail)
	server.router.POST("/verify-email", userHandler.verifyEmail)
	server.router.POST("/reset-password", userHandler.resetPassword)
	server.router.POST("/change-password", userHandler.changePassword)

	// Product routes
	productHandler := NewProductHandler(server.ProductService)
	productRoutes := server.router.Group("/products")
	productRoutes.Use(middleware.IsAuthenticated)
	{
		productRoutes.POST("", productHandler.createProduct)
		productRoutes.GET("", productHandler.listProducts)
		productRoutes.GET("/:id", productHandler.getProduct)
		productRoutes.DELETE("/:id", productHandler.deleteProduct)

		productRoutes.PUT("/:id", productHandler.updateProduct)
		productRoutes.PATCH("/:id/image", productHandler.updateProductImage)
		productRoutes.PATCH("/:id/quantity", productHandler.updateProductQuantity)

		productRoutes.GET("/settings", productHandler.getProductSettings)
		productRoutes.PATCH("/settings", productHandler.updateProductSettings)
	}

	// Order routes
	orderHandler := NewOrderHandler(server.OrderService)
	orderRoutes := server.router.Group("/orders")
	orderRoutes.Use(middleware.IsAuthenticated)
	{
		orderRoutes.POST("", orderHandler.createOrder)
		orderRoutes.GET("", orderHandler.listOrders)
		orderRoutes.GET("/:id", orderHandler.getOrder)
		orderRoutes.DELETE("/:id", orderHandler.deleteOrder)
	}

	// ReportService routes
	//server.router.GET("/report", server.getReport)
}
