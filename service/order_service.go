package service

import (
	"context"
	"log/slog"
	"testtask/entity"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o entity.Order) (entity.Order, error)
	OrderWithDelay(ctx context.Context, o entity.Order) error
	DelayedOrders(ctx context.Context, model string, version string) (customers []int64, err error)
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

func (s *OrderService) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	quantity, err := s.robot.GetRobotQuantify(ctx, order.Model, order.Version)
	if err != nil {
		slog.Error("database error: %v", "error", err)
		return entity.Order{}, err
	}

	if quantity == 0 {
		err := s.order.OrderWithDelay(ctx, order)
		if err != nil {
			slog.Error("database error: %v", "error", err)
			return entity.Order{}, err
		}
	}

	order.CreatedAt = time.Now()

	order, err = s.order.CreateOrder(ctx, order)
	if err != nil {
		slog.Error("database error: %v", "error", err)
		return entity.Order{}, err
	}

	return order, err
}
