package common

// Platform CI/CD 平台类型
type Platform string

// 支持的 CI/CD 平台
const (
	PlatformGitHubActions Platform = "github_actions"
	PlatformMock          Platform = "mock"
)

// ConfigType 配置类型
type ConfigType string

// 支持的配置类型
const (
	ConfigTypeYAML ConfigType = "yaml"
	ConfigTypeJSON ConfigType = "json"
)

// PipelineConfig CI/CD 管道配置
type PipelineConfig struct {
	Platform   Platform   `json:"platform"`
	ConfigType ConfigType `json:"config_type"`
	Content    string     `json:"content"`
	Filename   string     `json:"filename"`
}
