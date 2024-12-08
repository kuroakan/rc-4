package entity

import (
	"errors"
	"time"
)

const dateTemplate = "2006-01-02 15:04:05"

type Robot struct {
	ID      int64     `json:"id"`
	Model   string    `json:"model"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
}

type RobotToCreate struct {
	Model   string `json:"model"`
	Version string `json:"version"`
	Created string `json:"created"`
}

func RobotCreateAPI(rtc RobotToCreate) (Robot, error) {
	t, err := time.Parse(dateTemplate, rtc.Created)
	if err != nil {
		return Robot{}, errors.New("invalid creation time")
	}

	robot := Robot{
		Model:   rtc.Model,
		Version: rtc.Version,
		Created: t,
	}

	if rtc.Model == "R2" {
		switch rtc.Version {
		case "D1", "D2", "D3", "D4":
			return robot, nil
		default:
			return Robot{}, errors.New("invalid version")
		}
	}

	if rtc.Model == "13" {
		switch rtc.Version {
		case "X1", "X3", "X4", "X5":
			return robot, nil
		default:
			return Robot{}, errors.New("invalid version")
		}
	}

	return Robot{}, errors.New("invalid model")
}
