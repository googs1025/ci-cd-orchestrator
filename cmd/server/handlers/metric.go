package handlers

import (
	"net/http"
)

// MetricHandler 指标处理器
type MetricHandler struct{}

// NewMetricHandler 创建指标处理器实例
func NewMetricHandler() *MetricHandler {
	return &MetricHandler{}
}

// ListMetrics 获取项目指标列表
func (h *MetricHandler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	// 简单返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","data":[],"message":"获取项目指标列表成功"}`))
}

// GetExecutionMetrics 获取执行指标详情
func (h *MetricHandler) GetExecutionMetrics(w http.ResponseWriter, r *http.Request) {
	// 简单返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","data":[],"message":"获取执行指标详情成功"}`))
}
