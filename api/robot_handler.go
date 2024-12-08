package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

	h.sendScvFile(counts, w, r)
}

func (h *RobotHandler) sendScvFile(counts map[string]map[string]int64, w http.ResponseWriter, r *http.Request) {
	filename := "output.csv"
	writer, file, err := createCSVWriter(filename)
	if err != nil {
		fmt.Println("Error creating CSV writer:", err)
		return
	}
	defer file.Close()

	headers := []string{"Model", "Version", "Week quantity"}
	writeCSVRecord(writer, headers)

	var data [][]string

	for k, v := range counts {
		for k1, v1 := range v {
			data = append(data, []string{k, k1, strconv.FormatInt(v1, 10)})
		}
	}

	for _, record := range data {
		writeCSVRecord(writer, record)
	}
	// Flush the writer and check for any errors
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"data.csv\"")

	file.Seek(0, 0)
	http.ServeFile(w, r, file.Name())
}

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing record to CSV:", err)
	}
}
