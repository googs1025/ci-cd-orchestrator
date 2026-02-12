package detector

import (
	"path/filepath"
	"strings"
)

// Detector 技术栈检测器接口
type Detector interface {
	Detect(files []string) (string, string, string, string, error)
}

// TechStackDetector 技术栈检测器实现
type TechStackDetector struct{}

// NewTechStackDetector 创建技术栈检测器实例
func NewTechStackDetector() Detector {
	return &TechStackDetector{}
}

// Detect 检测技术栈
func (d *TechStackDetector) Detect(files []string) (language, framework, buildTool, testFramework string, err error) {
	// 统计文件类型
	fileCount := make(map[string]int)
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		fileName := strings.ToLower(filepath.Base(file))

		// 统计文件扩展名
		if ext != "" {
			fileCount[ext]++
		}

		// 统计关键文件名
		fileCount[fileName]++
	}

	// 检测编程语言
	language = d.detectLanguage(fileCount)

	// 检测框架
	framework = d.detectFramework(fileCount, language)

	// 检测构建工具
	buildTool = d.detectBuildTool(fileCount, language)

	// 检测测试框架
	testFramework = d.detectTestFramework(fileCount, language)

	return
}

// detectLanguage 检测编程语言
func (d *TechStackDetector) detectLanguage(fileCount map[string]int) string {
	// 检查关键文件
	if fileCount["go.mod"] > 0 {
		return "Go"
	}
	if fileCount["package.json"] > 0 {
		return "JavaScript"
	}
	if fileCount["pom.xml"] > 0 || fileCount["build.gradle"] > 0 {
		return "Java"
	}
	if fileCount["requirements.txt"] > 0 || fileCount["setup.py"] > 0 {
		return "Python"
	}
	if fileCount["Cargo.toml"] > 0 {
		return "Rust"
	}
	if fileCount[".csproj"] > 0 {
		return "C#"
	}

	// 检查文件扩展名
	if fileCount[".go"] > 0 {
		return "Go"
	}
	if fileCount[".js"] > 0 || fileCount[".ts"] > 0 {
		return "JavaScript"
	}
	if fileCount[".java"] > 0 {
		return "Java"
	}
	if fileCount[".py"] > 0 {
		return "Python"
	}
	if fileCount[".rs"] > 0 {
		return "Rust"
	}
	if fileCount[".cpp"] > 0 || fileCount[".c"] > 0 {
		return "C++"
	}
	if fileCount[".cs"] > 0 {
		return "C#"
	}

	return "Unknown"
}

// detectFramework 检测框架
func (d *TechStackDetector) detectFramework(fileCount map[string]int, language string) string {
	switch language {
	case "JavaScript":
		if fileCount["package.json"] > 0 {
			// 检查 package.json 内容
			if d.hasDependency("react", fileCount) {
				return "React"
			}
			if d.hasDependency("vue", fileCount) {
				return "Vue"
			}
			if d.hasDependency("angular", fileCount) {
				return "Angular"
			}
			if d.hasDependency("svelte", fileCount) {
				return "Svelte"
			}
			if d.hasDependency("express", fileCount) {
				return "Express"
			}
			if d.hasDependency("nestjs", fileCount) {
				return "NestJS"
			}
		}
	case "Java":
		if fileCount["pom.xml"] > 0 || fileCount["build.gradle"] > 0 {
			if d.hasDependency("spring", fileCount) {
				return "Spring"
			}
		}
	case "Python":
		if fileCount["requirements.txt"] > 0 {
			if d.hasDependency("django", fileCount) {
				return "Django"
			}
			if d.hasDependency("fastapi", fileCount) {
				return "FastAPI"
			}
		}
	case "Go":
		if fileCount["go.mod"] > 0 {
			if d.hasDependency("gin", fileCount) {
				return "Gin"
			}
		}
	}

	return ""
}

// detectBuildTool 检测构建工具
func (d *TechStackDetector) detectBuildTool(fileCount map[string]int, language string) string {
	switch language {
	case "JavaScript":
		if fileCount["package.json"] > 0 {
			return "npm"
		}
	case "Java":
		if fileCount["pom.xml"] > 0 {
			return "Maven"
		}
		if fileCount["build.gradle"] > 0 {
			return "Gradle"
		}
	case "Python":
		if fileCount["requirements.txt"] > 0 {
			return "pip"
		}
	case "Go":
		if fileCount["go.mod"] > 0 {
			return "go mod"
		}
	case "Rust":
		if fileCount["Cargo.toml"] > 0 {
			return "cargo"
		}
	}

	return ""
}

// detectTestFramework 检测测试框架
func (d *TechStackDetector) detectTestFramework(fileCount map[string]int, language string) string {
	switch language {
	case "JavaScript":
		if d.hasDependency("jest", fileCount) {
			return "Jest"
		}
		if d.hasDependency("mocha", fileCount) {
			return "Mocha"
		}
		if d.hasDependency("vitest", fileCount) {
			return "Vitest"
		}
	case "Java":
		if d.hasDependency("junit", fileCount) {
			return "JUnit"
		}
	case "Python":
		if d.hasDependency("pytest", fileCount) {
			return "pytest"
		}
	case "Go":
		// Go 标准库包含测试框架
		if fileCount[".go"] > 0 {
			return "testing"
		}
	}

	return ""
}

// hasDependency 检查是否有特定依赖
func (d *TechStackDetector) hasDependency(depName string, fileCount map[string]int) bool {
	// 这里简化处理，实际应该读取 package.json 等文件的内容
	// 现在只是基于文件名和扩展名进行判断
	for fileName := range fileCount {
		if strings.Contains(strings.ToLower(fileName), depName) {
			return true
		}
	}

	// 尝试读取 package.json 文件内容
	if fileCount["package.json"] > 0 {
		// 这里应该找到实际的 package.json 文件路径并读取内容
		// 现在暂时返回 false
	}

	return false
}
