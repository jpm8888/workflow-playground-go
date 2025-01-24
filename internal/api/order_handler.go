package api

import (
	"app/myproj/internal/order"
	"app/myproj/pkg/workflow"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandler struct {
	workflowEngine *workflow.WorkflowEngine
	orderService   *order.OrderService
}

func NewOrderHandler(engine *workflow.WorkflowEngine, orderService *order.OrderService) *OrderHandler {
	return &OrderHandler{
		workflowEngine: engine,
		orderService:   orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input struct {
		UserID    uint `json:"user_id" binding:"required"`
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(input.UserID, input.ProductID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ConfirmPayment(c *gin.Context) {
	orderID := c.Param("id")
	// Implement payment confirmation logic here
}

func (h *OrderHandler) FulfillOrder(c *gin.Context) {
	orderID := c.Param("id")
	// Implement order fulfillment logic here
}
