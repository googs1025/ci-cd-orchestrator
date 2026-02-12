package validator

import (
	"fmt"
	"strings"

	"ci-cd-orchestrator/internal/cicd/common"
)

// Validator 配置验证器接口
type Validator interface {
	Validate(config *common.PipelineConfig) error
	ValidateGitHubActionsConfig(content string) error
	ValidateMockConfig(content string) error
}

// ConfigValidator 配置验证器实现
type ConfigValidator struct{}

// NewValidator 创建配置验证器实例
func NewValidator() Validator {
	return &ConfigValidator{}
}

// Validate 验证 CI/CD 管道配置
func (v *ConfigValidator) Validate(config *common.PipelineConfig) error {
	// 检查配置是否为空
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// 检查平台是否有效
	if config.Platform == "" {
		return fmt.Errorf("platform cannot be empty")
	}

	// 检查配置类型是否有效
	if config.ConfigType == "" {
		return fmt.Errorf("config_type cannot be empty")
	}

	// 检查配置内容是否为空
	if config.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	// 检查文件名是否为空
	if config.Filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// 根据平台验证配置
	switch config.Platform {
	case common.PlatformGitHubActions:
		return v.ValidateGitHubActionsConfig(config.Content)
	case common.PlatformMock:
		return v.ValidateMockConfig(config.Content)
	default:
		return fmt.Errorf("unsupported platform: %s", config.Platform)
	}
}

// ValidateGitHubActionsConfig 验证 GitHub Actions 配置
func (v *ConfigValidator) ValidateGitHubActionsConfig(content string) error {
	// 检查配置是否包含必要的字段
	lines := strings.Split(content, "\n")

	// 检查是否包含 name 字段
	hasName := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			hasName = true
			break
		}
	}
	if !hasName {
		return fmt.Errorf("GitHub Actions config must have a 'name' field")
	}

	// 检查是否包含 on 字段
	hasOn := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "on:") {
			hasOn = true
			break
		}
	}
	if !hasOn {
		return fmt.Errorf("GitHub Actions config must have an 'on' field")
	}

	// 检查是否包含 jobs 字段
	hasJobs := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "jobs:") {
			hasJobs = true
			break
		}
	}
	if !hasJobs {
		return fmt.Errorf("GitHub Actions config must have a 'jobs' field")
	}

	return nil
}

// ValidateMockConfig 验证 Mock 平台配置
func (v *ConfigValidator) ValidateMockConfig(content string) error {
	// 检查配置是否包含必要的字段
	lines := strings.Split(content, "\n")

	// 检查是否包含 name 字段
	hasName := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			hasName = true
			break
		}
	}
	if !hasName {
		return fmt.Errorf("Mock config must have a 'name' field")
	}

	// 检查是否包含 on 字段
	hasOn := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "on:") {
			hasOn = true
			break
		}
	}
	if !hasOn {
		return fmt.Errorf("Mock config must have an 'on' field")
	}

	// 检查是否包含 jobs 字段
	hasJobs := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "jobs:") {
			hasJobs = true
			break
		}
	}
	if !hasJobs {
		return fmt.Errorf("Mock config must have a 'jobs' field")
	}

	return nil
}
