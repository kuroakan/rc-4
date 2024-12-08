package entity

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

const minNameLen = 4

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Customer) Validate() error {
	if utf8.RuneCountInString(c.Name) < minNameLen {
		return fmt.Errorf("name should be at least %d characters long", minNameLen)
	}

	if c.Email == "" {
		return errors.New("empty email")
	}

	return nil
}
