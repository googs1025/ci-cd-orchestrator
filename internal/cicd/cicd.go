package cicd

import (
	"ci-cd-orchestrator/internal/cicd/common"
	"ci-cd-orchestrator/internal/cicd/template"
	"ci-cd-orchestrator/internal/cicd/validator"
	"ci-cd-orchestrator/internal/repository"
	"ci-cd-orchestrator/internal/techstack"
)

// Platform CI/CD 平台类型
type Platform = common.Platform

// 支持的 CI/CD 平台
const (
	PlatformGitHubActions Platform = common.PlatformGitHubActions
	PlatformMock          Platform = common.PlatformMock
)

// ConfigType 配置类型
type ConfigType = common.ConfigType

// 支持的配置类型
const (
	ConfigTypeYAML ConfigType = common.ConfigTypeYAML
	ConfigTypeJSON ConfigType = common.ConfigTypeJSON
)

// PipelineConfig CI/CD 管道配置
type PipelineConfig = common.PipelineConfig

// Generator CI/CD 管道配置生成器接口
type Generator interface {
	GenerateConfig(techStack *techstack.TechStack, platform Platform, templateID ...int) (*PipelineConfig, error)
	ValidateConfig(config *PipelineConfig) error
}

// generatorImpl CI/CD 管道配置生成器实现
type generatorImpl struct {
	templateManager template.Manager
	validator       validator.Validator
}

// NewGenerator 创建 CI/CD 管道配置生成器实例
func NewGenerator(templateRepo repository.TemplateRepository) Generator {
	return &generatorImpl{
		templateManager: template.NewTemplateManager(templateRepo),
		validator:       validator.NewValidator(),
	}
}

// GenerateConfig 生成 CI/CD 管道配置
func (g *generatorImpl) GenerateConfig(techStack *techstack.TechStack, platform Platform, templateID ...int) (*PipelineConfig, error) {
	var tmpl *template.Template
	var err error

	// 如果提供了 templateID，直接使用该模板
	if len(templateID) > 0 && templateID[0] > 0 {
		tmpl, err = g.templateManager.GetTemplateByID(templateID[0])
		if err != nil {
			return nil, err
		}
	} else {
		// 否则根据技术栈和平台选择模板
		tmpl, err = g.templateManager.GetTemplate(techStack, platform)
		if err != nil {
			return nil, err
		}
	}

	// 生成配置
	config := &PipelineConfig{
		Platform:   platform,
		ConfigType: tmpl.ConfigType,
		Content:    tmpl.Content,
		Filename:   tmpl.Filename,
	}

	// 验证配置
	if err := g.validator.Validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfig 验证 CI/CD 管道配置
func (g *generatorImpl) ValidateConfig(config *PipelineConfig) error {
	return g.validator.Validate(config)
}
