package main

import (
	"fmt"
	"github.com/wneessen/go-mail"
	"log/slog"
	"os"
	"testtask/api"
	"testtask/bootstrap"
	"testtask/repository"
	"testtask/service"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := bootstrap.NewConfig()
	if err != nil {
		logger.Error("Problem with config load: ", "error", err)
		return
	}

	db, err := bootstrap.DBConnect(cfg)
	if err != nil {
		logger.Error("Problem with Postgres connection", "error", err)
		return
	}
	defer db.Close()

	logger.Info("postgres DB connection status: OK")

	client, err := mail.NewClient(cfg.Mail.Host, mail.WithPort(cfg.Mail.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Mail.Username), mail.WithPassword(cfg.Mail.Password))
	if err != nil {
		logger.Error("Problem with mail client connection", "error", err)
		return
	}

	orderRepo := repository.NewOrderRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	robotRepo := repository.NewRobotRepository(db)

	mailService := service.NewEmailService(client, cfg.Mail.From)
	robotService := service.NewRobotService(robotRepo)
	orderService := service.NewOrderService(orderRepo, robotService)
	customerService := service.NewCustomerService(customerRepo)
	notifierService := service.NewNotifier(orderService, mailService, robotService)

	orderHandler := api.NewOrderHandler(logger, orderService)
	customerHandler := api.NewCustomerHandler(logger, customerService)
	robotHandler := api.NewRobotHandler(logger, robotService)

	go func() {
		for {
			slog.Info("notifier started")
			notifierService.Notify()
			time.Sleep(time.Hour)
		}
	}()

	server := api.NewServer(cfg.HTTPPort, customerHandler, orderHandler, robotHandler)

	logger.Info(fmt.Sprintf("Server is listening at port: %s", cfg.HTTPPort))

	err = server.Start()
	if err != nil {
		logger.Error("Server start error", "error", err)
		return
	}
}
