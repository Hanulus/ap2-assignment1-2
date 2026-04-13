package domain

import "time"

// Order statuses
const (
	StatusPending   = "Pending"
	StatusPaid      = "Paid"
	StatusFailed    = "Failed"
	StatusCancelled = "Cancelled"
)

// Order is the core entity of the Order Service
type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64 // in cents, e.g. 1000 = $10.00
	Status     string
	CreatedAt  time.Time
}

// CanBeCancelled checks business rule: only Pending orders can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == StatusPending
}
