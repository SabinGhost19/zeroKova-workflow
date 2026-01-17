package clients

import (
	"log"

	"github.com/test-workflow/api-gateway/internal/config"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/inventory"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/notification"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients holds all gRPC client connections
type Clients struct {
	OrderClient        order.OrderServiceClient
	InventoryClient    inventory.InventoryServiceClient
	NotificationClient notification.NotificationServiceClient

	orderConn        *grpc.ClientConn
	inventoryConn    *grpc.ClientConn
	notificationConn *grpc.ClientConn
}

// NewClients creates new gRPC client connections
func NewClients(cfg *config.Config) (*Clients, error) {
	clients := &Clients{}

	// Connect to Order Service
	orderConn, err := grpc.Dial(
		cfg.OrderServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	clients.orderConn = orderConn
	clients.OrderClient = order.NewOrderServiceClient(orderConn)
	log.Printf("Connected to Order Service at %s", cfg.OrderServiceAddr)

	// Connect to Inventory Service
	inventoryConn, err := grpc.Dial(
		cfg.InventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	clients.inventoryConn = inventoryConn
	clients.InventoryClient = inventory.NewInventoryServiceClient(inventoryConn)
	log.Printf("Connected to Inventory Service at %s", cfg.InventoryServiceAddr)

	// Connect to Notification Service
	notificationConn, err := grpc.Dial(
		cfg.NotificationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	clients.notificationConn = notificationConn
	clients.NotificationClient = notification.NewNotificationServiceClient(notificationConn)
	log.Printf("Connected to Notification Service at %s", cfg.NotificationServiceAddr)

	return clients, nil
}

// Close closes all gRPC connections
func (c *Clients) Close() {
	if c.orderConn != nil {
		c.orderConn.Close()
	}
	if c.inventoryConn != nil {
		c.inventoryConn.Close()
	}
	if c.notificationConn != nil {
		c.notificationConn.Close()
	}
}
