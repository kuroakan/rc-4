package entity

import (
	"errors"
	"time"
)

type Order struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	Model      string    `json:"model"`
	Version    string    `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
}

func (o *Order) Validate() error {
	if o.CustomerID == 0 {
		err := errors.New("invalid customer id")
		return err
	}

	if o.Model == "" {
		err := errors.New("empty model field")
		return err
	}

	if o.Version == "" {
		err := errors.New("empty version field")
		return err
	}

	return nil
}
