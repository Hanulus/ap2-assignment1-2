package usecase

import "payment-service/internal/domain"

// PaymentRepository is the port for persisting payments
type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByOrderID(orderID string) (*domain.Payment, error)
	FindByAmountRange(min, max int64) ([]*domain.Payment, error) // new
}
