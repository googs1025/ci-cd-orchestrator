package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ci-cd-orchestrator/internal/db"
	"ci-cd-orchestrator/internal/repository"
	"ci-cd-orchestrator/internal/techstack"
)

// TechStackHandler 技术栈处理器
type TechStackHandler struct {
	projectRepo *repository.ProjectRepository
}

// NewTechStackHandler 创建技术栈处理器实例
func NewTechStackHandler() *TechStackHandler {
	return &TechStackHandler{
		projectRepo: repository.NewProjectRepository(db.GetDB()),
	}
}

// AnalyzeProject 分析项目技术栈
func (h *TechStackHandler) AnalyzeProject(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取项目 ID
	path := r.URL.Path
	// 提取 /api/v1/projects/{id}/analyze 中的 id
	parts := strings.Split(path, "/")
	var projectID string
	for i, part := range parts {
		if part == "projects" && i+1 < len(parts) {
			projectID = parts[i+1]
			break
		}
	}

	if projectID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 从数据库中获取项目信息
	id, err := strconv.Atoi(projectID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	project, err := h.projectRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取项目信息失败: ` + err.Error() + `"}`))
		return
	}

	// 使用项目的实际路径
	projectPath := project.Path
	if projectPath == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"项目路径为空"}`))
		return
	}

	// 调用技术栈识别模块
	recognizer := techstack.NewRecognizer()
	result, err := recognizer.Recognize(projectPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"分析项目技术栈失败: ` + err.Error() + `"}`))
		return
	}

	// 移除 files 字段，因为不需要
	result.TechStack.Files = []string{}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    result,
		"message": "分析项目技术栈成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetTechStack 获取项目技术栈
func (h *TechStackHandler) GetTechStack(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取项目 ID
	path := r.URL.Path
	// 提取 /api/v1/projects/{id}/tech-stack 中的 id
	parts := strings.Split(path, "/")
	var projectID string
	for i, part := range parts {
		if part == "projects" && i+1 < len(parts) {
			projectID = parts[i+1]
			break
		}
	}

	if projectID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 从数据库中获取项目信息
	id, err := strconv.Atoi(projectID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	project, err := h.projectRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取项目信息失败: ` + err.Error() + `"}`))
		return
	}

	// 使用项目的实际路径
	projectPath := project.Path
	if projectPath == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"项目路径为空"}`))
		return
	}

	// 调用技术栈识别模块
	recognizer := techstack.NewRecognizer()
	result, err := recognizer.Recognize(projectPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取项目技术栈失败: ` + err.Error() + `"}`))
		return
	}

	// 移除 files 字段，因为不需要
	result.TechStack.Files = []string{}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    result.TechStack,
		"message": "获取项目技术栈成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}
