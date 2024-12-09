package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

func (n *Notifier) Notify() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errs []error

	subject := "Update about your order!"
	var text string

	orders, err := n.order.Orders(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("notifier get order: %s", err), "error", err)
		return
	}

	for _, order := range orders {
		quantity, err := n.robot.GetRobotQuantity(ctx, order.Model, order.Version)
		if err != nil {
			slog.Error(fmt.Sprintf("notifier get quantity: %s", err), "error", err)
			return
		}

		if quantity <= 1 {
			continue
		}

		text = fmt.Sprintf(`Good afternoon!
You recently inquired about our Model %s, Version %s robot.
This robot is now in stock. If this option is suitable for you, please contact us`, order.Model, order.Version)

		err = n.sender.SendMail(order.CustomerEmail, text, subject)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = n.order.RemoveOrder(ctx, order.ID)
		if err != nil {
			slog.Error(fmt.Sprintf("notifier remove order: %s", err), "error", err)
			return
		}
	}

	if len(errs) > 0 {
		slog.Error("notifier", "error", errors.Join(errs...))
	}
}
