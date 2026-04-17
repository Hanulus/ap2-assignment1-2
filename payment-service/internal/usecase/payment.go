package usecase

import (
	"errors"

	"github.com/google/uuid"
	"payment-service/internal/domain"
)

// PaymentUseCase contains all payment business logic
type PaymentUseCase struct {
	repo PaymentRepository
}

func NewPaymentUseCase(repo PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

// Authorize processes a payment request.
// Returns the payment record (check Status to see if Authorized or Declined).
func (uc *PaymentUseCase) Authorize(orderID string, amount int64) (*domain.Payment, error) {
	status := domain.StatusAuthorized

	// Business rule: decline payments above the limit
	if amount > domain.MaxAmount {
		status = domain.StatusDeclined
	}

	payment := &domain.Payment{
		ID:            uuid.NewString(),
		OrderID:       orderID,
		TransactionID: uuid.NewString(), // unique transaction ID
		Amount:        amount,
		Status:        status,
	}

	if err := uc.repo.Save(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// GetByOrderID returns the payment for a given order
func (uc *PaymentUseCase) GetByOrderID(orderID string) (*domain.Payment, error) {
	return uc.repo.FindByOrderID(orderID)
}

// ListByAmountRange returns payments filtered by amount range.
// Pass 0 for min or max to skip that limit.
func (uc *PaymentUseCase) ListByAmountRange(min, max int64) ([]*domain.Payment, error) {
	// Validate: if both are set, min must be <= max
	if min > 0 && max > 0 && min > max {
		return nil, errors.New("min_amount cannot be greater than max_amount")
	}
	return uc.repo.FindByAmountRange(min, max)
}
