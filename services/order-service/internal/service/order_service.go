package service

import (
	"fmt"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/client"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/repository"
)

// OrderService defines the interface for order-related operations.
type OrderService interface {
	CreateOrder(userID uint, request CreateOrderRequest) (*OrderResponse, error)
	GetOrderByID(userID uint, orderID uint) (*OrderResponse, error)
	GetOrders(userID uint) ([]OrderResponse, error)
}

// orderService is the concrete implementation of the OrderService interface.
type orderService struct {
	orderRepository repository.OrderRepository
	productClient   client.ProductClient
}

// NewOrderService creates a new instance of orderService with the provided order repository and product client.
func NewOrderService(orderRepository repository.OrderRepository, productClient client.ProductClient) OrderService {
	return &orderService{
		orderRepository: orderRepository,
		productClient:   productClient,
	}
}

// CreateOrder creates a new order for the specified user with the provided order request.
func (s *orderService) CreateOrder(userID uint, request CreateOrderRequest) (*OrderResponse, error) {
	// Validate the order request
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// Create a new order and calculate the total amount
	order := &model.Order{
		UserID: userID,
		Status: "PLACED",
	}

	var items []model.OrderItem
	var totalAmount float64

	for _, requestItem := range request.Items {

		if requestItem.Quantity <= 0 {
			return nil, fmt.Errorf(
				"quantity must be greater than zero for product %d",
				requestItem.ProductID,
			)
		}

		product, err := s.productClient.GetProduct(requestItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get product %d: %w",
				requestItem.ProductID,
				err,
			)
		}

		if product.StockQuantity < requestItem.Quantity {
			return nil, fmt.Errorf(
				"insufficient stock for product %d",
				requestItem.ProductID,
			)
		}

		subtotal := product.Price * float64(requestItem.Quantity)

		items = append(items, model.OrderItem{
			ProductID: requestItem.ProductID,
			Quantity:  requestItem.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
		})

		totalAmount += subtotal
	}

	order.TotalAmount = totalAmount

	if err := s.orderRepository.CreateWithTransaction(order, items); err != nil {
		return nil, err
	}

	order.Items = items

	return toOrderResponse(order), nil
}

// GetOrderByID retrieves an order by its ID for the specified user.
func (s *orderService) GetOrderByID(userID uint, orderID uint) (*OrderResponse, error) {

	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	return toOrderResponse(order), nil
}

// GetOrders retrieves all orders for the specified user.
func (s *orderService) GetOrders(userID uint) ([]OrderResponse, error) {
	orders, err := s.orderRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]OrderResponse, 0, len(orders))

	for _, order := range orders {
		responses = append(
			responses,
			*toOrderResponse(&order),
		)
	}

	return responses, nil
}

// toOrderResponse converts an Order model to an OrderResponse DTO.
func toOrderResponse(order *model.Order) *OrderResponse {

	items := make([]OrderItemResponse, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	return &OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       items,
	}
}
