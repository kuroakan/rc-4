package service

import (
	"context"
	"log/slog"
	"testtask/entity"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	Orders(ctx context.Context, model, version string) (orders []entity.Order, err error)
	RemoveOrder(ctx context.Context, id int64) error
}

type Roboter interface {
	GetRobotQuantity(ctx context.Context, model, version string) (int64, error)
}

type OrderService struct {
	order OrderRepository
	robot RobotRepository
}

func NewOrderService(order OrderRepository, robot RobotRepository) *OrderService {
	return &OrderService{
		order: order,
		robot: robot,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order entity.Order) (int64, error) {
	quantity, err := s.robot.GetRobotQuantify(ctx, order.Model, order.Version)
	if err != nil {
		slog.Error("database error: %v", "error", err)
		return 0, err
	}

	if quantity == 0 {
		order.CreatedAt = time.Now()

		err = s.order.CreateOrder(ctx, order)
		if err != nil {
			slog.Error("database error: %v", "error", err)
			return 0, err
		}
	}

	return quantity, err
}

func (s *OrderService) Orders(ctx context.Context, model, version string) ([]entity.Order, error) {
	orders, err := s.order.Orders(ctx, model, version)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) RemoveOrder(ctx context.Context, id int64) error {
	return s.order.RemoveOrder(ctx, id)
}
