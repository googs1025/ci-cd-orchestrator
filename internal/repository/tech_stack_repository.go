package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// TechStackRepository 技术栈仓库
type TechStackRepository struct {
	db *sql.DB
}

// NewTechStackRepository 创建技术栈仓库实例
func NewTechStackRepository(db *sql.DB) *TechStackRepository {
	return &TechStackRepository{db: db}
}

// Create 创建技术栈
func (r *TechStackRepository) Create(techStack *models.TechStack) error {
	query := `
		INSERT INTO tech_stacks (project_id, language, framework, build_tool, test_framework, dependencies, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		techStack.ProjectID,
		techStack.Language,
		techStack.Framework,
		techStack.BuildTool,
		techStack.TestFramework,
		techStack.Dependencies,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	techStack.ID = int(id)
	techStack.CreatedAt = now

	return nil
}

// GetByProjectID 根据项目 ID 获取技术栈
func (r *TechStackRepository) GetByProjectID(projectID int) (*models.TechStack, error) {
	query := `
		SELECT id, project_id, language, framework, build_tool, test_framework, dependencies, created_at
		FROM tech_stacks
		WHERE project_id = ?
	`

	row := r.db.QueryRow(query, projectID)

	var techStack models.TechStack
	err := row.Scan(
		&techStack.ID,
		&techStack.ProjectID,
		&techStack.Language,
		&techStack.Framework,
		&techStack.BuildTool,
		&techStack.TestFramework,
		&techStack.Dependencies,
		&techStack.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &techStack, nil
}

// Update 更新技术栈
func (r *TechStackRepository) Update(techStack *models.TechStack) error {
	query := `
		UPDATE tech_stacks
		SET language = ?, framework = ?, build_tool = ?, test_framework = ?, dependencies = ?
		WHERE project_id = ?
	`

	_, err := r.db.Exec(
		query,
		techStack.Language,
		techStack.Framework,
		techStack.BuildTool,
		techStack.TestFramework,
		techStack.Dependencies,
		techStack.ProjectID,
	)

	return err
}

// DeleteByProjectID 根据项目 ID 删除技术栈
func (r *TechStackRepository) DeleteByProjectID(projectID int) error {
	query := `DELETE FROM tech_stacks WHERE project_id = ?`
	_, err := r.db.Exec(query, projectID)
	return err
}
