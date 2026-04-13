package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"order-service/internal/usecase"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	uc *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

type createOrderRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	ItemName   string `json:"item_name" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.uc.CreateOrder(req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		// Use case validation error (e.g. amount <= 0)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If payment failed/service unavailable, return 503
	if order.Status == "Failed" {
		c.JSON(http.StatusServiceUnavailable, order)
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.uc.GetOrder(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// CancelOrder handles PATCH /orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	order, err := h.uc.CancelOrder(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
