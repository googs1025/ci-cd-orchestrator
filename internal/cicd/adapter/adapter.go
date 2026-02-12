package adapter

import (
	"fmt"

	"ci-cd-orchestrator/internal/cicd/common"
)

// Platform 平台类型
type Platform = common.Platform

// PipelineConfig CI/CD 管道配置
type PipelineConfig = common.PipelineConfig

// Adapter 平台适配器接口
type Adapter interface {
	ConvertToPlatform(config *PipelineConfig) (*PipelineConfig, error)
	ConvertFromPlatform(config *PipelineConfig) (*PipelineConfig, error)
	GetPlatform() Platform
}

// GitHubActionsAdapter GitHub Actions 适配器
type GitHubActionsAdapter struct{}

// NewGitHubActionsAdapter 创建 GitHub Actions 适配器实例
func NewGitHubActionsAdapter() Adapter {
	return &GitHubActionsAdapter{}
}

// ConvertToPlatform 转换为 GitHub Actions 配置
func (a *GitHubActionsAdapter) ConvertToPlatform(config *PipelineConfig) (*PipelineConfig, error) {
	// 如果已经是 GitHub Actions 配置，直接返回
	if config.Platform == common.PlatformGitHubActions {
		return config, nil
	}

	// 这里可以实现从其他平台转换为 GitHub Actions 配置的逻辑
	// 现在暂时返回一个新的 GitHub Actions 配置
	return &PipelineConfig{
		Platform:   common.PlatformGitHubActions,
		ConfigType: common.ConfigTypeYAML,
		Content:    config.Content, // 这里应该转换内容
		Filename:   ".github/workflows/ci.yml",
	}, nil
}

// ConvertFromPlatform 从 GitHub Actions 配置转换
func (a *GitHubActionsAdapter) ConvertFromPlatform(config *PipelineConfig) (*PipelineConfig, error) {
	// 检查配置是否是 GitHub Actions 配置
	if config.Platform != common.PlatformGitHubActions {
		return nil, fmt.Errorf("config is not a GitHub Actions config")
	}

	// 这里可以实现从 GitHub Actions 转换为其他平台配置的逻辑
	// 现在暂时返回原配置
	return config, nil
}

// GetPlatform 获取平台类型
func (a *GitHubActionsAdapter) GetPlatform() Platform {
	return common.PlatformGitHubActions
}

// MockAdapter Mock 平台适配器
type MockAdapter struct{}

// NewMockAdapter 创建 Mock 平台适配器实例
func NewMockAdapter() Adapter {
	return &MockAdapter{}
}

// ConvertToPlatform 转换为 Mock 平台配置
func (a *MockAdapter) ConvertToPlatform(config *PipelineConfig) (*PipelineConfig, error) {
	// 如果已经是 Mock 平台配置，直接返回
	if config.Platform == common.PlatformMock {
		return config, nil
	}

	// 这里可以实现从其他平台转换为 Mock 平台配置的逻辑
	// 现在暂时返回一个新的 Mock 平台配置
	return &PipelineConfig{
		Platform:   common.PlatformMock,
		ConfigType: common.ConfigTypeYAML,
		Content:    config.Content, // 这里应该转换内容
		Filename:   "mock-ci.yml",
	}, nil
}

// ConvertFromPlatform 从 Mock 平台配置转换
func (a *MockAdapter) ConvertFromPlatform(config *PipelineConfig) (*PipelineConfig, error) {
	// 检查配置是否是 Mock 平台配置
	if config.Platform != common.PlatformMock {
		return nil, fmt.Errorf("config is not a Mock config")
	}

	// 这里可以实现从 Mock 平台转换为其他平台配置的逻辑
	// 现在暂时返回原配置
	return config, nil
}

// GetPlatform 获取平台类型
func (a *MockAdapter) GetPlatform() Platform {
	return common.PlatformMock
}

// NewAdapter 根据平台类型创建适配器实例
func NewAdapter(platform Platform) (Adapter, error) {
	switch platform {
	case common.PlatformGitHubActions:
		return NewGitHubActionsAdapter(), nil
	case common.PlatformMock:
		return NewMockAdapter(), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
