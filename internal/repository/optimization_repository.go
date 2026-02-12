package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// OptimizationRepository 优化建议仓库
type OptimizationRepository struct {
	db *sql.DB
}

// NewOptimizationRepository 创建优化建议仓库实例
func NewOptimizationRepository(db *sql.DB) *OptimizationRepository {
	return &OptimizationRepository{db: db}
}

// Create 创建优化建议
func (r *OptimizationRepository) Create(optimization *models.Optimization) error {
	query := `
		INSERT INTO optimizations (project_id, type, description, suggestion, applied, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		optimization.ProjectID,
		optimization.Type,
		optimization.Description,
		optimization.Suggestion,
		optimization.Applied,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	optimization.ID = int(id)
	optimization.CreatedAt = now

	return nil
}

// GetByProjectID 根据项目 ID 获取优化建议
func (r *OptimizationRepository) GetByProjectID(projectID int) ([]*models.Optimization, error) {
	query := `
		SELECT id, project_id, type, description, suggestion, applied, created_at
		FROM optimizations
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var optimizations []*models.Optimization
	for rows.Next() {
		var optimization models.Optimization
		err := rows.Scan(
			&optimization.ID,
			&optimization.ProjectID,
			&optimization.Type,
			&optimization.Description,
			&optimization.Suggestion,
			&optimization.Applied,
			&optimization.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		optimizations = append(optimizations, &optimization)
	}

	return optimizations, nil
}

// MarkAsApplied 将优化建议标记为已应用
func (r *OptimizationRepository) MarkAsApplied(id int) error {
	query := `
		UPDATE optimizations
		SET applied = true
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	return err
}

// GetUnappliedByProjectID 获取项目的未应用优化建议
func (r *OptimizationRepository) GetUnappliedByProjectID(projectID int) ([]*models.Optimization, error) {
	query := `
		SELECT id, project_id, type, description, suggestion, applied, created_at
		FROM optimizations
		WHERE project_id = ? AND applied = false
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var optimizations []*models.Optimization
	for rows.Next() {
		var optimization models.Optimization
		err := rows.Scan(
			&optimization.ID,
			&optimization.ProjectID,
			&optimization.Type,
			&optimization.Description,
			&optimization.Suggestion,
			&optimization.Applied,
			&optimization.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		optimizations = append(optimizations, &optimization)
	}

	return optimizations, nil
}
