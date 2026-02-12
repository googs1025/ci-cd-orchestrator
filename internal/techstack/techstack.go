package techstack

import (
	"ci-cd-orchestrator/internal/techstack/analyzer"
	"ci-cd-orchestrator/internal/techstack/detector"
	"ci-cd-orchestrator/internal/techstack/scanner"
)

// TechStack 技术栈信息
type TechStack struct {
	Language      string            `json:"language"`
	Framework     string            `json:"framework"`
	BuildTool     string            `json:"build_tool"`
	TestFramework string            `json:"test_framework"`
	Dependencies  map[string]string `json:"dependencies"`
	Files         []string          `json:"files"`
}

// Result 技术栈识别结果
type Result struct {
	ProjectPath string    `json:"project_path"`
	TechStack   TechStack `json:"tech_stack"`
	Confidence  float64   `json:"confidence"`
	Errors      []string  `json:"errors"`
}

// Recognizer 技术栈识别器接口
type Recognizer interface {
	Recognize(projectPath string) (*Result, error)
}

// NewRecognizer 创建技术栈识别器实例
func NewRecognizer() Recognizer {
	return &recognizerImpl{}
}

// recognizerImpl 技术栈识别器实现
type recognizerImpl struct{}

// Recognize 识别项目技术栈
func (r *recognizerImpl) Recognize(projectPath string) (*Result, error) {
	result := &Result{
		ProjectPath: projectPath,
		TechStack: TechStack{
			Language:      "",
			Framework:     "",
			BuildTool:     "",
			TestFramework: "",
			Dependencies:  make(map[string]string),
			Files:         []string{},
		},
		Confidence: 0.0,
		Errors:     []string{},
	}

	// 1. 文件扫描
	scanner := scanner.NewFileScanner()
	files, err := scanner.Scan(projectPath)
	if err != nil {
		result.Errors = append(result.Errors, "文件扫描失败: "+err.Error())
		return result, err
	}
	result.TechStack.Files = files

	// 2. 技术栈检测
	detector := detector.NewTechStackDetector()
	language, framework, buildTool, testFramework, err := detector.Detect(files)
	if err != nil {
		result.Errors = append(result.Errors, "技术栈检测失败: "+err.Error())
	}
	result.TechStack.Language = language
	result.TechStack.Framework = framework
	result.TechStack.BuildTool = buildTool
	result.TechStack.TestFramework = testFramework

	// 3. 依赖分析
	analyzer := analyzer.NewDependencyAnalyzer()
	dependencies, err := analyzer.Analyze(files)
	if err != nil {
		result.Errors = append(result.Errors, "依赖分析失败: "+err.Error())
	} else {
		result.TechStack.Dependencies = dependencies
	}

	// 4. 计算置信度
	result.Confidence = r.calculateConfidence(result)

	return result, nil
}

// calculateConfidence 计算识别置信度
func (r *recognizerImpl) calculateConfidence(result *Result) float64 {
	confidence := 0.0

	// 根据识别到的信息计算置信度
	if result.TechStack.Language != "" && result.TechStack.Language != "Unknown" {
		confidence += 0.3
	}

	if result.TechStack.Framework != "" {
		confidence += 0.2
	}

	if result.TechStack.BuildTool != "" {
		confidence += 0.2
	}

	if result.TechStack.TestFramework != "" {
		confidence += 0.1
	}

	if len(result.TechStack.Dependencies) > 0 {
		confidence += 0.1
	}

	if len(result.TechStack.Files) > 0 {
		confidence += 0.1
	}

	return confidence
}
