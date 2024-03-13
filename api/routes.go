package api

import (
	"github.com/chizidotdev/shop/api/cart"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/api/order"
	"github.com/chizidotdev/shop/api/product"
	"github.com/chizidotdev/shop/api/seed"
	"github.com/chizidotdev/shop/api/store"
	"github.com/chizidotdev/shop/api/user"
	"github.com/gin-gonic/gin"
)

// createRoutes creates all the routes for the server
func createRoutes(server *Server) {
	router := server.router

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
	orderHandler := order.NewOrderHandler(server.pgStore)

	seedHandler := seed.NewSeedHandler(server.pgStore)

	// Auth routes
	router.POST("/register", userHandler.CreateUser)
	router.POST("/login", userHandler.Login)
	router.GET("/login/google", userHandler.GoogleLogin)
	router.GET("/callback", userHandler.GoogleCallback)
	router.POST("/logout", userHandler.Logout)
	router.GET("/users/me", middleware.IsAuthenticated, userHandler.GetUser)

	{
		storeRouter := router.Group("/stores")
		{
			storeRouter.POST("/search", storeHandler.SearchStores)
			storeRouter.GET("", storeHandler.ListStores)
			storeRouter.GET("/:storeID", storeHandler.GetStore)
			storeRouter.GET("/:storeID/products", productHandler.ListStoreProducts)
			router.GET("/products/:productID", productHandler.GetProduct)
		}

		orderRouter := router.Group("/orders")
		{
			orderRouter.Use(middleware.IsAuthenticated)
			orderRouter.POST("", orderHandler.CreateOrder)
			orderRouter.GET("", orderHandler.ListUserOrders)
			// orderRoutes.DELETE("/:id", orderHandler.DeleteOrder)
		}

		cartRoutes := router.Group("/cart")
		{
			cartRoutes.Use(middleware.IsAuthenticated)
			cartRoutes.GET("", cartHandler.GetCart)
			cartRoutes.POST("", cartHandler.AddToCart)
			cartRoutes.PATCH("/:cartID", cartHandler.UpdateCart)
			cartRoutes.DELETE("/:cartID", cartHandler.DeleteCart)
		}
	}

	// vendor routes
	vendorRouter := router.Group("/vendor")
	{
		vendorRouter.Use(middleware.IsVendor)

		// vendor store
		vendorStoreRouter := vendorRouter.Group("/stores")
		{
			vendorStoreRouter.GET("", storeHandler.GetUserStore)
			vendorStoreRouter.PUT("", storeHandler.UpdateUserStore)
			vendorStoreRouter.POST("", storeHandler.CreateStore)
		}

		// vendor products
		vendorProductRouter := vendorRouter.Group("/products")
		{
			vendorProductRouter.POST("/seed", seedHandler.SeedStore)
			vendorProductRouter.GET("", productHandler.ListStoreProducts)
			vendorProductRouter.POST("", productHandler.CreateProduct)
			vendorProductRouter.PATCH("/:productID", productHandler.UpdateProduct)
			vendorProductRouter.DELETE("/:productID", productHandler.DeleteProduct)
		}

		// vendor orders
		vendorOrderRouter := vendorRouter.Group("/orders")
		{
			vendorOrderRouter.GET("", orderHandler.ListStoreOrders)
			vendorOrderRouter.PATCH("/:orderItemID", orderHandler.UpdateOrderItem)
		}
	}
}
