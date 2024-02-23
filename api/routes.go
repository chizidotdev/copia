package api

import "github.com/gin-gonic/gin"

// createRoutes creates all the routes for the server
func createRoutes(server *Server) {
	// Create root ping route
	server.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Aidmedium",
		})
	})

	// User routes
	// userHandler := NewUserHandler(server.UserService)
	// server.router.GET("/user", middleware.IsAuthenticated, userHandler.getUser)
	// server.router.POST("/register", userHandler.createUser)
	// server.router.POST("/login", userHandler.login)
	// server.router.GET("/login/google", userHandler.loginWithSSO)
	// server.router.GET("/callback", userHandler.ssoCallback)
	// server.router.POST("/logout", userHandler.logout)

	// server.router.POST("/send-verification-email", userHandler.sendVerificationEmail)
	// server.router.POST("/verify-email", userHandler.verifyEmail)
	// server.router.POST("/reset-password", userHandler.resetPassword)
	// server.router.POST("/change-password", userHandler.changePassword)

	// // Product routes
	// productHandler := NewProductHandler(server.ProductService)
	// productRoutes := server.router.Group("/products")
	// productRoutes.Use(middleware.IsAuthenticated)
	// {
	// 	productRoutes.POST("", productHandler.createProduct)
	// 	productRoutes.GET("", productHandler.listProducts)
	// 	productRoutes.GET("/:id", productHandler.getProduct)
	// 	productRoutes.DELETE("/:id", productHandler.deleteProduct)

	// 	productRoutes.PUT("/:id", productHandler.updateProduct)
	// 	productRoutes.PATCH("/:id/image", productHandler.updateProductImage)
	// 	productRoutes.PATCH("/:id/quantity", productHandler.updateProductQuantity)

	// 	productRoutes.GET("/settings", productHandler.getProductSettings)
	// 	productRoutes.PATCH("/settings", productHandler.updateProductSettings)
	// }

	// // Order routes
	// orderHandler := NewOrderHandler(server.OrderService)
	// orderRoutes := server.router.Group("/orders")
	// orderRoutes.Use(middleware.IsAuthenticated)
	// {
	// 	orderRoutes.POST("", orderHandler.createOrder)
	// 	orderRoutes.GET("", orderHandler.listOrders)
	// 	orderRoutes.GET("/:id", orderHandler.getOrder)
	// 	orderRoutes.DELETE("/:id", orderHandler.deleteOrder)
	// }

	// ReportService routes
	//server.router.GET("/report", server.getReport)
}
