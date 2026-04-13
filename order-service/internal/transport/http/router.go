package http

import "github.com/gin-gonic/gin"

// NewRouter registers all routes for the Order Service
func NewRouter(h *OrderHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrder)
	r.PATCH("/orders/:id/cancel", h.CancelOrder)

	return r
}
