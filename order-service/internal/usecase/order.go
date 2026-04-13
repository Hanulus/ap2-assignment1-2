package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"order-service/internal/domain"
)

// OrderUseCase contains all business logic for orders
type OrderUseCase struct {
	repo    OrderRepository
	payment PaymentClient
}

func NewOrderUseCase(repo OrderRepository, payment PaymentClient) *OrderUseCase {
	return &OrderUseCase{repo: repo, payment: payment}
}

// CreateOrder creates a Pending order, calls Payment Service, then updates status
func (uc *OrderUseCase) CreateOrder(customerID, itemName string, amount int64) (*domain.Order, error) {
	// Business rule: amount must be positive
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	order := &domain.Order{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	// Save order as Pending first
	if err := uc.repo.Save(order); err != nil {
		return nil, err
	}

	// Call Payment Service to authorize the payment
	_, err := uc.payment.Authorize(order.ID, order.Amount)
	if err != nil {
		// Payment failed or service unavailable — mark as Failed
		_ = uc.repo.UpdateStatus(order.ID, domain.StatusFailed)
		order.Status = domain.StatusFailed
		return order, nil
	}

	// Payment authorized — mark as Paid
	_ = uc.repo.UpdateStatus(order.ID, domain.StatusPaid)
	order.Status = domain.StatusPaid
	return order, nil
}

// GetOrder fetches order by ID
func (uc *OrderUseCase) GetOrder(id string) (*domain.Order, error) {
	return uc.repo.FindByID(id)
}

// CancelOrder cancels a Pending order; Paid orders cannot be cancelled
func (uc *OrderUseCase) CancelOrder(id string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Business rule: only Pending orders can be cancelled
	if !order.CanBeCancelled() {
		return nil, errors.New("only pending orders can be cancelled")
	}

	if err := uc.repo.UpdateStatus(id, domain.StatusCancelled); err != nil {
		return nil, err
	}
	order.Status = domain.StatusCancelled
	return order, nil
}
