package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ci-cd-orchestrator/internal/db"
	"ci-cd-orchestrator/internal/execution"
	"ci-cd-orchestrator/internal/models"
	"ci-cd-orchestrator/internal/repository"
)

// OptimizationHandler 优化建议处理器
type OptimizationHandler struct {
	optimizationRepo *repository.OptimizationRepository
	executionRepo    *repository.ExecutionRepository
	metricRepo       *repository.MetricRepository
	executionManager execution.Manager
}

// NewOptimizationHandler 创建优化建议处理器实例
func NewOptimizationHandler(executionManager execution.Manager) *OptimizationHandler {
	return &OptimizationHandler{
		optimizationRepo: repository.NewOptimizationRepository(db.GetDB()),
		executionRepo:    repository.NewExecutionRepository(db.GetDB()),
		metricRepo:       repository.NewMetricRepository(db.GetDB()),
		executionManager: executionManager,
	}
}

// AnalyzeOptimization 分析优化建议
func (h *OptimizationHandler) AnalyzeOptimization(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	idStr = strings.TrimSuffix(idStr, "/analyze-optimization")
	projectID := idStr

	// 获取项目的执行历史
	executions, err := h.executionManager.ListExecutions(projectID, 100, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取执行历史失败: ` + err.Error() + `"}`))
		return
	}

	// 如果没有执行历史，返回提示
	if len(executions) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","data":{"suggestions":[],"metrics":null},"message":"暂无执行历史，无法生成优化建议"}`))
		return
	}

	// 如果执行历史少于5次，返回提示
	if len(executions) < 5 {
		// 计算关键指标
		metrics := h.calculateExecutionMetrics(executions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"suggestions": []map[string]interface{}{
					{
						"type":        "info",
						"description": "执行历史不足",
						"suggestion":  "需要至少5次执行历史才能生成优化建议",
					},
				},
				"metrics": metrics,
			},
			"message": "分析优化建议成功",
		}

		data, _ := json.Marshal(response)
		w.Write(data)
		return
	}

	// 分析执行历史，生成优化建议
	suggestions := h.analyzeExecutions(executions)

	// 保存优化建议到数据库
	for _, suggestion := range suggestions {
		optimization := &models.Optimization{
			ProjectID:   1, // 暂时使用固定值，实际应该从数据库中获取项目ID
			Type:        suggestion["type"].(string),
			Description: suggestion["description"].(string),
			Suggestion:  suggestion["suggestion"].(string),
			Applied:     false,
		}
		h.optimizationRepo.Create(optimization)
	}

	// 计算关键指标
	metrics := h.calculateExecutionMetrics(executions)

	// 确保 suggestions 不为 nil
	if suggestions == nil {
		suggestions = []map[string]interface{}{}
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"suggestions": suggestions,
			"metrics":     metrics,
		},
		"message": "分析优化建议成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// ListOptimizationSuggestions 获取优化建议列表
