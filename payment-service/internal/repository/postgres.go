package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"payment-service/internal/domain"
)

// PostgresPaymentRepo implements PaymentRepository using PostgreSQL
type PostgresPaymentRepo struct {
	db *sql.DB
}

func NewPostgresPaymentRepo(db *sql.DB) *PostgresPaymentRepo {
	return &PostgresPaymentRepo{db: db}
}

// Save stores a new payment record
func (r *PostgresPaymentRepo) Save(payment *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, payment.ID, payment.OrderID, payment.TransactionID, payment.Amount, payment.Status)
	return err
}

// FindByOrderID looks up a payment by the order it belongs to
func (r *PostgresPaymentRepo) FindByOrderID(orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id = $1`
	row := r.db.QueryRow(query, orderID)

	var p domain.Payment
	err := row.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("payment not found")
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByAmountRange returns payments filtered by amount range.
// min=0 means no lower bound, max=0 means no upper bound.
func (r *PostgresPaymentRepo) FindByAmountRange(min, max int64) ([]*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if min > 0 {
		query += fmt.Sprintf(" AND amount >= $%d", argIdx)
		args = append(args, min)
		argIdx++
	}
	if max > 0 {
		query += fmt.Sprintf(" AND amount <= $%d", argIdx)
		args = append(args, max)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status); err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}
	return payments, nil
}
