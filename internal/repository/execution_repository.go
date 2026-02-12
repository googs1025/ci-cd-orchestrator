package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// ExecutionRepository 执行历史仓库
type ExecutionRepository struct {
	db *sql.DB
}

// NewExecutionRepository 创建执行历史仓库实例
func NewExecutionRepository(db *sql.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

// Create 创建执行历史
func (r *ExecutionRepository) Create(execution *models.Execution) error {
	query := `
		INSERT INTO executions (project_id, pipeline_id, status, start_time, end_time, duration, logs, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		execution.ProjectID,
		execution.PipelineID,
		execution.Status,
		execution.StartTime,
		execution.EndTime,
		execution.Duration,
		execution.Logs,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	execution.ID = int(id)
	execution.CreatedAt = now

	return nil
}

// GetByProjectID 根据项目 ID 获取执行历史
func (r *ExecutionRepository) GetByProjectID(projectID int) ([]*models.Execution, error) {
	query := `
		SELECT id, project_id, pipeline_id, status, start_time, end_time, duration, logs, created_at
		FROM executions
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*models.Execution
	for rows.Next() {
		var execution models.Execution
		err := rows.Scan(
			&execution.ID,
			&execution.ProjectID,
			&execution.PipelineID,
			&execution.Status,
			&execution.StartTime,
			&execution.EndTime,
			&execution.Duration,
			&execution.Logs,
			&execution.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &execution)
	}

	return executions, nil
}

// GetByID 根据 ID 获取执行历史
func (r *ExecutionRepository) GetByID(id int) (*models.Execution, error) {
	query := `
		SELECT id, project_id, pipeline_id, status, start_time, end_time, duration, logs, created_at
		FROM executions
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var execution models.Execution
	err := row.Scan(
		&execution.ID,
		&execution.ProjectID,
		&execution.PipelineID,
		&execution.Status,
		&execution.StartTime,
		&execution.EndTime,
		&execution.Duration,
		&execution.Logs,
		&execution.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &execution, nil
}

// UpdateStatus 更新执行状态
func (r *ExecutionRepository) UpdateStatus(id int, status string, endTime time.Time, duration int, logs string) error {
	query := `
		UPDATE executions
		SET status = ?, end_time = ?, duration = ?, logs = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, status, endTime, duration, logs, id)
	return err
}
