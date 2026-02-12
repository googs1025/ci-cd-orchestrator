package execution

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// CIConfig CI 配置结构
type CIConfig struct {
	Jobs map[string]Job `yaml:"jobs"`
}

// Job 任务结构
type Job struct {
	RunsOn    string    `yaml:"runs-on,omitempty"`
	Resources Resources `yaml:"resources,omitempty"`
	Steps     []Step    `yaml:"steps,omitempty"`
}

// Step 步骤结构
type Step struct {
	Name string   `yaml:"name"`
	Run  string   `yaml:"run"`
	Uses string   `yaml:"uses,omitempty"`
	With WithData `yaml:"with,omitempty"`
}

// WithData 步骤参数结构
type WithData map[string]interface{}

// Resources 资源配置结构
type Resources struct {
	Limits   ResourceLimit `yaml:"limits"`
	Requests ResourceLimit `yaml:"requests"`
}

// ResourceLimit 资源限制结构
type ResourceLimit struct {
	CPU    interface{} `yaml:"cpu"`
	Memory string      `yaml:"memory"`
}

// MockEngine Mock CI 执行引擎
type MockEngine struct {
	executions map[string]*Execution
	mutex      sync.RWMutex
}

// NewMockEngine 创建 Mock CI 执行引擎实例
func NewMockEngine() Engine {
	return &MockEngine{
		executions: make(map[string]*Execution),
	}
}

// Execute 执行模拟 CI 流程
func (e *MockEngine) Execute(executionID string, options ExecutionOptions) error {
	e.mutex.Lock()
	execution, exists := e.executions[executionID]
	if !exists {
		e.mutex.Unlock()
		return fmt.Errorf("execution not found: %s", executionID)
	}

	// 更新状态为运行中
	execution.Status = StatusRunning
	execution.StartTime = time.Now()
	e.mutex.Unlock()

	// 异步执行模拟流程
	go func() {
		e.simulateExecution(executionID, options)
	}()

	return nil
}

// parseMemory 解析内存字符串为字节数
func parseMemory(memoryStr string) int64 {
	if memoryStr == "" {
		return 1024 * 1024 * 1024 // 默认 1GB
	}

	// 简单实现，支持 GB 和 MB 单位
	var value int64
	var unit string
	fmt.Sscanf(memoryStr, "%d%s", &value, &unit)

	switch unit {
	case "G", "GB":
		return value * 1024 * 1024 * 1024
	case "M", "MB":
		return value * 1024 * 1024
	default:
		return value
	}
}

// parseCPU 解析 CPU 值
func parseCPU(cpu interface{}) float64 {
	switch v := cpu.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		var value float64
		fmt.Sscanf(v, "%f", &value)
		return value
	default:
		return 1.0 // 默认 1 核
	}
}

// Stop 停止执行
func (e *MockEngine) Stop(executionID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	if execution.Status != StatusRunning {
		return fmt.Errorf("execution is not running: %s", execution.Status)
	}

	// 更新状态为已取消
	execution.Status = StatusCancelled
	execution.EndTime = time.Now()
	execution.Duration = int64(time.Since(execution.StartTime).Seconds())

	// 添加取消日志
	if execution.Logs == nil {
		execution.Logs = []LogEntry{}
	}
	execution.Logs = append(execution.Logs, LogEntry{
		ID:          uuid.New().String(),
		ExecutionID: executionID,
		Timestamp:   time.Now(),
		Level:       "info",
		Stage:       "cancellation",
		Message:     "Execution cancelled by user",
	})

	return nil
}

// GetStatus 获取执行状态
func (e *MockEngine) GetStatus(executionID string) (*Execution, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("execution not found: %s", executionID)
	}

	// 返回执行的副本，避免并发修改问题
	copy := *execution
	if execution.Logs != nil {
		copy.Logs = make([]LogEntry, len(execution.Logs))
		for i, log := range execution.Logs {
			copy.Logs[i] = log
		}
	}

	return &copy, nil
}

// RegisterExecution 注册执行记录
func (e *MockEngine) RegisterExecution(execution *Execution) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.executions[execution.ID] = execution
}

