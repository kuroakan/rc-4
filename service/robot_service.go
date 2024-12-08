package service

import (
	"context"
	"fmt"
	"github.com/wneessen/go-mail"
	"log/slog"
	"os"
	"testtask/entity"
)

type RobotRepository interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	GetRobotQuantify(ctx context.Context, model string, version string) (int64, error)
	RobotsCreatedInAWeek(ctx context.Context) (map[string]map[string]int64, error)
}

type RobotService struct {
	robot    RobotRepository
	order    OrderRepository
	customer CustomerRepository
}

func NewRobotService(robot RobotRepository, order OrderRepository, customer CustomerRepository) *RobotService {
	return &RobotService{robot: robot, order: order, customer: customer}
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
	ids, err := s.order.DeleteOrderWithNotification(ctx, model, version)
	if err != nil {
		return err
	}

	for _, id := range ids {
		customer, err := s.customer.CustomerByID(ctx, id)
		if err != nil {
			return err
		}

		err = s.sendMail(customer.Email)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RobotService) RobotsCreatedThisWeek(ctx context.Context) (map[string]map[string]int64, error) {
	counts, err := s.robot.RobotsCreatedInAWeek(ctx)
	if err != nil {
		return nil, err
	}

	return counts, nil
}

func (s *RobotService) sendMail(email string) error {
	message := mail.NewMsg()
	if err := message.From(os.Getenv("MAIL_ADDRESS")); err != nil {
		return fmt.Errorf("failed to send From address: %s", err)
	}

	if err := message.To(email); err != nil {
		return fmt.Errorf("failed to send To address: %s", err)
	}

	message.Subject("Update about your order!")

	message.SetBodyString(mail.TypeTextPlain, "Good afternoon!\n"+
		"You recently inquired about our Model X, Version Y robot. \n"+
		"This robot is now in stock. If this option is suitable for you, please contact us")
	client, err := mail.NewClient(os.Getenv("MAIL_HOST"), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("MAIL_USERNAME")), mail.WithPassword(os.Getenv("MAIL_PASSWORD")))
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}

	return nil
}
