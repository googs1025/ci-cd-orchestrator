package execution

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// ManagerImpl 执行管理器实现
type ManagerImpl struct {
	engines    map[string]Engine
	executions map[string]*Execution
	mutex      sync.RWMutex
}

// NewManager 创建执行管理器实例
func NewManager() Manager {
	return &ManagerImpl{
		engines:    make(map[string]Engine),
		executions: make(map[string]*Execution),
	}
}

// RegisterEngine 注册执行引擎
func (m *ManagerImpl) RegisterEngine(platform string, engine Engine) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.engines[platform] = engine
}

// CreateExecution 创建新的执行
func (m *ManagerImpl) CreateExecution(projectID, platform, triggerType string, options ExecutionOptions) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查平台是否支持
	_, exists := m.engines[platform]
	if !exists {
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	// 创建执行记录
	executionID := uuid.New().String()
	execution := &Execution{
		ID:           executionID,
		ProjectID:    projectID,
		Platform:     platform,
		Status:       StatusPending,
		TriggerType:  triggerType,
		TriggerInfo:  map[string]interface{}{},
		PlatformData: map[string]interface{}{},
		Metrics: Metrics{
			StageDurations: make(map[string]int64),
		},
		Logs: []LogEntry{},
	}

	// 存储执行记录
	m.executions[executionID] = execution

	// 注册到对应的引擎
	if mockEngine, ok := m.engines[platform].(*MockEngine); ok {
		mockEngine.RegisterExecution(execution)
	}

	return executionID, nil
}

// StartExecution 开始执行
func (m *ManagerImpl) StartExecution(executionID string) error {
	m.mutex.RLock()
	execution, exists := m.executions[executionID]
	if !exists {
		m.mutex.RUnlock()
		return fmt.Errorf("execution not found: %s", executionID)
	}

	engine, engineExists := m.engines[execution.Platform]
	if !engineExists {
		m.mutex.RUnlock()
		return fmt.Errorf("engine not found for platform: %s", execution.Platform)
	}
	m.mutex.RUnlock()

	// 检查执行状态
	if execution.Status != StatusPending {
		return fmt.Errorf("execution is not in pending status: %s", execution.Status)
	}

	// 构建执行选项
	options := ExecutionOptions{
		TotalDuration:   10,
		StageDurations:  nil,
		Result:          StatusSuccess,
		GenerateMetrics: true,
		GenerateLogs:    true,
		ResourceUsage: ResourceUsage{
			CpuUsage:    50.0,
			MemoryUsage: 60.0,
		},
	}

	// 启动执行
	return engine.Execute(executionID, options)
}

// StopExecution 停止执行
func (m *ManagerImpl) StopExecution(executionID string) error {
	m.mutex.RLock()
	execution, exists := m.executions[executionID]
	if !exists {
		m.mutex.RUnlock()
		return fmt.Errorf("execution not found: %s", executionID)
	}

	engine, engineExists := m.engines[execution.Platform]
	if !engineExists {
		m.mutex.RUnlock()
		return fmt.Errorf("engine not found for platform: %s", execution.Platform)
	}
	m.mutex.RUnlock()

	// 停止执行
	return engine.Stop(executionID)
}

// GetExecution 获取执行详情
func (m *ManagerImpl) GetExecution(executionID string) (*Execution, error) {
	m.mutex.RLock()
	execution, exists := m.executions[executionID]
	if !exists {
		m.mutex.RUnlock()
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	engine, engineExists := m.engines[execution.Platform]
	if !engineExists {
		m.mutex.RUnlock()
		return nil, fmt.Errorf("engine not found for platform: %s", execution.Platform)
	}
	m.mutex.RUnlock()

	// 从引擎获取最新状态
	return engine.GetStatus(executionID)
}

// ListExecutions 列出项目的执行历史
func (m *ManagerImpl) ListExecutions(projectID string, limit, offset int) ([]*Execution, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 过滤出指定项目的执行
	var executions []*Execution
	for _, execution := range m.executions {
		if execution.ProjectID == projectID {
			executions = append(executions, execution)
		}
	}

	// 应用分页
	if offset >= len(executions) {
		return []*Execution{}, nil
	}

	end := offset + limit
	if end > len(executions) {
		end = len(executions)
	}

	return executions[offset:end], nil
}
