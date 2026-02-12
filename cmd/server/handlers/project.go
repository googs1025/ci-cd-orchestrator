package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ci-cd-orchestrator/internal/db"
	"ci-cd-orchestrator/internal/models"
	"ci-cd-orchestrator/internal/repository"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectRepo *repository.ProjectRepository
}

// NewProjectHandler 创建项目处理器实例
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectRepo: repository.NewProjectRepository(db.GetDB()),
	}
}

// ListProjects 获取项目列表
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectRepo.GetAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取项目列表失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    projects,
		"message": "获取项目列表成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// CreateProject 创建项目
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"请求参数错误: ` + err.Error() + `"}`))
		return
	}

	// 验证参数
	if project.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"项目名称不能为空"}`))
		return
	}

	// 创建项目
	if err := h.projectRepo.Create(&project); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"创建项目失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":  "success",
		"data":    project,
		"message": "创建项目成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetProject 获取项目详情
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 获取项目
	project, err := h.projectRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取项目详情失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    project,
		"message": "获取项目详情成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// UpdateProject 更新项目
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 解析请求体
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"请求参数错误: ` + err.Error() + `"}`))
		return
	}

	// 设置项目ID
	project.ID = id

	// 验证参数
	if project.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"项目名称不能为空"}`))
		return
	}

	// 更新项目
	if err := h.projectRepo.Update(&project); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"更新项目失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    project,
		"message": "更新项目成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// DeleteProject 删除项目
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取项目ID
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的项目ID"}`))
		return
	}

	// 删除项目
	if err := h.projectRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"删除项目失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    nil,
		"message": "删除项目成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}
