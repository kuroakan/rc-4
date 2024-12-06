package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"testtask/entity"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer entity.Customer) (entity.Customer, error)
}

type CustomerHandler struct {
	logger   *slog.Logger
	customer CustomerService
}

func NewCustomerHandler(logger *slog.Logger, customer CustomerService) *CustomerHandler {
	return &CustomerHandler{logger: logger, customer: customer}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	customer := entity.Customer{}
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	err = customer.Validate()
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	customer, err = h.customer.CreateCustomer(ctx, customer)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	sendResponse(w, customer)
}
