package handlers

import (
	"ci-cd-orchestrator/internal/execution"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ExecutionHandler 执行处理器
type ExecutionHandler struct {
	manager execution.Manager
}

// NewExecutionHandler 创建执行处理器实例
func NewExecutionHandler(manager execution.Manager) *ExecutionHandler {
	return &ExecutionHandler{
		manager: manager,
	}
}

// ExecutePipeline 执行管道
func (h *ExecutionHandler) ExecutePipeline(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取平台参数，默认为 GitHub Actions
	platform := r.URL.Query().Get("platform")
	if platform == "" {
		platform = "github_actions"
	}

	// 从路径参数中获取项目 ID
	path := r.URL.Path
	// 提取 /api/v1/projects/{id}/execute 中的 id
	parts := strings.Split(path, "/")
	var projectID string
	for i, part := range parts {
		if part == "projects" && i+1 < len(parts) {
			projectID = parts[i+1]
			break
		}
	}

	// 读取 CI 配置文件内容
	ciConfigContent := ""
	if platform == "mock" {
		// 这里可以从项目目录中读取 CI 配置文件
		// 简化实现，直接使用默认的 Go 项目 CI 配置
		ciConfigContent = `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    # 资源配置
    resources:
      limits:
        cpu: 2
        memory: 4G
      requests:
        cpu: 1
        memory: 2G
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20
    - name: Build
      run: go build -v ./...
    - name: Test
      run: go test -v ./...
`
	}

	// 创建执行
	executionID, err := h.manager.CreateExecution(projectID, platform, "manual", execution.ExecutionOptions{
		TotalDuration:   10,
		GenerateMetrics: true,
		GenerateLogs:    true,
		CIConfigContent: ciConfigContent,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"创建执行失败: ` + err.Error() + `"}`))
		return
	}

	// 启动执行
	if err := h.manager.StartExecution(executionID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"启动执行失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"execution_id": executionID,
			"platform":     platform,
			"status":       "running",
		},
		"message": "执行管道成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetExecution 获取执行详情
func (h *ExecutionHandler) GetExecution(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取执行 ID
	executionID := r.URL.Path[len("/api/v1/executions/"):]

	// 获取执行详情
	execution, err := h.manager.GetExecution(executionID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取执行详情失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    execution,
		"message": "获取执行详情成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// StopExecution 停止执行
func (h *ExecutionHandler) StopExecution(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取执行 ID
	executionID := r.URL.Path[len("/api/v1/executions/") : len("/api/v1/executions/")+36]

	// 停止执行
	if err := h.manager.StopExecution(executionID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"停止执行失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"execution_id": executionID,
			"status":       "cancelled",
		},
		"message": "停止执行成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// ListExecutions 获取执行历史
func (h *ExecutionHandler) ListExecutions(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取项目 ID
	path := r.URL.Path
	// 提取 /api/v1/projects/{id}/executions 中的 id
	parts := strings.Split(path, "/")
	var projectID string
	for i, part := range parts {
		if part == "projects" && i+1 < len(parts) {
			projectID = parts[i+1]
			break
		}
	}

	// 从查询参数中获取分页参数
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	// 获取执行历史
	executions, err := h.manager.ListExecutions(projectID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取执行历史失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    executions,
		"message": "获取执行历史成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetExecutionMetrics 获取执行指标
func (h *ExecutionHandler) GetExecutionMetrics(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取执行 ID
	executionID := r.URL.Path[len("/api/v1/executions/") : len("/api/v1/executions/")+36]

	// 获取执行详情
	execution, err := h.manager.GetExecution(executionID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取执行详情失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    execution.Metrics,
		"message": "获取执行指标成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetExecutionLogs 获取执行日志
func (h *ExecutionHandler) GetExecutionLogs(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取执行 ID
	executionID := r.URL.Path[len("/api/v1/executions/") : len("/api/v1/executions/")+36]

	// 获取执行详情
	execution, err := h.manager.GetExecution(executionID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取执行详情失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    execution.Logs,
		"message": "获取执行日志成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}
