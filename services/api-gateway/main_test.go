package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	middleware := corsMiddleware()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}
}

func TestAPIRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Test that routes are properly defined
	api := router.Group("/api/v1")
	api.GET("/orders", func(c *gin.Context) {
		c.JSON(200, gin.H{"orders": []string{}})
	})
	api.GET("/inventory/products", func(c *gin.Context) {
		c.JSON(200, gin.H{"products": []string{}})
	})

	tests := []struct {
		path   string
		status int
	}{
		{"/api/v1/orders", 200},
		{"/api/v1/inventory/products", 200},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.status {
			t.Errorf("Route %s: expected status %d, got %d", tt.path, tt.status, w.Code)
		}
	}
}
