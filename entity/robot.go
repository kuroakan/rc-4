package entity

import (
	"errors"
	"time"
)

type Robot struct {
	ID        int64     `json:"id"`
	Model     string    `json:"model"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Robot) Validate() error {
	if r.Model == "R2" {
		switch r.Version {
		case "D1":
			return nil
		case "D2":
			return nil
		case "D3":
			return nil
		case "D4":
			return nil
		default:
			return errors.New("invalid version")
		}
	}

	if r.Model == "13" {
		switch r.Version {
		case "X1":
			return nil
		case "X3":
			return nil
		case "X4":
			return nil
		case "X5":
			return nil
		default:
			return errors.New("invalid version")
		}
	}

	return errors.New("invalid model")
}
