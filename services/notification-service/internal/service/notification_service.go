package service

import (
	"fmt"

	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/event"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) HandleOrderCreated(
	orderEvent event.OrderCreatedEvent,
) error {

	// TODO : Implement actual notification logic (e.g., sending email, SMS, push notification, etc.)
	fmt.Printf(
		"NOTIFICATION: Order #%d created successfully for User #%d. Total amount: %.2f\n",
		orderEvent.OrderID,
		orderEvent.UserID,
		orderEvent.TotalAmount,
	)

	return nil
}
