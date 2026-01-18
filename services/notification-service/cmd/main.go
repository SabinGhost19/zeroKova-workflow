package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/database"
	"notification-service/internal/health"
	"notification-service/internal/service"
	pb "notification-service/proto"

	"google.golang.org/grpc"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	grpcPort := getEnv("GRPC_PORT", "50053")
	httpPort := getEnv("HTTP_PORT", "8083")

	// Start HTTP health server in goroutine
	healthServer := health.NewServer(":" + httpPort)
	go func() {
		if err := healthServer.Start(); err != nil {
			log.Fatalf("Failed to start health server: %v", err)
		}
	}()

	// Initialize database
	db, err := database.GetInstance()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	notificationService := service.NewNotificationServiceServer(db)
	pb.RegisterNotificationServiceServer(grpcServer, notificationService)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting gRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
