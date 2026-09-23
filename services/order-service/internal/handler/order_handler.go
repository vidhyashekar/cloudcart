package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Create handles the creation of a new order. It expects a JSON payload containing the order details and the user ID from the context.
func (h *OrderHandler) Create(c *gin.Context) {
	var request service.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	order, err := h.orderService.CreateOrder(
		userID.(uint),
		request,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetByID handles the retrieval of an order by its ID. It expects the order ID as a URL parameter and the user ID from the context.
func (h *OrderHandler) GetByID(c *gin.Context) {
	orderID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	order, err := h.orderService.GetOrderByID(
		userID.(uint),
		uint(orderID),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetAll handles the retrieval of all orders for the authenticated user. It expects the user ID from the context.
func (h *OrderHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	orders, err := h.orderService.GetOrders(userID.(uint))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}
