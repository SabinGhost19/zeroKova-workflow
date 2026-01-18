package service

import (
	"context"
	"fmt"
	"log"

	"notification-service/internal/database"
	pb "notification-service/proto"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
	db *database.DB
}

func NewNotificationServiceServer(db *database.DB) *NotificationServiceServer {
	return &NotificationServiceServer{db: db}
}

func (s *NotificationServiceServer) SendOrderNotification(ctx context.Context, req *pb.OrderNotificationRequest) (*pb.StatusResponse, error) {
	message := fmt.Sprintf("Order %s for %s - Status: %s, Total: $%.2f",
		req.OrderId, req.CustomerName, req.Status, req.TotalAmount)

	id, err := s.db.SaveNotification("ORDER", message)
	if err != nil {
		log.Printf("Failed to save order notification: %v", err)
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send notification: %v", err),
		}, nil
	}

	log.Printf("Order notification sent: %s", id)
	return &pb.StatusResponse{
		Success: true,
		Message: fmt.Sprintf("Notification sent successfully with ID: %s", id),
	}, nil
}

func (s *NotificationServiceServer) SendStockAlert(ctx context.Context, req *pb.StockAlertRequest) (*pb.StatusResponse, error) {
	message := fmt.Sprintf("Stock Alert [%s]: %s - Current quantity: %d",
		req.AlertType, req.ProductName, req.CurrentQuantity)

	id, err := s.db.SaveNotification("STOCK_ALERT", message)
	if err != nil {
		log.Printf("Failed to save stock alert: %v", err)
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send alert: %v", err),
		}, nil
	}

	log.Printf("Stock alert sent: %s", id)
	return &pb.StatusResponse{
		Success: true,
		Message: fmt.Sprintf("Alert sent successfully with ID: %s", id),
	}, nil
}

func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error) {
	notifications, total, err := s.db.GetNotifications(req.Limit, req.Offset)
	if err != nil {
		log.Printf("Failed to get notifications: %v", err)
		return &pb.NotificationsResponse{
			Notifications: []*pb.Notification{},
			Total:         0,
		}, nil
	}

	pbNotifications := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		pbNotifications[i] = &pb.Notification{
			Id:        n.ID,
			Type:      n.Type,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Sent:      n.Sent,
		}
	}

	return &pb.NotificationsResponse{
		Notifications: pbNotifications,
		Total:         total,
	}, nil
}
