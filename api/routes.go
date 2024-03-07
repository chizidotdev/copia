package api

import (
	"github.com/chizidotdev/shop/api/cart"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/api/product"
	"github.com/chizidotdev/shop/api/seed"
	"github.com/chizidotdev/shop/api/store"
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

	userHandler := user.NewUserHandler(server.pgStore)
	storeHandler := store.NewStoreHandler(server.pgStore)
	productHandler := product.NewProductHandler(server.pgStore)
	cartHandler := cart.NewCartHandler(server.pgStore)
	seedHandler := seed.NewSeedHandler(server.pgStore)

	// Auth routes
	server.router.POST("/register", userHandler.CreateUser)
	server.router.POST("/login", userHandler.Login)
	server.router.GET("/login/google", userHandler.GoogleLogin)
	server.router.GET("/callback", userHandler.GoogleCallback)
	server.router.POST("/logout", userHandler.Logout)

	// user routes
	userRoutes := server.router.Group("/users")
	{
		userRoutes.Use(middleware.IsAuthenticated)
		userRoutes.GET("/me", userHandler.GetUser)

		// user store routes
		userStoreRoutes := userRoutes.Group("/store")
		{
			userStoreRoutes.GET("", storeHandler.GetUserStore)
			userStoreRoutes.PUT("", storeHandler.UpdateUserStore)
			userStoreRoutes.POST("", storeHandler.CreateStore)
		}

		userCartRoutes := userRoutes.Group("/cart")
		{
			userCartRoutes.GET("", cartHandler.GetCart)
			userCartRoutes.POST("", cartHandler.AddToCart)
			userCartRoutes.PATCH("/:cartID", cartHandler.UpdateCart)
			userCartRoutes.DELETE("/:cartID", cartHandler.DeleteCart)
		}
	}

	storeRoutes := server.router.Group("/stores")
	{
		storeRoutes.GET("", storeHandler.ListStores)
		storeRoutes.GET("/:storeID", storeHandler.GetStore)

		storeRoutes.POST("/:storeID/seed", middleware.IsAuthenticated, seedHandler.SeedStore)

		// store product routes
		storeProductRoutes := storeRoutes.Group("/:storeID/products")
		{
			storeProductRoutes.GET("", productHandler.ListUserProducts)
			storeProductRoutes.Use(middleware.IsAuthenticated)
			storeProductRoutes.POST("", productHandler.CreateProduct)
			storeProductRoutes.PATCH("/:productID", productHandler.UpdateProduct)
			storeProductRoutes.DELETE("/:productID", productHandler.DeleteProduct)
		}
	}

	// Product routes
	productRoutes := server.router.Group("/products")
	{
		productRoutes.GET("/:productID", productHandler.GetProduct)
	}

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
