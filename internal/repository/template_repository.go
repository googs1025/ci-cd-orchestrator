package repository

import (
	"database/sql"
	"time"
)

// Template 模板模型
type Template struct {
	ID         int       `json:"id"`
	Platform   string    `json:"platform"`
	Language   string    `json:"language"`
	Framework  string    `json:"framework"`
	Content    string    `json:"content"`
	Filename   string    `json:"filename"`
	ConfigType string    `json:"config_type"`
	IsBuiltin  bool      `json:"is_builtin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TemplateRepository 模板仓库接口
type TemplateRepository interface {
	Create(template *Template) error
	GetByID(id int) (*Template, error)
	GetAll() ([]*Template, error)
	Update(template *Template) error
	Delete(id int) error
	GetByPlatformAndLanguage(platform string, language, framework string) (*Template, error)
	ResetBuiltinTemplates() error
}

// TemplateRepositoryImpl 模板仓库实现
type TemplateRepositoryImpl struct {
	db *sql.DB
}

// NewTemplateRepository 创建模板仓库实例
func NewTemplateRepository(db *sql.DB) TemplateRepository {
	return &TemplateRepositoryImpl{
		db: db,
	}
}

// Create 创建模板
func (r *TemplateRepositoryImpl) Create(template *Template) error {
	query := `
		INSERT INTO templates (platform, language, framework, content, filename, config_type, is_builtin, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, template.Platform, template.Language, template.Framework, template.Content, template.Filename, template.ConfigType, template.IsBuiltin, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	template.ID = int(id)
	template.CreatedAt = now
	template.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取模板
func (r *TemplateRepositoryImpl) GetByID(id int) (*Template, error) {
	query := `
		SELECT id, platform, language, framework, content, filename, config_type, is_builtin, created_at, updated_at
		FROM templates
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var template Template
	err := row.Scan(
		&template.ID,
		&template.Platform,
		&template.Language,
		&template.Framework,
		&template.Content,
		&template.Filename,
		&template.ConfigType,
		&template.IsBuiltin,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// GetAll 获取所有模板
func (r *TemplateRepositoryImpl) GetAll() ([]*Template, error) {
	query := `
		SELECT id, platform, language, framework, content, filename, config_type, is_builtin, created_at, updated_at
		FROM templates
		ORDER BY is_builtin DESC, platform, language, framework
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*Template
	for rows.Next() {
		var template Template
		err := rows.Scan(
			&template.ID,
			&template.Platform,
			&template.Language,
			&template.Framework,
			&template.Content,
			&template.Filename,
			&template.ConfigType,
			&template.IsBuiltin,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// Update 更新模板
func (r *TemplateRepositoryImpl) Update(template *Template) error {
	query := `
		UPDATE templates
		SET platform = ?, language = ?, framework = ?, content = ?, filename = ?, config_type = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query, template.Platform, template.Language, template.Framework, template.Content, template.Filename, template.ConfigType, now, template.ID)
	if err != nil {
		return err
	}

	template.UpdatedAt = now
	return nil
}

// Delete 删除模板
func (r *TemplateRepositoryImpl) Delete(id int) error {
	// 检查是否是内置模板
	var isBuiltin bool
	err := r.db.QueryRow("SELECT is_builtin FROM templates WHERE id = ?", id).Scan(&isBuiltin)
	if err != nil {
		return err
	}

	if isBuiltin {
		return nil // 内置模板不可删除
	}

	query := `DELETE FROM templates WHERE id = ?`
	_, err = r.db.Exec(query, id)
	return err
}

// GetByPlatformAndLanguage 根据平台、语言和框架获取模板
func (r *TemplateRepositoryImpl) GetByPlatformAndLanguage(platform string, language, framework string) (*Template, error) {
	query := `
		SELECT id, platform, language, framework, content, filename, config_type, is_builtin, created_at, updated_at
		FROM templates
		WHERE platform = ? AND language = ? AND framework = ?
		ORDER BY is_builtin DESC
		LIMIT 1
	`

	row := r.db.QueryRow(query, platform, language, framework)

	var template Template
	err := row.Scan(
		&template.ID,
		&template.Platform,
		&template.Language,
		&template.Framework,
		&template.Content,
		&template.Filename,
		&template.ConfigType,
		&template.IsBuiltin,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// ResetBuiltinTemplates 重置内置模板
func (r *TemplateRepositoryImpl) ResetBuiltinTemplates() error {
	// 这里将在后续实现，用于重置内置模板到默认状态
	return nil
}
