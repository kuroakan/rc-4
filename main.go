package main

import (
	"log/slog"
	"net/http"
	"os"
	"testtask/api"
	"testtask/bootstrap"
	"testtask/repository"
	"testtask/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := bootstrap.NewConfig()
	if err != nil {
		logger.Error("Problem with config load: ", "error", err)
		return
	}

	errorList := cfg.Validate()
	if errorList != nil {
		logger.Error("Config validation error ", "error", errorList)
		return
	}

	db, err := bootstrap.DBConnect(cfg)
	if err != nil {
		logger.Error("Problem with Postgres connection", "error", err)
		return
	}
	defer db.Close()

	logger.Info("postgres DB connection status: OK")

	orderRepo := repository.NewOrderRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	robotRepo := repository.NewRobotRepository(db)

	orderService := service.NewOrderService(orderRepo, robotRepo)
	customerService := service.NewCustomerService(customerRepo, robotRepo, orderRepo)
	robotService := service.NewRobotService(robotRepo, orderRepo)

	orderHandler := api.NewOrderHandler(logger, orderService)
	customerHandler := api.NewCustomerHandler(logger, customerService)
	robotHandler := api.NewRobotHandler(logger, robotService)

	server := api.NewServer(cfg.HTTPPort, http.DefaultServeMux, customerHandler, orderHandler, robotHandler)

	err = server.Start()
	if err != nil {
		logger.Error("Server start error", "error", err)
		return
	}
}