func (h *OptimizationHandler) ListOptimizationSuggestions(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	idStr = strings.TrimSuffix(idStr, "/optimization-suggestions")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 获取项目的优化建议
	optimizations, err := h.optimizationRepo.GetByProjectID(projectID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取优化建议列表失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    optimizations,
		"message": "获取优化建议列表成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// ApplyOptimization 应用优化建议
func (h *OptimizationHandler) ApplyOptimization(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	idStr = strings.TrimSuffix(idStr, "/apply-optimization")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 解析请求体
	var request struct {
		OptimizationID int `json:"optimization_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"请求参数错误: ` + err.Error() + `"}`))
		return
	}

	// 标记优化建议为已应用
	if err := h.optimizationRepo.MarkAsApplied(request.OptimizationID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"应用优化建议失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    nil,
		"message": "应用优化建议成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// analyzeExecutions 分析执行历史，生成优化建议
func (h *OptimizationHandler) analyzeExecutions(executions []*execution.Execution) []map[string]interface{} {
	var suggestions []map[string]interface{}

	// 分析执行时长
	if len(executions) >= 5 {
		totalDuration := int64(0)
		for _, exec := range executions {
			totalDuration += exec.Duration
		}
		averageDuration := float64(totalDuration) / float64(len(executions))

		if averageDuration > 15 {
			suggestions = append(suggestions, map[string]interface{}{
				"type":        "performance",
				"description": "执行时长过长",
				"suggestion":  "考虑优化构建和测试流程，减少执行时间",
			})
		}
	}

	// 分析成功/失败率
	successCount := 0
	failureCount := 0
	for _, exec := range executions {
		if exec.Status == "success" {
			successCount++
		} else if exec.Status == "failed" {
			failureCount++
		}
	}

	successRate := float64(successCount) / float64(len(executions))
	if successRate < 0.8 {
		suggestions = append(suggestions, map[string]interface{}{
			"type":        "reliability",
			"description": "成功率过低",
			"suggestion":  "分析失败原因，修复构建和测试中的问题",
		})
	}

	// 分析资源使用情况
	if len(executions) > 0 {
		totalCpuUsage := 0.0
		totalMemoryUsage := 0.0
		totalCpuRequest := 0.0
		totalMemoryRequest := 0.0
		totalCpuLimit := 0.0
		totalMemoryLimit := 0.0
		executionCount := 0

		for _, exec := range executions {
			if exec.PlatformData != nil {
				if resourceConfig, ok := exec.PlatformData["resource_config"].(map[string]interface{}); ok {
					if cpuData, ok := resourceConfig["cpu"].(map[string]interface{}); ok {
						if usagePercent, ok := cpuData["usage_percent"].(float64); ok {
							totalCpuUsage += usagePercent
						}
						if request, ok := cpuData["request"].(float64); ok {
							totalCpuRequest += request
						}
						if limit, ok := cpuData["limit"].(float64); ok {
							totalCpuLimit += limit
						}
					}

					if memoryData, ok := resourceConfig["memory"].(map[string]interface{}); ok {
						if usagePercent, ok := memoryData["usage_percent"].(float64); ok {
							totalMemoryUsage += usagePercent
						}
						if request, ok := memoryData["request"].(float64); ok {
							totalMemoryRequest += request
						}
						if limit, ok := memoryData["limit"].(float64); ok {
							totalMemoryLimit += limit
						}
					}

					executionCount++
				}
			}
		}

		if executionCount > 0 {
			averageCpuUsage := totalCpuUsage / float64(executionCount)
			averageMemoryUsage := totalMemoryUsage / float64(executionCount)
			averageCpuRequest := totalCpuRequest / float64(executionCount)
			averageMemoryRequest := totalMemoryRequest / float64(executionCount)
			averageCpuLimit := totalCpuLimit / float64(executionCount)
			averageMemoryLimit := totalMemoryLimit / float64(executionCount)

			// CPU 分析
			if averageCpuRequest > 0 {
				if averageCpuUsage < 40 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "CPU 使用率过低",
						"suggestion":  "考虑减小 CI 配置中的 CPU 请求值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageCpuUsage) + "，请求值为 " + fmt.Sprintf("%.2f", averageCpuRequest) + " 核",
					})
				} else if averageCpuUsage < 60 && averageCpuRequest > 1.0 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "CPU 使用率适中",
						"suggestion":  "考虑适当减小 CI 配置中的 CPU 请求值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageCpuUsage) + "，请求值为 " + fmt.Sprintf("%.2f", averageCpuRequest) + " 核",
					})
				}
			}

			// 内存分析
			if averageMemoryRequest > 0 {
				if averageMemoryUsage < 40 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "内存使用率过低",
						"suggestion":  "考虑减小 CI 配置中的内存请求值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageMemoryUsage) + "，请求值为 " + fmt.Sprintf("%.2f", averageMemoryRequest) + " GB",
					})
				} else if averageMemoryUsage < 60 && averageMemoryRequest > 1.5 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "内存使用率适中",
						"suggestion":  "考虑适当减小 CI 配置中的内存请求值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageMemoryUsage) + "，请求值为 " + fmt.Sprintf("%.2f", averageMemoryRequest) + " GB",
					})
				}
			}

			// 分析 Limit 配置
			if averageCpuLimit > 0 && averageCpuUsage < 50 {
				if averageCpuLimit > averageCpuRequest*1.5 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "CPU Limit 设置过高",
						"suggestion":  "考虑减小 CI 配置中的 CPU Limit 值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageCpuUsage) + "，Limit 值为 " + fmt.Sprintf("%.2f", averageCpuLimit) + " 核",
					})
				}
			}

			if averageMemoryLimit > 0 && averageMemoryUsage < 50 {
				if averageMemoryLimit > averageMemoryRequest*1.5 {
					suggestions = append(suggestions, map[string]interface{}{
						"type":        "resource",
						"description": "内存 Limit 设置过高",
						"suggestion":  "考虑减小 CI 配置中的内存 Limit 值，当前平均使用率为 " + fmt.Sprintf("%.2f%%", averageMemoryUsage) + "，Limit 值为 " + fmt.Sprintf("%.2f", averageMemoryLimit) + " GB",
					})
				}
			}
		}
	}

	// 总是生成一些优化建议，即使没有明显的问题
	if len(suggestions) == 0 {
		suggestions = append(suggestions, map[string]interface{}{
			"type":        "info",
			"description": "执行历史分析",
			"suggestion":  "当前执行历史表现良好，资源使用率合理，建议继续保持",
		})
	}

	return suggestions
}

