package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/test-workflow/api-gateway/internal/grpc/clients"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/notification"

	"github.com/gin-gonic/gin"
)

// GetNotifications handles GET /api/v1/notifications
func GetNotifications(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

		grpcReq := &notification.GetNotificationsRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.NotificationClient.GetNotifications(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
