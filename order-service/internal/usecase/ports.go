package usecase

import "order-service/internal/domain"

// OrderRepository is the port that the use case depends on.
// The actual implementation lives in the repository layer.
type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status string) error
}

// PaymentClient is the port for calling the Payment Service.
type PaymentClient interface {
	Authorize(orderID string, amount int64) (string, error) // returns transactionID
}
