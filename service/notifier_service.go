package service

import (
	"context"
	"errors"
	"fmt"
	"testtask/entity"
)

type OrderSv interface {
	Orders(ctx context.Context) ([]entity.Order, error)
	RemoveOrder(ctx context.Context, id int64) error
}

type SenderSv interface {
	SendMail(email, text, subject string) error
}

type RobotSv interface {
	GetRobotQuantity(ctx context.Context, model, version string) (int64, error)
}

type Notifier struct {
	order  OrderSv
	sender SenderSv
	robot  RobotSv
}

func NewNotifier(order OrderSv, sender SenderSv, robot RobotSv) *Notifier {
	return &Notifier{order: order, sender: sender, robot: robot}
}

func (n *Notifier) NotifyCustomers() error {
	ctx := context.Background() //TODO gracefull shutdown

	orders, err := n.order.Orders(ctx)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	var errs []error

	for _, order := range orders {
		quantity, err := n.robot.GetRobotQuantity(ctx, order.Model, order.Version)
		if err != nil {
			return fmt.Errorf("get quantity: %w", err)
		}

		if quantity == 0 {
			continue
		}

		text := fmt.Sprintf(`Good afternoon!
You recently inquired about our Model %s, Version %s robot.
This robot is now in stock. If this option is suitable for you, please contact us`, order.Model, order.Version)

		err = n.sender.SendMail(order.CustomerEmail, text, "Update about your order!")
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = n.order.RemoveOrder(ctx, order.ID)
		if err != nil {
			return fmt.Errorf("remove order: %w", err)
		}
	}
	return errors.Join(errs...)
}
