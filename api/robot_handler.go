package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"testtask/entity"
	"time"
)

type RobotService interface {
	CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error)
	RobotsCreatedThisWeek(ctx context.Context) (map[string]map[string]int64, error)
}

type RobotHandler struct {
	logger *slog.Logger
	robot  RobotService
}

func NewRobotHandler(logger *slog.Logger, robot RobotService) *RobotHandler {
	return &RobotHandler{logger: logger, robot: robot}
}

type CreateRobotRequest struct {
	Model   string `json:"model"`
	Version string `json:"version"`
	Created string `json:"created"`
}

func (h *RobotHandler) CreateRobot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	var rtc CreateRobotRequest

	err := json.NewDecoder(r.Body).Decode(&rtc)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	rob, err := robotFromAPI(rtc)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	robot, err := h.robot.CreateRobot(ctx, rob)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	sendResponse(w, robot)
}

func robotFromAPI(rtc CreateRobotRequest) (entity.Robot, error) {
	const dateLayout = "2006-01-02 15:04:05"

	t, err := time.Parse(dateLayout, rtc.Created)
	if err != nil {
		return entity.Robot{}, errors.New("invalid creation time")
	}

	robot := entity.Robot{
		Model:   rtc.Model,
		Version: rtc.Version,
		Created: t,
	}

	if rtc.Model == "R2" {
		switch rtc.Version {
		case "D1", "D2", "D3", "D4":
			return robot, nil
		default:
			return entity.Robot{}, errors.New("invalid version")
		}
	}

	if rtc.Model == "13" {
		switch rtc.Version {
		case "X1", "X3", "X4", "X5":
			return robot, nil
		default:
			return entity.Robot{}, errors.New("invalid version")
		}
	}

	return entity.Robot{}, errors.New("invalid model")
}

func (h *RobotHandler) RobotsCreatedThisWeek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	counts, err := h.robot.RobotsCreatedThisWeek(ctx)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	err = handleCSV(w, counts)
	if err != nil {
		sendError(ctx, w, err)
		return
	}
}

func handleCSV(w http.ResponseWriter, counts map[string]map[string]int64) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=weekly report %s.csv", time.Now().Format("2006-01-02")))

	csvWriter := csv.NewWriter(w)

	data := generateCSVData(counts)

	if err := csvWriter.Write([]string{"Model", "Version", "Week Quantity"}); err != nil {
		return err
	}

	for _, record := range data {
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func generateCSVData(counts map[string]map[string]int64) [][]string {
	var data [][]string

	for k, v := range counts {
		for k1, v1 := range v {
			data = append(data, []string{k, k1, strconv.FormatInt(v1, 10)})
		}
	}
	return data
}
