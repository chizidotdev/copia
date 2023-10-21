package api

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

	server.router.POST("/user", server.createUser)
	server.router.POST("/login", server.login)
	server.router.GET("/login/google", server.loginWithSSO)
	server.router.GET("/callback", server.ssoCallback)
	server.router.GET("/logout", server.logout)
	server.router.GET("/user", server.isAuthenticated, server.getUser)

	//server.router.Use(server.isAuthenticated)
	// OrderService routes
	orderRoutes := server.router.Group("/orders")
	orderRoutes.Use(server.isAuthenticated)
	{
		orderRoutes.POST("", server.handleCreateOrder)
		orderRoutes.GET("", server.listOrders)
		orderRoutes.GET("/:id", server.handleGetOrderByID)
		orderRoutes.PUT("/:id", server.updateOrder)
		orderRoutes.DELETE("/:id", server.deleteOrder)
	}

	// ReportService routes
	server.router.GET("/report", server.getReport)
}
