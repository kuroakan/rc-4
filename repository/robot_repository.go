package repository

import (
	"context"
	"database/sql"
	"testtask/entity"
)

type RobotRepository struct {
	db *sql.DB
}

func NewRobotRepository(db *sql.DB) *RobotRepository {
	return &RobotRepository{db: db}
}

func (r *RobotRepository) CreateRobot(ctx context.Context, robot entity.Robot) (entity.Robot, error) {
	q := "INSERT INTO robots(model, version, created_at) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, q, robot.Model, robot.Version, robot.CreatedAt).Scan(&robot.ID)
	if err != nil {
		return entity.Robot{}, err
	}

	return robot, nil
}

func (r *RobotRepository) GetRobotQuantify(ctx context.Context, model string, version string) (quantity int64, err error) {
	q := "SELECT count(*) from robots WHERE model = $1 AND version = $2"

	err = r.db.QueryRowContext(ctx, q, model, version).Scan(&quantity)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (r *RobotRepository) RobotsCreatedInAWeek(ctx context.Context) (map[string]map[string]int64, error) {
	q := "SELECT model, version from robots WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]map[string]int64)

	for rows.Next() {
		var model, version string
		if err := rows.Scan(&model, &version); err != nil {
			return nil, err
		}

		if _, exists := counts[model]; !exists {
			counts[model] = make(map[string]int64)
		}

		counts[model][version]++
	}

	return counts, nil
}
