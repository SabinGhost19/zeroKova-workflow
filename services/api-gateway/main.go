// API Gateway - Entry point for all client requests
// Version: 1.0.0
package main

import (
	"log"
	"os"

	"github.com/test-workflow/api-gateway/internal/config"
	"github.com/test-workflow/api-gateway/internal/handlers"
	"github.com/test-workflow/api-gateway/internal/grpc/clients"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize gRPC clients
	grpcClients, err := clients.NewClients(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer grpcClients.Close()

	// Initialize Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder(grpcClients))
			orders.GET("", handlers.ListOrders(grpcClients))
			orders.GET("/:id", handlers.GetOrder(grpcClients))
			orders.PUT("/:id/status", handlers.UpdateOrderStatus(grpcClients))
		}

		// Inventory routes
		inventory := api.Group("/inventory")
		{
			inventory.GET("/products", handlers.ListProducts(grpcClients))
			inventory.POST("/products", handlers.AddProduct(grpcClients))
			inventory.GET("/products/:id", handlers.GetStock(grpcClients))
			inventory.PUT("/products/:id/stock", handlers.UpdateStock(grpcClients))
		}

		// Notifications routes
		notifications := api.Group("/notifications")
		{
			notifications.GET("", handlers.GetNotifications(grpcClients))
		}
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