// simulateExecution 模拟执行流程
func (e *MockEngine) simulateExecution(executionID string, options ExecutionOptions) {
	// 模拟阶段执行
	stages := []string{"init", "build", "test", "deploy", "complete"}
	stageEndTimes := make(map[string]time.Time)

	for _, stage := range stages {
		// 检查是否已被停止
		e.mutex.RLock()
		currentExecution := e.executions[executionID]
		if currentExecution.Status == StatusCancelled {
			e.mutex.RUnlock()
			return
		}
		e.mutex.RUnlock()

		// 添加阶段开始日志
		e.addLog(executionID, "info", stage, fmt.Sprintf("Starting %s stage", stage))

		// 模拟阶段执行
		duration := 1 + rand.Intn(3) // 1-4 秒

		// 如果是build或test阶段，执行多个step
		if (stage == "build" || stage == "test") && options.CIConfigContent != "" {
			// 解析CI配置
			var config CIConfig
			err := yaml.Unmarshal([]byte(options.CIConfigContent), &config)
			if err == nil && len(config.Jobs) > 0 {
				// 执行第一个job的所有steps
				for _, job := range config.Jobs {
					if len(job.Steps) > 0 {
						// 执行每个step
						for i, step := range job.Steps {
							// 添加step开始日志
							e.addLogWithStep(executionID, "info", stage, fmt.Sprintf("step-%d", i+1), fmt.Sprintf("Starting step: %s", step.Name))

							// 模拟step执行
							stepDuration := 1 + rand.Intn(2) // 1-3 秒
							time.Sleep(time.Duration(stepDuration) * time.Second)

							// 添加step完成日志
							e.addLogWithStep(executionID, "info", stage, fmt.Sprintf("step-%d", i+1), fmt.Sprintf("Completed step: %s in %d seconds", step.Name, stepDuration))

							// 检查是否需要模拟失败
							if options.Result == StatusFailed && options.FailureStage == stage {
								e.failExecution(executionID, stage, options.FailureReason, options.CIConfigContent)
								return
							}
						}
						continue
					}
				}
			}
		} else {
			// 普通阶段执行
			time.Sleep(time.Duration(duration) * time.Second)
		}

		stageEndTimes[stage] = time.Now()

		// 检查是否需要模拟失败
		if options.Result == StatusFailed && options.FailureStage == stage {
			e.failExecution(executionID, stage, options.FailureReason, options.CIConfigContent)
			return
		}

		// 添加阶段完成日志
		e.addLog(executionID, "info", stage, fmt.Sprintf("Completed %s stage in %d seconds", stage, duration))
	}

	// 执行成功完成
	e.completeExecution(executionID, stageEndTimes, options)
}

// addLogWithStep 添加带步骤信息的日志条目
func (e *MockEngine) addLogWithStep(executionID, level, stage, step, message string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return
	}

	if execution.Logs == nil {
		execution.Logs = []LogEntry{}
	}

	execution.Logs = append(execution.Logs, LogEntry{
		ID:          uuid.New().String(),
		ExecutionID: executionID,
		Timestamp:   time.Now(),
		Level:       level,
		Stage:       stage,
		Step:        step,
		Message:     message,
	})
}

// addLog 添加日志条目
func (e *MockEngine) addLog(executionID, level, stage, message string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return
	}

	if execution.Logs == nil {
		execution.Logs = []LogEntry{}
	}

	execution.Logs = append(execution.Logs, LogEntry{
		ID:          uuid.New().String(),
		ExecutionID: executionID,
		Timestamp:   time.Now(),
		Level:       level,
		Stage:       stage,
		Message:     message,
	})
}

// failExecution 模拟执行失败
func (e *MockEngine) failExecution(executionID, stage, reason string, ciConfigContent string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return
	}

	// 更新状态为失败
	execution.Status = StatusFailed
	execution.EndTime = time.Now()
	execution.Duration = int64(time.Since(execution.StartTime).Seconds())

	// 添加失败日志
	if execution.Logs == nil {
		execution.Logs = []LogEntry{}
	}
	execution.Logs = append(execution.Logs, LogEntry{
		ID:          uuid.New().String(),
		ExecutionID: executionID,
		Timestamp:   time.Now(),
		Level:       "error",
		Stage:       stage,
		Message:     fmt.Sprintf("Failed at %s stage: %s", stage, reason),
	})

	// 生成失败指标
	e.generateMetrics(execution, false, ciConfigContent)
}

