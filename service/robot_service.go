package service

import (
	"context"
	"testtask/entity"
)

type RobotRepository interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	GetRobotQuantify(ctx context.Context, model string, version string) (int64, error)
	RobotsCreatedInAWeek(ctx context.Context) (map[string]map[string]int64, error)
}

type RobotService struct {
	robot RobotRepository
}

func NewRobotService(robot RobotRepository) *RobotService {
	return &RobotService{robot: robot}
}

func (s *RobotService) CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error) {
	robot, err := s.robot.CreateRobot(ctx, robot)
	if err != nil {
		return entity.Robot{}, err
	}

	return robot, nil
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
