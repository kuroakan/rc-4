package entity

import (
	"errors"
	"time"
)

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Customer) Validate() error {
	if c.Name == "" || len(c.Name) < 4 {
		return errors.New("name is too short")
	}

	if c.Email == "" {
		return errors.New("email already taken")
	}

	return nil
}
