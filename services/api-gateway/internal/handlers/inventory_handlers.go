package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/test-workflow/api-gateway/internal/grpc/clients"
	"github.com/test-workflow/api-gateway/internal/grpc/proto/inventory"

	"github.com/gin-gonic/gin"
)

// AddProductRequest represents the request body for adding a product
type AddProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Quantity    int32   `json:"quantity" binding:"required,min=0"`
}

// UpdateStockRequest represents the request body for updating stock
type UpdateStockRequest struct {
	QuantityChange int32 `json:"quantity_change" binding:"required"`
}

// ListProducts handles GET /api/v1/inventory/products
func ListProducts(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

		grpcReq := &inventory.ListProductsRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.InventoryClient.ListProducts(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// AddProduct handles POST /api/v1/inventory/products
func AddProduct(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req AddProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &inventory.AddProductRequest{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Quantity:    req.Quantity,
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.InventoryClient.AddProduct(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

// GetStock handles GET /api/v1/inventory/products/:id
func GetStock(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productID := ctx.Param("id")

		grpcReq := &inventory.GetStockRequest{
			ProductId: productID,
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.InventoryClient.GetStock(grpcCtx, grpcReq)
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

// UpdateStock handles PUT /api/v1/inventory/products/:id/stock
func UpdateStock(c *clients.Clients) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productID := ctx.Param("id")

		var req UpdateStockRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &inventory.UpdateStockRequest{
			ProductId:      productID,
			QuantityChange: req.QuantityChange,
		}

		grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.InventoryClient.UpdateStock(grpcCtx, grpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
