package http

import "github.com/gin-gonic/gin"

// NewRouter registers all routes for the Payment Service
func NewRouter(h *PaymentHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/payments", h.Authorize)
	r.GET("/payments/:order_id", h.GetByOrderID)

	return r
}
