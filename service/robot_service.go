package service

import (
	"context"
	"log/slog"
	"testtask/entity"
	"time"
)

type RobotRepository interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	GetRobotQuantify(ctx context.Context, model string, version string) (int64, error)
	RobotsCreatedInAWeek(ctx context.Context) (map[string]map[string]int64, error)
}

type RobotService struct {
	robot RobotRepository
	order OrderRepository
}

func NewRobotService(robot RobotRepository, order OrderRepository) *RobotService {
	return &RobotService{robot: robot, order: order}
}

func (s *RobotService) CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error) {
	err := robot.Validate()
	if err != nil {
		return entity.Robot{}, err
	}

	robot.CreatedAt = time.Now().UTC()

	robot, err = s.robot.CreateRobot(ctx, robot)
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
	return nil
}

func (s *RobotService) RobotsCreatedThisWeek(ctx context.Context) (map[string]map[string]int64, error) {
	counts, err := s.robot.RobotsCreatedInAWeek(ctx)
	if err != nil {
		return nil, err
	}

	return counts, nil
}
