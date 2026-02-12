package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Analyzer 依赖分析器接口
type Analyzer interface {
	Analyze(files []string) (map[string]string, error)
}

// DependencyAnalyzer 依赖分析器实现
type DependencyAnalyzer struct{}

// NewDependencyAnalyzer 创建依赖分析器实例
func NewDependencyAnalyzer() Analyzer {
	return &DependencyAnalyzer{}
}

// Analyze 分析依赖
func (a *DependencyAnalyzer) Analyze(files []string) (map[string]string, error) {
	dependencies := make(map[string]string)

	for _, file := range files {
		fileName := strings.ToLower(filepath.Base(file))

		switch fileName {
		case "package.json":
			deps, err := a.analyzePackageJSON(file)
			if err == nil {
				for k, v := range deps {
					dependencies[k] = v
				}
			}
		case "go.mod":
			deps, err := a.analyzeGoMod(file)
			if err == nil {
				for k, v := range deps {
					dependencies[k] = v
				}
			}
		case "requirements.txt":
			deps, err := a.analyzeRequirementsTxt(file)
			if err == nil {
				for k, v := range deps {
					dependencies[k] = v
				}
			}
		case "pom.xml":
			deps, err := a.analyzePomXML(file)
			if err == nil {
				for k, v := range deps {
					dependencies[k] = v
				}
			}
		case "cargo.toml":
			deps, err := a.analyzeCargoToml(file)
			if err == nil {
				for k, v := range deps {
					dependencies[k] = v
				}
			}
		}
	}

	return dependencies, nil
}

// analyzePackageJSON 分析 package.json 文件
func (a *DependencyAnalyzer) analyzePackageJSON(filePath string) (map[string]string, error) {
	dependencies := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies, err
	}

	// 解析 JSON
	var data struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return dependencies, err
	}

	// 合并依赖
	for k, v := range data.Dependencies {
		dependencies[k] = v
	}
	for k, v := range data.DevDependencies {
		dependencies[k] = v
	}

	return dependencies, nil
}

// analyzeGoMod 分析 go.mod 文件
func (a *DependencyAnalyzer) analyzeGoMod(filePath string) (map[string]string, error) {
	dependencies := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies, err
	}

	// 解析 go.mod 文件
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require (") {
			// 开始 require 块
			continue
		}
		if strings.HasPrefix(line, ")") {
			// 结束 require 块
			continue
		}
		if strings.HasPrefix(line, "require ") {
			// 单行 require
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[1]
				version := parts[2]
				dependencies[name] = version
			}
		}
		if !strings.HasPrefix(line, "//") && strings.Contains(line, " ") {
			// 可能是依赖项
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				version := parts[1]
				dependencies[name] = version
			}
		}
	}

	return dependencies, nil
}

// analyzeRequirementsTxt 分析 requirements.txt 文件
func (a *DependencyAnalyzer) analyzeRequirementsTxt(filePath string) (map[string]string, error) {
	dependencies := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies, err
	}

	// 解析 requirements.txt 文件
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 处理依赖项
		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			version := strings.TrimSpace(parts[1])
			dependencies[name] = version
		} else {
			// 没有版本号
			name := strings.TrimSpace(line)
			dependencies[name] = ""
		}
	}

	return dependencies, nil
}

// analyzePomXML 分析 pom.xml 文件
func (a *DependencyAnalyzer) analyzePomXML(filePath string) (map[string]string, error) {
	dependencies := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies, err
	}

	// 简单解析 XML（实际项目中应该使用 XML 解析库）
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	inDependency := false
	var groupID, artifactID, version string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "<dependency>") {
			inDependency = true
			groupID = ""
			artifactID = ""
			version = ""
		}

		if inDependency {
			if strings.Contains(line, "<groupId>") {
				groupID = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "<groupId>", ""), "</groupId>", ""))
			}
			if strings.Contains(line, "<artifactId>") {
				artifactID = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "<artifactId>", ""), "</artifactId>", ""))
			}
			if strings.Contains(line, "<version>") {
				version = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "<version>", ""), "</version>", ""))
			}
		}

		if strings.Contains(line, "</dependency>") {
			inDependency = false
			if groupID != "" && artifactID != "" {
				name := groupID + ":" + artifactID
				dependencies[name] = version
			}
		}
	}

	return dependencies, nil
}

// analyzeCargoToml 分析 Cargo.toml 文件
func (a *DependencyAnalyzer) analyzeCargoToml(filePath string) (map[string]string, error) {
	dependencies := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies, err
	}

	// 简单解析 Cargo.toml（实际项目中应该使用 TOML 解析库）
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	inDependencies := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "[dependencies]") {
			inDependencies = true
			continue
		}

		if strings.HasPrefix(line, "[") && !strings.Contains(line, "dependencies") {
			inDependencies = false
			continue
		}

		if inDependencies && line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[0])
				version := strings.TrimSpace(parts[1])
				// 移除引号
				version = strings.Trim(version, `"'`)
				dependencies[name] = version
			}
		}
	}

	return dependencies, nil
}
