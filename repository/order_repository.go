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

func (r *OrderRepository) CreateOrder(ctx context.Context, o entity.Order) (entity.Order, error) {
	q := "INSERT INTO orders(customer_id, robot_model, robot_version, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, o.CustomerID, o.Model, o.Version, o.CreatedAt).Scan(&o.ID)
	if err != nil {
		return entity.Order{}, err
	}

	return o, nil
}

func (r *OrderRepository) OrderWithDelay(ctx context.Context, o entity.Order) error {
	q := "INSERT INTO delayed_order(customer_id, model, version) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, q, o.CustomerID, o.Model, o.Version)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) DelayedOrders(ctx context.Context, model string, version string) (customers []int64, err error) {
	q := "SELECT customer_id FROM delayed_order WHERE model = $1 AND version = $2"

	rows, err := r.db.QueryContext(ctx, q, model, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		customers = append(customers, id)
	}

	return customers, nil
}

func (r *OrderRepository) RemoveDelayedOrder(ctx context.Context, customerID int64, model, version string) error {
	q := "DELETE FROM delayed_order WHERE customer_id = $1 AND model = $2 AND version = $3"

	_, err := r.db.ExecContext(ctx, q, customerID, model, version)
	if err != nil {
		return err
	}

	return nil
}
