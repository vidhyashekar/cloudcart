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

// CreateOrder creates a new order for the specified user based on the provided request.
func (s *orderService) CreateOrder(userID uint, request CreateOrderRequest) (*OrderResponse, error) {

	// 1. Validate order
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// 2. Prepare order
	order := &model.Order{
		UserID: userID,
		Status: "PLACED",
	}

	var items []model.OrderItem
	var totalAmount float64

	// Keep track of products whose stock was successfully decreased.
	var decreasedProducts []CreateOrderItemRequest

	// 3. Validate products and calculate total
	for _, requestItem := range request.Items {

		if requestItem.Quantity <= 0 {
			return nil, fmt.Errorf(
				"quantity must be greater than zero for product %d",
				requestItem.ProductID,
			)
		}

		// Get product details from Product Service.
		product, err := s.productClient.GetProduct(
			requestItem.ProductID,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to get product %d: %w",
				requestItem.ProductID,
				err,
			)
		}

		// Check stock.
		if product.StockQuantity < requestItem.Quantity {
			return nil, fmt.Errorf(
				"insufficient stock for product %d",
				requestItem.ProductID,
			)
		}

		// Calculate item subtotal using the price
		// returned by Product Service.
		subtotal := product.Price * float64(requestItem.Quantity)

		// Create order item.
		items = append(items, model.OrderItem{
			ProductID: requestItem.ProductID,
			Quantity:  requestItem.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
		})

		// Calculate order total.
		totalAmount += subtotal
	}

	order.TotalAmount = totalAmount

	// 4. Reserve/decrease stock.
	for _, item := range request.Items {

		if err := s.productClient.DecreaseStock(
			item.ProductID,
			item.Quantity,
		); err != nil {

			// Compensate stock that was already decreased.
			for _, decreased := range decreasedProducts {
				_ = s.productClient.IncreaseStock(
					decreased.ProductID,
					decreased.Quantity,
				)
			}

			return nil, fmt.Errorf(
				"failed to reserve stock for product %d: %w",
				item.ProductID,
				err,
			)
		}

		// Remember successful stock updates.
		decreasedProducts = append(
			decreasedProducts,
			item,
		)
	}

	// 5. Create Order + Order Items in one DB transaction.
	if err := s.orderRepository.CreateWithTransaction(
		order,
		items,
	); err != nil {
		// Order creation failed.
		// Restore all stock that was already decreased.
		for _, item := range decreasedProducts {
			_ = s.productClient.IncreaseStock(
				item.ProductID,
				item.Quantity,
			)
		}

		return nil, fmt.Errorf(
			"failed to create order: %w",
			err,
		)
	}

	// 6. Attach items to order for response.
	order.Items = items

	// 7. Return response.
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
