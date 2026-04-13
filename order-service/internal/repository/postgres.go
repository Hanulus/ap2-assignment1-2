package repository

import (
	"database/sql"
	"errors"
	"time"

	"order-service/internal/domain"
)

// PostgresOrderRepo implements OrderRepository using PostgreSQL
type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

// Save inserts a new order into the database
func (r *PostgresOrderRepo) Save(order *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, order.ID, order.CustomerID, order.ItemName, order.Amount, order.Status, order.CreatedAt)
	return err
}

// FindByID looks up an order by its ID
func (r *PostgresOrderRepo) FindByID(id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var order domain.Order
	var createdAt time.Time
	err := row.Scan(&order.ID, &order.CustomerID, &order.ItemName, &order.Amount, &order.Status, &createdAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}
	order.CreatedAt = createdAt
	return &order, nil
}

// UpdateStatus changes the order status
func (r *PostgresOrderRepo) UpdateStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}
