package repository

import (
	"context"
	"database/sql"
	"testtask/entity"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, c entity.Customer) (entity.Customer, error) {
	q := "INSERT INTO customers(name, email, created_at) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, c.Name, c.Email, c.CreatedAt).Scan(&c.ID)
	if err != nil {
		return entity.Customer{}, err
	}

	return c, nil
}

func (r *CustomerRepository) CustomerByEmail(ctx context.Context, email string) (customer entity.Customer, err error) {
	q := "SELECT id, name, email, created_at FROM customers WHERE email = $1"

	err = r.db.QueryRowContext(ctx, q, email).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepository) CustomerByID(ctx context.Context, id int) (customer entity.Customer, err error) {
	q := "SELECT id, name, email, created_at FROM customers WHERE id = $1"

	err = r.db.QueryRowContext(ctx, q, id).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}
