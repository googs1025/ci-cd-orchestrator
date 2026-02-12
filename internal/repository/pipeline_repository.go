package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// PipelineRepository 管道配置仓库
type PipelineRepository struct {
	db *sql.DB
}

// NewPipelineRepository 创建管道配置仓库实例
func NewPipelineRepository(db *sql.DB) *PipelineRepository {
	return &PipelineRepository{db: db}
}

// Create 创建管道配置
func (r *PipelineRepository) Create(pipeline *models.Pipeline) error {
	query := `
		INSERT INTO pipelines (project_id, platform, config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, pipeline.ProjectID, pipeline.Platform, pipeline.Config, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	pipeline.ID = int(id)
	pipeline.CreatedAt = now
	pipeline.UpdatedAt = now

	return nil
}

// GetByProjectID 根据项目 ID 获取管道配置
func (r *PipelineRepository) GetByProjectID(projectID int) ([]*models.Pipeline, error) {
	query := `
		SELECT id, project_id, platform, config, created_at, updated_at
		FROM pipelines
		WHERE project_id = ?
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pipelines []*models.Pipeline
	for rows.Next() {
		var pipeline models.Pipeline
		err := rows.Scan(
			&pipeline.ID,
			&pipeline.ProjectID,
			&pipeline.Platform,
			&pipeline.Config,
			&pipeline.CreatedAt,
			&pipeline.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, &pipeline)
	}

	return pipelines, nil
}

// GetByID 根据 ID 获取管道配置
func (r *PipelineRepository) GetByID(id int) (*models.Pipeline, error) {
	query := `
		SELECT id, project_id, platform, config, created_at, updated_at
		FROM pipelines
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var pipeline models.Pipeline
	err := row.Scan(
		&pipeline.ID,
		&pipeline.ProjectID,
		&pipeline.Platform,
		&pipeline.Config,
		&pipeline.CreatedAt,
		&pipeline.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

// Update 更新管道配置
func (r *PipelineRepository) Update(pipeline *models.Pipeline) error {
	query := `
		UPDATE pipelines
		SET config = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query, pipeline.Config, now, pipeline.ID)
	if err != nil {
		return err
	}

	pipeline.UpdatedAt = now

	return nil
}

// Delete 删除管道配置
func (r *PipelineRepository) Delete(id int) error {
	query := `DELETE FROM pipelines WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
