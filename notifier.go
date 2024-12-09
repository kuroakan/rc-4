package main

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

func notifier(sender SenderSv, ors OrderSv, rs RobotSv) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errs []error

	subject := "Update about your order!"
	var text string

	orders, err := ors.Orders(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("notifier: %s", err), "error", err)
		return
	}

	for _, order := range orders {
		quantity, err := rs.GetRobotQuantity(ctx, order.Model, order.Version)
		if err != nil {
			slog.Error(fmt.Sprintf("notifier: %s", err), "error", err)
			return
		}

		if quantity <= 1 {
			continue
		}

		text = fmt.Sprintf(`Good afternoon!
You recently inquired about our Model %s, Version %s robot.
This robot is now in stock. If this option is suitable for you, please contact us`, order.Model, order.Version)

		err = sender.SendMail(order.CustomerEmail, text, subject)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = ors.RemoveOrder(ctx, order.ID)
		if err != nil {
			slog.Error(fmt.Sprintf("notifier: %s", err), "error", err)
			return
		}
	}
	if len(errs) > 0 {
		slog.Error("notifier", "error", errors.Join(errs...))
	}
}
