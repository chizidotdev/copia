package api

import (
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/api/user"
	"github.com/gin-gonic/gin"
)

// createRoutes creates all the routes for the server
func createRoutes(server *Server) {
	// Create root ping route
	server.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Aidmedium",
		})
	})

	// User routes
	userHandler := user.NewUserHandler(server.pgStore)
	server.router.GET("/user", middleware.IsAuthenticated, userHandler.GetUser)
	server.router.GET("/login/google", userHandler.GoogleLogin)
	server.router.GET("/callback", userHandler.GoogleCallback)
	server.router.POST("/logout", userHandler.Logout)

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
