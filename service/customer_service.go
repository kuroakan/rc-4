package service

import (
	"context"
	"errors"
	"fmt"
	"testtask/entity"
	"time"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, c entity.Customer) (entity.Customer, error)
	CustomerByEmail(ctx context.Context, email string) (customer entity.Customer, err error)
	CustomerByID(ctx context.Context, id int64) (customer entity.Customer, err error)
}

type CustomerService struct {
	customer CustomerRepository
	robot    RobotRepository
	order    OrderRepository
}

func NewCustomerService(customer CustomerRepository, robot RobotRepository, order OrderRepository) *CustomerService {
	return &CustomerService{
		customer: customer,
		robot:    robot,
		order:    order,
	}
}

func (c *CustomerService) CreateCustomer(ctx context.Context, customer entity.Customer) (entity.Customer, error) {
	_, err := c.customer.CustomerByEmail(ctx, customer.Email)
	if err == nil {
		return entity.Customer{}, fmt.Errorf("customer with that email: %w", entity.ErrAlreadyExist)
	}

	if !errors.Is(err, entity.ErrNotFound) {
		return entity.Customer{}, err
	}

	customer.CreatedAt = time.Now()

	customer, err = c.customer.CreateCustomer(ctx, customer)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}
