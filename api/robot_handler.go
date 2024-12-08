package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"testtask/entity"
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

func (h *RobotHandler) CreateRobot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	rtc := entity.RobotToCreate{}

	err := json.NewDecoder(r.Body).Decode(&rtc)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	rob, err := entity.RobotCreateAPI(rtc)
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

func (h *RobotHandler) RobotsCreatedThisWeek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", h.logger)

	counts, err := h.robot.RobotsCreatedThisWeek(ctx)
	if err != nil {
		sendError(ctx, w, err)
		return
	}

	handleCSV(ctx, w, counts)
}

func handleCSV(ctx context.Context, w http.ResponseWriter, counts map[string]map[string]int64) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=weekly report.csv")

	csvWriter := csv.NewWriter(w)

	data := generateCSVData(counts)
	for _, record := range data {
		if err := csvWriter.Write(record); err != nil {
			sendError(ctx, w, err)
			return
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		sendError(ctx, w, err)
		return
	}
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
