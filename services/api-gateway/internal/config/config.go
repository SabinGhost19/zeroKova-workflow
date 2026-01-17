package config

import "os"

// Config holds the application configuration
type Config struct {
	ServerPort           string
	OrderServiceAddr     string
	InventoryServiceAddr string
	NotificationServiceAddr string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:              getEnv("SERVER_PORT", "8080"),
		OrderServiceAddr:        getEnv("ORDER_SERVICE_ADDR", "order-service:50051"),
		InventoryServiceAddr:    getEnv("INVENTORY_SERVICE_ADDR", "inventory-service:50052"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "notification-service:50053"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
