package repository

import (
	"database/sql"
	"time"

	"ci-cd-orchestrator/internal/models"
)

// ProjectRepository 项目仓库
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository 创建项目仓库实例
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create 创建项目
func (r *ProjectRepository) Create(project *models.Project) error {
	query := `
		INSERT INTO projects (name, path, description, repository_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, project.Name, project.Path, project.Description, project.RepositoryURL, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	project.ID = int(id)
	project.CreatedAt = now
	project.UpdatedAt = now

	return nil
}

// GetByID 根据 ID 获取项目
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	query := `
		SELECT id, name, path, description, repository_url, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var project models.Project
	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Path,
		&project.Description,
		&project.RepositoryURL,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetAll 获取所有项目
func (r *ProjectRepository) GetAll() ([]*models.Project, error) {
	query := `
		SELECT id, name, path, description, repository_url, created_at, updated_at
		FROM projects
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Path,
			&project.Description,
			&project.RepositoryURL,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

// Update 更新项目
func (r *ProjectRepository) Update(project *models.Project) error {
	query := `
		UPDATE projects
		SET name = ?, path = ?, description = ?, repository_url = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query, project.Name, project.Path, project.Description, project.RepositoryURL, now, project.ID)
	if err != nil {
		return err
	}

	project.UpdatedAt = now
	return nil
}

// Delete 删除项目
func (r *ProjectRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
