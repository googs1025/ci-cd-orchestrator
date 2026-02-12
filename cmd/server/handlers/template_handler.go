package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ci-cd-orchestrator/internal/repository"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	templateRepo repository.TemplateRepository
}

// NewTemplateHandler 创建模板处理器实例
func NewTemplateHandler(templateRepo repository.TemplateRepository) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
	}
}

// CreateTemplate 创建模板
func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template repository.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"请求参数错误: ` + err.Error() + `"}`))
		return
	}

	// 验证必要字段
	if template.Platform == "" || template.Content == "" || template.Filename == "" || template.ConfigType == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"平台、内容、文件名和配置类型为必填字段"}`))
		return
	}

	// 非内置模板
	template.IsBuiltin = false

	if err := h.templateRepo.Create(&template); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"创建模板失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":  "success",
		"data":    template,
		"message": "创建模板成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetTemplate 获取单个模板
func (h *TemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取模板ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var templateID string
	for i, part := range parts {
		if part == "templates" && i+1 < len(parts) {
			templateID = parts[i+1]
			break
		}
	}

	if templateID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	id, err := strconv.Atoi(templateID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	template, err := h.templateRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取模板失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    template,
		"message": "获取模板成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// ListTemplates 获取所有模板
func (h *TemplateHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.templateRepo.GetAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取模板列表失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    templates,
		"message": "获取模板列表成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// UpdateTemplate 更新模板
func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取模板ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var templateID string
	for i, part := range parts {
		if part == "templates" && i+1 < len(parts) {
			templateID = parts[i+1]
			break
		}
	}

	if templateID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	id, err := strconv.Atoi(templateID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	var template repository.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"请求参数错误: ` + err.Error() + `"}`))
		return
	}

	// 验证必要字段
	if template.Platform == "" || template.Content == "" || template.Filename == "" || template.ConfigType == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"平台、内容、文件名和配置类型为必填字段"}`))
		return
	}

	// 确保ID一致
	template.ID = id

	// 保留内置模板标记
	existingTemplate, err := h.templateRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取模板失败: ` + err.Error() + `"}`))
		return
	}
	template.IsBuiltin = existingTemplate.IsBuiltin

	if err := h.templateRepo.Update(&template); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"更新模板失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    template,
		"message": "更新模板成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// DeleteTemplate 删除模板
func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取模板ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var templateID string
	for i, part := range parts {
		if part == "templates" && i+1 < len(parts) {
			templateID = parts[i+1]
			break
		}
	}

	if templateID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	id, err := strconv.Atoi(templateID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	if err := h.templateRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"删除模板失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    nil,
		"message": "删除模板成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// ResetTemplate 重置内置模板
func (h *TemplateHandler) ResetTemplate(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取模板ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var templateID string
	for i, part := range parts {
		if part == "templates" && i+1 < len(parts) {
			templateID = parts[i+1]
			break
		}
	}

	if templateID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","data":null,"message":"无效的模板ID"}`))
		return
	}

	// 这里简化处理，实际应该重置为默认模板内容
	// 后续实现

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    nil,
		"message": "重置模板成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// InitializeBuiltinTemplates 初始化内置模板
func (h *TemplateHandler) InitializeBuiltinTemplates() error {
	// 这里将在后续实现，用于初始化内置模板
	return nil
}
