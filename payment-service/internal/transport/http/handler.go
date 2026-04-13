package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

// PaymentHandler handles HTTP requests for payments
type PaymentHandler struct {
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

type authorizeRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}

// Authorize handles POST /payments
func (h *PaymentHandler) Authorize(c *gin.Context) {
	var req authorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.uc.Authorize(req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return 402 if payment was declined (amount too high)
	if payment.Status == domain.StatusDeclined {
		c.JSON(http.StatusPaymentRequired, payment)
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetByOrderID handles GET /payments/:order_id
func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	payment, err := h.uc.GetByOrderID(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}
