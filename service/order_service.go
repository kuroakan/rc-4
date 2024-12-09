package service

import (
	"context"
	"log/slog"
	"testtask/entity"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	OrdersByModelVersion(ctx context.Context, model, version string) (orders []entity.Order, err error)
	Orders(ctx context.Context) ([]entity.Order, error)
	RemoveOrder(ctx context.Context, id int64) error
}

type RobotServ interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	RobotsCreatedThisWeek(ctx context.Context) (map[string]map[string]int64, error)
	GetRobotQuantity(ctx context.Context, model, version string) (int64, error)
}

type OrderService struct {
	order OrderRepository
	robot RobotServ
}

func NewOrderService(order OrderRepository, robot RobotServ) *OrderService {
	return &OrderService{
		order: order,
		robot: robot,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order entity.Order) (int64, error) {
	quantity, err := s.robot.GetRobotQuantity(ctx, order.Model, order.Version)
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

func (s *OrderService) OrdersByModelVersion(ctx context.Context, model, version string) ([]entity.Order, error) {
	orders, err := s.order.OrdersByModelVersion(ctx, model, version)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) RemoveOrder(ctx context.Context, id int64) error {
	return s.order.RemoveOrder(ctx, id)
}

func (s *OrderService) Orders(ctx context.Context) ([]entity.Order, error) {
	return s.order.Orders(ctx)
}
