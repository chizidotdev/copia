package app

import (
	"github.com/gin-gonic/gin"
)

// createRoutes creates all the routes for the server
func createRoutes(server *Server) {
	// Create root ping route
	server.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Copia!",
		})
	})

	server.router.Use(server.isAuth)
	// OrderService routes
	orderRoutes := server.router.Group("/orders")
	{
		orderRoutes.POST("", server.createOrder)
		orderRoutes.GET("", server.listOrders)
		orderRoutes.GET("/:id", server.getOrderByID)
		orderRoutes.PUT("/:id", server.updateOrder)
		orderRoutes.DELETE("/:id", server.deleteOrder)
	}

	// ReportService routes
	server.router.GET("/inventory", server.getReport)
}
