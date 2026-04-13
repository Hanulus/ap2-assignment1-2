package domain

// Payment statuses
const (
	StatusAuthorized = "Authorized"
	StatusDeclined   = "Declined"

	// MaxAmount is the business limit: amounts above this are declined
	MaxAmount int64 = 100000 // 1000 units in cents
)

// Payment is the core entity of the Payment Service
type Payment struct {
	ID            string
	OrderID       string
	TransactionID string
	Amount        int64 // in cents
	Status        string
}
