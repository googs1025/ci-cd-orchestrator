package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// MetricRepository 指标仓库
type MetricRepository struct {
	db *sql.DB
}

// NewMetricRepository 创建指标仓库实例
func NewMetricRepository(db *sql.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// Create 创建指标
func (r *MetricRepository) Create(metric *models.Metric) error {
	query := `
		INSERT INTO metrics (execution_id, name, value, created_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, metric.ExecutionID, metric.Name, metric.Value, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	metric.ID = int(id)
	metric.CreatedAt = now

	return nil
}

// GetByExecutionID 根据执行 ID 获取指标
func (r *MetricRepository) GetByExecutionID(executionID int) ([]*models.Metric, error) {
	query := `
		SELECT id, execution_id, name, value, created_at
		FROM metrics
		WHERE execution_id = ?
	`

	rows, err := r.db.Query(query, executionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*models.Metric
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(
			&metric.ID,
			&metric.ExecutionID,
			&metric.Name,
			&metric.Value,
			&metric.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// GetByProjectID 根据项目 ID 获取指标（通过执行历史关联）
func (r *MetricRepository) GetByProjectID(projectID int) ([]*models.Metric, error) {
	query := `
		SELECT m.id, m.execution_id, m.name, m.value, m.created_at
		FROM metrics m
		JOIN executions e ON m.execution_id = e.id
		WHERE e.project_id = ?
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*models.Metric
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(
			&metric.ID,
			&metric.ExecutionID,
			&metric.Name,
			&metric.Value,
			&metric.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}
