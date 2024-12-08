package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testtask/entity"
)

type RobotRepository interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	GetRobotQuantify(ctx context.Context, model string, version string) (int64, error)
	RobotsCreatedInAWeek(ctx context.Context) (map[string]map[string]int64, error)
}

type Sender interface {
	SendMail(email, text, subject string) error
}

type Orderer interface {
	Orders(ctx context.Context, model, version string) ([]entity.Order, error)
	RemoveOrder(ctx context.Context, id int64) error
}

type RobotService struct {
	robot  RobotRepository
	order  Orderer
	sender Sender
}

func NewRobotService(robot RobotRepository, order Orderer, sender Sender) *RobotService {
	return &RobotService{robot: robot, order: order, sender: sender}
}

func (s *RobotService) CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error) {
	robot, err := s.robot.CreateRobot(ctx, robot)
	if err != nil {
		return entity.Robot{}, err
	}

	err = s.sendMessageAboutRobotCreation(ctx, robot.Model, robot.Version)
	if err != nil {
		slog.Error("error sending message to customer", "error", err)
	}

	return robot, nil
}

func (s *RobotService) sendMessageAboutRobotCreation(ctx context.Context, model string, version string) error {
	orders, err := s.order.Orders(ctx, model, version)
	if err != nil {
		return fmt.Errorf("get orders: %w", err)
	}

	text := fmt.Sprintf(`Good afternoon!
You recently inquired about our Model %s, Version %s robot.
This robot is now in stock. If this option is suitable for you, please contact us`, model, version)

	subject := "Update about your order!"

	var errs []error

	for _, order := range orders {
		err := s.sender.SendMail(order.CustomerEmail, text, subject)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = s.order.RemoveOrder(ctx, order.ID)
		if err != nil {
			return fmt.Errorf("remove order %d: %w", order.ID, err)
		}
	}

	return errors.Join(errs...)
}

func (s *RobotService) RobotsCreatedThisWeek(ctx context.Context) (map[string]map[string]int64, error) {
	counts, err := s.robot.RobotsCreatedInAWeek(ctx)
	if err != nil {
		return nil, err
	}

	return counts, nil
}

func (s *RobotService) GetRobotQuantity(ctx context.Context, model, version string) (int64, error) {
	return s.robot.GetRobotQuantify(ctx, model, version)
}