// completeExecution 模拟执行成功完成
func (e *MockEngine) completeExecution(executionID string, stageEndTimes map[string]time.Time, options ExecutionOptions) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	execution, exists := e.executions[executionID]
	if !exists {
		return
	}

	// 更新状态为成功
	execution.Status = StatusSuccess
	execution.EndTime = time.Now()
	execution.Duration = int64(time.Since(execution.StartTime).Seconds())

	// 添加完成日志
	if execution.Logs == nil {
		execution.Logs = []LogEntry{}
	}
	execution.Logs = append(execution.Logs, LogEntry{
		ID:          uuid.New().String(),
		ExecutionID: executionID,
		Timestamp:   time.Now(),
		Level:       "info",
		Stage:       "complete",
		Message:     fmt.Sprintf("Execution completed successfully in %d seconds", execution.Duration),
	})

	// 生成成功指标
	e.generateMetrics(execution, true, options.CIConfigContent)
}

// generateMetrics 生成执行指标
func (e *MockEngine) generateMetrics(execution *Execution, success bool, ciConfigContent string) {
	// 基础指标
	execution.Metrics.TotalDuration = execution.Duration
	execution.Metrics.StageDurations = map[string]int64{
		"init":     2,
		"build":    3,
		"test":     3,
		"deploy":   1,
		"complete": 1,
	}

	// 成功/失败率
	if success {
		execution.Metrics.SuccessRate = 1.0
	} else {
		execution.Metrics.SuccessRate = 0.0
	}

	// 解析 CI 配置
	var cpuLimit, memoryLimit float64
	var cpuRequest, memoryRequest float64

	if ciConfigContent != "" {
		var config CIConfig
		err := yaml.Unmarshal([]byte(ciConfigContent), &config)
		if err == nil {
			// 提取资源配置
			for _, job := range config.Jobs {
				if job.Resources.Limits.CPU != nil {
					cpuLimit = parseCPU(job.Resources.Limits.CPU)
				} else {
					cpuLimit = 2.0 // 默认 2 核
				}

				if job.Resources.Requests.CPU != nil {
					cpuRequest = parseCPU(job.Resources.Requests.CPU)
				} else {
					cpuRequest = 1.0 // 默认 1 核
				}

				if job.Resources.Limits.Memory != "" {
					memoryLimitBytes := parseMemory(job.Resources.Limits.Memory)
					memoryLimit = float64(memoryLimitBytes) / (1024 * 1024 * 1024) // 转换为 GB
				} else {
					memoryLimit = 4.0 // 默认 4GB
				}

				if job.Resources.Requests.Memory != "" {
					memoryRequestBytes := parseMemory(job.Resources.Requests.Memory)
					memoryRequest = float64(memoryRequestBytes) / (1024 * 1024 * 1024) // 转换为 GB
				} else {
					memoryRequest = 2.0 // 默认 2GB
				}
				break
			}
		}
	} else {
		// 默认值
		cpuLimit = 2.0
		cpuRequest = 1.0
		memoryLimit = 4.0
		memoryRequest = 2.0
	}

	// 基于配置生成资源使用情况
	// CPU 使用率：配置请求的 30-95%
	cpuUsagePercent := 30.0 + rand.Float64()*65.0 // 30-95%
	execution.Metrics.CpuUsage = cpuUsagePercent

	// 内存使用率：配置请求的 40-95%
	memoryUsagePercent := 40.0 + rand.Float64()*55.0 // 40-95%
	execution.Metrics.MemoryUsage = memoryUsagePercent

	// 测试覆盖率（模拟）
	execution.Metrics.TestCoverage = 50.0 + rand.Float64()*45.0 // 50-95%

	// 构建大小（模拟）
	execution.Metrics.BuildSize = int64(1024 * 1024 * (5 + rand.Intn(95))) // 5-100 MB

	// 部署时间（模拟）
	execution.Metrics.DeploymentTime = 1 + int64(rand.Intn(10)) // 1-10 seconds

	// 添加资源配置与实际使用的比较信息
	execution.PlatformData = map[string]interface{}{
		"resource_config": map[string]interface{}{
			"cpu": map[string]interface{}{
				"limit":         cpuLimit,
				"request":       cpuRequest,
				"actual":        cpuUsagePercent / 100.0 * cpuRequest, // 实际使用的 CPU
				"usage_percent": cpuUsagePercent,
			},
			"memory": map[string]interface{}{
				"limit":         memoryLimit,
				"request":       memoryRequest,
				"actual":        memoryUsagePercent / 100.0 * memoryRequest, // 实际使用的内存
				"usage_percent": memoryUsagePercent,
			},
		},
	}
}
