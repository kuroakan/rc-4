package entity

import (
	"time"
)

type Robot struct {
	ID      int64     `json:"id"`
	Model   string    `json:"model"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
}
