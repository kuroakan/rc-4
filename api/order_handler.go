package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"testtask/entity"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error)
}

type OrderHandler struct {
	logger *slog.Logger
	order  OrderService
}

func NewOrderHandler(logger *slog.Logger, order OrderService) *OrderHandler {
	return &OrderHandler{logger: logger, order: order}
}

func (h *OrderHandler) OrderRobot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	order := entity.Order{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	order, err = h.order.CreateOrder(ctx, order)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	sendResponse(w, order)
}
