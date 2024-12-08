package repository

import (
	"context"
	"database/sql"
	"testtask/entity"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order entity.Order) error {
	q := "INSERT INTO orders(customer_id, robot_model, robot_version, created_at) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, q, order.CustomerID, order.Model, order.Version, order.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) DeleteOrderWithNotification(ctx context.Context, model, version string) ([]int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := "SELECT customer_id FROM orders WHERE robot_model = $1 AND robot_version = $2"

	rows, err := tx.QueryContext(ctx, q, model, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var result []int64

	for rows.Next() {
		rows.Scan(&id)

		result = append(result, id)
	}

	q = "DELETE FROM orders WHERE robot_model = $1 AND robot_version = $2"

	_, err = tx.ExecContext(ctx, q, model, version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}
