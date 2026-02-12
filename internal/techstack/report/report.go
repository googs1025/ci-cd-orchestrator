package report

import (
	"encoding/json"
	"fmt"

	"ci-cd-orchestrator/internal/techstack"
)

// Reporter 报告生成器接口
type Reporter interface {
	Generate(result *techstack.Result) (string, error)
	GenerateJSON(result *techstack.Result) (string, error)
	GenerateMarkdown(result *techstack.Result) (string, error)
}

// ReportGenerator 报告生成器实现
type ReportGenerator struct{}

// NewReportGenerator 创建报告生成器实例
func NewReportGenerator() Reporter {
	return &ReportGenerator{}
}

// Generate 生成报告
func (r *ReportGenerator) Generate(result *techstack.Result) (string, error) {
	// 默认生成 JSON 格式报告
	return r.GenerateJSON(result)
}

// GenerateJSON 生成 JSON 格式报告
func (r *ReportGenerator) GenerateJSON(result *techstack.Result) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GenerateMarkdown 生成 Markdown 格式报告
func (r *ReportGenerator) GenerateMarkdown(result *techstack.Result) (string, error) {
	var markdown string

	markdown += fmt.Sprintf("# 技术栈识别报告\n\n")
	markdown += fmt.Sprintf("## 项目信息\n\n")
	markdown += fmt.Sprintf("- 项目路径: %s\n", result.ProjectPath)
	markdown += fmt.Sprintf("- 识别置信度: %.2f%%\n\n", result.Confidence*100)

	markdown += fmt.Sprintf("## 技术栈信息\n\n")
	markdown += fmt.Sprintf("- 编程语言: %s\n", result.TechStack.Language)
	markdown += fmt.Sprintf("- 框架: %s\n", result.TechStack.Framework)
	markdown += fmt.Sprintf("- 构建工具: %s\n", result.TechStack.BuildTool)
	markdown += fmt.Sprintf("- 测试框架: %s\n\n", result.TechStack.TestFramework)

	if len(result.TechStack.Dependencies) > 0 {
		markdown += fmt.Sprintf("## 依赖项\n\n")
		markdown += fmt.Sprintf("| 依赖名称 | 版本 |\n")
		markdown += fmt.Sprintf("|---------|------|\n")

		for dep, version := range result.TechStack.Dependencies {
			if version == "" {
				version = "latest"
			}
			markdown += fmt.Sprintf("| %s | %s |\n", dep, version)
		}
		markdown += fmt.Sprintf("\n")
	}

	if len(result.TechStack.Files) > 0 {
		markdown += fmt.Sprintf("## 相关文件\n\n")
		for _, file := range result.TechStack.Files {
			markdown += fmt.Sprintf("- %s\n", file)
		}
		markdown += fmt.Sprintf("\n")
	}

	if len(result.Errors) > 0 {
		markdown += fmt.Sprintf("## 错误信息\n\n")
		for _, err := range result.Errors {
			markdown += fmt.Sprintf("- %s\n", err)
		}
	}

	return markdown, nil
}
