package execution

import (
	"time"
)

// Status 执行状态
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusSuccess   = "success"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

// Execution 执行记录
type Execution struct {
	ID           string                 `json:"id"`
	ProjectID    string                 `json:"project_id"`
	Platform     string                 `json:"platform"`
	Status       string                 `json:"status"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Duration     int64                  `json:"duration"`
	TriggerType  string                 `json:"trigger_type"`
	TriggerInfo  map[string]interface{} `json:"trigger_info"`
	PlatformData map[string]interface{} `json:"platform_data"`
	Metrics      Metrics                `json:"metrics"`
	Logs         []LogEntry             `json:"logs,omitempty"`
}

// Metrics 执行指标
type Metrics struct {
	TotalDuration  int64            `json:"total_duration"`
	StageDurations map[string]int64 `json:"stage_durations"`
	SuccessRate    float64          `json:"success_rate"`
	CpuUsage       float64          `json:"cpu_usage,omitempty"`
	MemoryUsage    float64          `json:"memory_usage,omitempty"`
	TestCoverage   float64          `json:"test_coverage,omitempty"`
	BuildSize      int64            `json:"build_size,omitempty"`
	DeploymentTime int64            `json:"deployment_time,omitempty"`
}

// LogEntry 日志条目
type LogEntry struct {
	ID          string                 `json:"id"`
	ExecutionID string                 `json:"execution_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Stage       string                 `json:"stage"`
	Step        string                 `json:"step,omitempty"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// ExecutionOptions 执行选项
type ExecutionOptions struct {
	TotalDuration   int            `json:"total_duration"`
	StageDurations  map[string]int `json:"stage_durations"`
	Result          string         `json:"result"`
	FailureStage    string         `json:"failure_stage"`
	FailureReason   string         `json:"failure_reason"`
	GenerateMetrics bool           `json:"generate_metrics"`
	GenerateLogs    bool           `json:"generate_logs"`
	ResourceUsage   ResourceUsage  `json:"resource_usage"`
	CIConfigContent string         `json:"ci_config_content"` // CI 配置文件内容
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

// Engine 执行引擎接口
type Engine interface {
	Execute(executionID string, options ExecutionOptions) error
	Stop(executionID string) error
	GetStatus(executionID string) (*Execution, error)
}

// Manager 执行管理器接口
type Manager interface {
	CreateExecution(projectID, platform, triggerType string, options ExecutionOptions) (string, error)
	StartExecution(executionID string) error
	StopExecution(executionID string) error
	GetExecution(executionID string) (*Execution, error)
	ListExecutions(projectID string, limit, offset int) ([]*Execution, error)
	RegisterEngine(platform string, engine Engine)
}