// calculateExecutionMetrics 计算执行历史的关键指标
func (h *OptimizationHandler) calculateExecutionMetrics(executions []*execution.Execution) map[string]interface{} {
	metrics := make(map[string]interface{})

	// 计算平均执行时长
	totalDuration := int64(0)
	for _, exec := range executions {
		totalDuration += exec.Duration
	}
	averageDuration := float64(totalDuration) / float64(len(executions))
	metrics["average_duration"] = averageDuration

	// 计算成功/失败率
	successCount := 0
	failureCount := 0
	for _, exec := range executions {
		if exec.Status == "success" {
			successCount++
		} else if exec.Status == "failed" {
			failureCount++
		}
	}
	successRate := float64(successCount) / float64(len(executions))
	metrics["success_rate"] = successRate
	metrics["failure_rate"] = float64(failureCount) / float64(len(executions))

	// 计算最近 5 次执行的趋势
	if len(executions) >= 5 {
		recentExecutions := executions[len(executions)-5:]
		recentDurations := make([]int64, len(recentExecutions))
		recentCpuUsages := make([]float64, len(recentExecutions))
		recentMemoryUsages := make([]float64, len(recentExecutions))

		for i, exec := range recentExecutions {
			recentDurations[i] = exec.Duration

			// 提取 CPU 使用率
			if exec.Metrics.CpuUsage > 0 {
				recentCpuUsages[i] = exec.Metrics.CpuUsage
			} else if exec.PlatformData != nil {
				if resourceConfig, ok := exec.PlatformData["resource_config"].(map[string]interface{}); ok {
					if cpuData, ok := resourceConfig["cpu"].(map[string]interface{}); ok {
						if usagePercent, ok := cpuData["usage_percent"].(float64); ok {
							recentCpuUsages[i] = usagePercent
						}
					}
				}
			}

			// 提取内存使用率
			if exec.Metrics.MemoryUsage > 0 {
				recentMemoryUsages[i] = exec.Metrics.MemoryUsage
			} else if exec.PlatformData != nil {
				if resourceConfig, ok := exec.PlatformData["resource_config"].(map[string]interface{}); ok {
					if memoryData, ok := resourceConfig["memory"].(map[string]interface{}); ok {
						if usagePercent, ok := memoryData["usage_percent"].(float64); ok {
							recentMemoryUsages[i] = usagePercent
						}
					}
				}
			}
		}

		metrics["recent_durations"] = recentDurations
		metrics["recent_cpu_usages"] = recentCpuUsages
		metrics["recent_memory_usages"] = recentMemoryUsages
	}

	return metrics
}

// calculateKeyMetrics 计算关键指标
func (h *OptimizationHandler) calculateKeyMetrics(executions []*models.Execution) map[string]interface{} {
	metrics := make(map[string]interface{})

	// 计算平均执行时长
	totalDuration := 0
	for _, exec := range executions {
		totalDuration += exec.Duration
	}
	averageDuration := float64(totalDuration) / float64(len(executions))
	metrics["average_duration"] = averageDuration

	// 计算成功/失败率
	successCount := 0
	failureCount := 0
	for _, exec := range executions {
		if exec.Status == "success" {
			successCount++
		} else if exec.Status == "failed" {
			failureCount++
		}
	}
	successRate := float64(successCount) / float64(len(executions))
	metrics["success_rate"] = successRate
	metrics["failure_rate"] = float64(failureCount) / float64(len(executions))

	// 计算最近 5 次执行的趋势
	if len(executions) >= 5 {
		recentExecutions := executions[len(executions)-5:]
		recentDurations := make([]int, len(recentExecutions))
		for i, exec := range recentExecutions {
			recentDurations[i] = exec.Duration
		}
		metrics["recent_durations"] = recentDurations
	}

	return metrics
}
