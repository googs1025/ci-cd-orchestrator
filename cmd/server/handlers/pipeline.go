package handlers

import (
	"ci-cd-orchestrator/internal/cicd"
	"ci-cd-orchestrator/internal/repository"
	"ci-cd-orchestrator/internal/techstack"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// PipelineHandler 管道配置处理器
type PipelineHandler struct {
	templateRepo repository.TemplateRepository
}

// NewPipelineHandler 创建管道配置处理器实例
func NewPipelineHandler(templateRepo repository.TemplateRepository) *PipelineHandler {
	return &PipelineHandler{
		templateRepo: templateRepo,
	}
}

// GeneratePipeline 生成管道配置
func (h *PipelineHandler) GeneratePipeline(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取平台参数，默认为 GitHub Actions
	platformStr := r.URL.Query().Get("platform")
	platform := cicd.PlatformGitHubActions
	if platformStr == "mock" {
		platform = cicd.PlatformMock
	}

	// 从请求中获取项目路径参数
	projectPath := r.URL.Query().Get("path")
	if projectPath == "" {
		projectPath = "."
	}

	// 从请求中获取模板ID参数
	var templateID int
	templateIDStr := r.URL.Query().Get("template_id")
	if templateIDStr != "" {
		if id, err := strconv.Atoi(templateIDStr); err == nil {
			templateID = id
		}
	}

	// 调用技术栈识别模块获取实际的技术栈信息
	recognizer := techstack.NewRecognizer()
	techStackResult, err := recognizer.Recognize(projectPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"分析项目技术栈失败: ` + err.Error() + `"}`))
		return
	}

	// 使用 CI/CD 生成器生成配置
	generator := cicd.NewGenerator(h.templateRepo)
	var config *cicd.PipelineConfig

	if templateID > 0 {
		config, err = generator.GenerateConfig(&techStackResult.TechStack, platform, templateID)
	} else {
		config, err = generator.GenerateConfig(&techStackResult.TechStack, platform)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"生成管道配置失败: ` + err.Error() + `"}`))
		return
	}

	// 将配置内容写入到项目目录中

	// 构建完整的文件路径
	configPath := filepath.Join(projectPath, config.Filename)

	// 检查文件是否存在
	if _, err := os.Stat(configPath); err == nil {
		// 文件已存在，不生成新配置
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":  "success",
			"data":    config,
			"message": "配置文件已存在，跳过生成",
		}

		data, _ := json.Marshal(response)
		w.Write(data)
		return
	}

	// 确保目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"创建配置目录失败: ` + err.Error() + `"}`))
		return
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(config.Content), 0644); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"写入配置文件失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    config,
		"message": "生成管道配置成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// GetPipeline 获取管道配置
func (h *PipelineHandler) GetPipeline(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取平台参数，默认为 GitHub Actions
	platformStr := r.URL.Query().Get("platform")
	platform := cicd.PlatformGitHubActions
	if platformStr == "mock" {
		platform = cicd.PlatformMock
	}

	// 从请求中获取技术栈信息（这里简化处理，暂时使用一个模拟的技术栈）
	// 实际应用中，应该从项目的技术栈分析结果中获取
	mockTechStack := &techstack.TechStack{
		Language:      "Go",
		Framework:     "",
		BuildTool:     "go mod",
		TestFramework: "testing",
		Dependencies: map[string]string{
			"github.com/gin-gonic/gin": "v1.9.0",
		},
		Files: []string{},
	}

	// 使用 CI/CD 生成器生成配置
	generator := cicd.NewGenerator(h.templateRepo)
	config, err := generator.GenerateConfig(mockTechStack, platform)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"获取管道配置失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    config,
		"message": "获取管道配置成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}

// UpdatePipeline 更新管道配置
func (h *PipelineHandler) UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取管道配置（这里简化处理，暂时使用一个模拟的配置）
	// 实际应用中，应该从请求体中解析配置
	mockConfig := &cicd.PipelineConfig{
		Platform:   cicd.PlatformGitHubActions,
		ConfigType: cicd.ConfigTypeYAML,
		Content: `name: Updated CI

on:
  push:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Build
      run: go build -v ./...
    - name: Test
      run: go test -v ./...
`,
		Filename: ".github/workflows/ci.yml",
	}

	// 使用 CI/CD 生成器验证配置
	generator := cicd.NewGenerator(h.templateRepo)
	err := generator.ValidateConfig(mockConfig)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","data":null,"message":"更新管道配置失败: ` + err.Error() + `"}`))
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"data":    mockConfig,
		"message": "更新管道配置成功",
	}

	data, _ := json.Marshal(response)
	w.Write(data)
}
