package models

import "time"

// Project 项目模型
type Project struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Description  string    `json:"description"`
	RepositoryURL string    `json:"repository_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TechStack 技术栈模型
type TechStack struct {
	ID            int       `json:"id"`
	ProjectID     int       `json:"project_id"`
	Language      string    `json:"language"`
	Framework     string    `json:"framework"`
	BuildTool     string    `json:"build_tool"`
	TestFramework string    `json:"test_framework"`
	Dependencies  string    `json:"dependencies"` // JSON 格式
	CreatedAt     time.Time `json:"created_at"`
}

// Pipeline 管道配置模型
type Pipeline struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Platform  string    `json:"platform"`
	Config    string    `json:"config"` // YAML 格式
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Execution 执行历史模型
type Execution struct {
	ID         int       `json:"id"`
	ProjectID  int       `json:"project_id"`
	PipelineID int       `json:"pipeline_id"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"` // 秒
	Logs       string    `json:"logs"`
	CreatedAt  time.Time `json:"created_at"`
}

// Metric 指标模型
type Metric struct {
	ID          int       `json:"id"`
	ExecutionID int       `json:"execution_id"`
	Name        string    `json:"name"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
}

// Optimization 优化建议模型
type Optimization struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Suggestion  string    `json:"suggestion"`
	Applied     bool      `json:"applied"`
	CreatedAt   time.Time `json:"created_at"`
}
