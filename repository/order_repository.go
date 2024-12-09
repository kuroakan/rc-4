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

func (r *OrderRepository) OrdersByModelVersion(ctx context.Context, model, version string) (orders []entity.Order, err error) {
	q := "SELECT o.id, o.customer_id, o.robot_model, o.robot_version, o.created_at, c.email FROM orders o INNER JOIN customers c ON o.customer_id = c.id WHERE o.robot_model = $1 AND o.robot_version = $2"

	rows, err := r.db.QueryContext(ctx, q, model, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order entity.Order

	for rows.Next() {
		err = rows.Scan(&order.ID, &order.CustomerID, &order.Model, &order.Version, &order.CreatedAt, &order.CustomerEmail)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) Orders(ctx context.Context) (orders []entity.Order, err error) {
	q := "SELECT id, customer_id, robot_model, robot_version, created_at FROM orders"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order entity.Order

	for rows.Next() {
		err = rows.Scan(&order.ID, &order.CustomerID, &order.Model, &order.Version, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) RemoveOrder(ctx context.Context, id int64) error {
	q := "DELETE FROM orders WHERE id = $1"

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}
