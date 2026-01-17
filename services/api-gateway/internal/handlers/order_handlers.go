package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/test-workflow/api-gateway/internal/grpc/clients"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/order"

	"github.com/gin-gonic/gin"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerName string      `json:"customer_name" binding:"required"`
	Items        []OrderItem `json:"items" binding:"required,min=1"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID   string  `json:"product_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int32   `json:"quantity" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

// UpdateStatusRequest represents the request body for updating order status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// CreateOrder handles POST /api/v1/orders
func CreateOrder(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateOrderRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert to gRPC request
		grpcItems := make([]*order.OrderItem, len(req.Items))
		for i, item := range req.Items {
			grpcItems[i] = &order.OrderItem{
				ProductId:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Price:       item.Price,
			}
		}

		grpcReq := &order.CreateOrderRequest{
			CustomerName: req.CustomerName,
			Items:        grpcItems,
		}

		// Call Order Service
		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.OrderClient.CreateOrder(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

// GetOrder handles GET /api/v1/orders/:id
func GetOrder(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderID := ctx.Param("id")

		grpcReq := &order.GetOrderRequest{
			OrderId: orderID,
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.OrderClient.GetOrder(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !resp.Success {
			ctx.JSON(http.StatusNotFound, gin.H{"error": resp.Message})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// ListOrders handles GET /api/v1/orders
func ListOrders(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

		grpcReq := &order.ListOrdersRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.OrderClient.ListOrders(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// UpdateOrderStatus handles PUT /api/v1/orders/:id/status
func UpdateOrderStatus(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderID := ctx.Param("id")

		var req UpdateStatusRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &order.UpdateOrderStatusRequest{
			OrderId: orderID,
			Status:  req.Status,
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.OrderClient.UpdateOrderStatus(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
