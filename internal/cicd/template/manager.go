package template

import (
	"fmt"

	"ci-cd-orchestrator/internal/cicd/common"
	"ci-cd-orchestrator/internal/repository"
	"ci-cd-orchestrator/internal/techstack"
)

// Template 配置模板
type Template struct {
	Platform   common.Platform   `json:"platform"`
	Language   string            `json:"language"`
	Framework  string            `json:"framework"`
	Content    string            `json:"content"`
	Filename   string            `json:"filename"`
	ConfigType common.ConfigType `json:"config_type"`
}

// Manager 模板管理器接口
type Manager interface {
	GetTemplate(techStack *techstack.TechStack, platform common.Platform) (*Template, error)
	GetTemplateByID(id int) (*Template, error)
	AddTemplate(template *Template) error
	ListTemplates() []*Template
}

// TemplateManager 模板管理器实现
type TemplateManager struct {
	templates    []*Template
	templateRepo repository.TemplateRepository
}

// NewTemplateManager 创建模板管理器实例
func NewTemplateManager(templateRepo repository.TemplateRepository) Manager {
	manager := &TemplateManager{
		templates:    []*Template{},
		templateRepo: templateRepo,
	}

	// 初始化默认模板
	manager.initDefaultTemplates()

	// 从数据库加载模板
	manager.loadTemplatesFromDB()

	return manager
}

// GetTemplate 获取适合的模板
func (m *TemplateManager) GetTemplate(techStack *techstack.TechStack, platform common.Platform) (*Template, error) {
	// 按优先级查找模板：
	// 1. 语言 + 框架 + 平台
	// 2. 语言 + 平台
	// 3. 平台默认模板

	// 查找语言 + 框架 + 平台的模板
	for _, template := range m.templates {
		if template.Platform == platform &&
			template.Language == techStack.Language &&
			template.Framework == techStack.Framework {
			return template, nil
		}
	}

	// 查找语言 + 平台的模板
	for _, template := range m.templates {
		if template.Platform == platform &&
			template.Language == techStack.Language &&
			template.Framework == "" {
			return template, nil
		}
	}

	// 查找平台默认模板
	for _, template := range m.templates {
		if template.Platform == platform &&
			template.Language == "" &&
			template.Framework == "" {
			return template, nil
		}
	}

	// 从数据库中查找模板
	if m.templateRepo != nil {
		dbTemplate, err := m.templateRepo.GetByPlatformAndLanguage(string(platform), techStack.Language, techStack.Framework)
		if err == nil {
			// 转换为Template类型
			template := &Template{
				Platform:   common.Platform(dbTemplate.Platform),
				Language:   dbTemplate.Language,
				Framework:  dbTemplate.Framework,
				Content:    dbTemplate.Content,
				Filename:   dbTemplate.Filename,
				ConfigType: common.ConfigType(dbTemplate.ConfigType),
			}
			// 添加到内存中
			m.templates = append(m.templates, template)
			return template, nil
		}
	}

	return nil, fmt.Errorf("no template found for platform %s and language %s", platform, techStack.Language)
}

// loadTemplatesFromDB 从数据库加载模板
func (m *TemplateManager) loadTemplatesFromDB() {
	if m.templateRepo == nil {
		return
	}

	dbTemplates, err := m.templateRepo.GetAll()
	if err != nil {
		return
	}

	// 转换为Template类型并添加到内存中
	for _, dbTemplate := range dbTemplates {
		template := &Template{
			Platform:   common.Platform(dbTemplate.Platform),
			Language:   dbTemplate.Language,
			Framework:  dbTemplate.Framework,
			Content:    dbTemplate.Content,
			Filename:   dbTemplate.Filename,
			ConfigType: common.ConfigType(dbTemplate.ConfigType),
		}
		m.templates = append(m.templates, template)
	}
}

// AddTemplate 添加模板
func (m *TemplateManager) AddTemplate(template *Template) error {
	// 添加到内存中
	m.templates = append(m.templates, template)

	// 保存到数据库中
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(template.Platform),
			Language:   template.Language,
			Framework:  template.Framework,
			Content:    template.Content,
			Filename:   template.Filename,
			ConfigType: string(template.ConfigType),
			IsBuiltin:  false,
		}
		return m.templateRepo.Create(dbTemplate)
	}

	return nil
}

// ListTemplates 列出所有模板
func (m *TemplateManager) ListTemplates() []*Template {
	return m.templates
}

// GetTemplateByID 根据ID获取模板
func (m *TemplateManager) GetTemplateByID(id int) (*Template, error) {
	// 从数据库中查找
	if m.templateRepo != nil {
		dbTemplate, err := m.templateRepo.GetByID(id)
		if err == nil {
			// 转换为Template类型
			template := &Template{
				Platform:   common.Platform(dbTemplate.Platform),
				Language:   dbTemplate.Language,
				Framework:  dbTemplate.Framework,
				Content:    dbTemplate.Content,
				Filename:   dbTemplate.Filename,
				ConfigType: common.ConfigType(dbTemplate.ConfigType),
			}
			return template, nil
		}
	}

	return nil, fmt.Errorf("template not found with ID: %d", id)
}

// initDefaultTemplates 初始化默认模板
func (m *TemplateManager) initDefaultTemplates() {
	// 检查数据库中是否已经有内置模板
	if m.templateRepo != nil {
		allTemplates, err := m.templateRepo.GetAll()
		if err == nil {
			hasBuiltin := false
			for _, template := range allTemplates {
				if template.IsBuiltin {
					hasBuiltin = true
					break
				}
			}
			if hasBuiltin {
				// 已经有内置模板，不需要初始化
				return
			}
		}
	}

	// GitHub Actions 模板
	m.initGitHubActionsTemplates()

	// Mock 模板
	m.initMockTemplates()
}

// initGitHubActionsTemplates 初始化 GitHub Actions 模板
func (m *TemplateManager) initGitHubActionsTemplates() {
	// Go 项目模板
	goTemplate := &Template{
		Platform:   common.PlatformGitHubActions,
		Language:   "Go",
		Framework:  "",
		Content:    getGoGitHubActionsTemplate(),
		Filename:   ".github/workflows/ci.yml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, goTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(goTemplate.Platform),
			Language:   goTemplate.Language,
			Framework:  goTemplate.Framework,
			Content:    goTemplate.Content,
			Filename:   goTemplate.Filename,
			ConfigType: string(goTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}

	// Java 项目模板
	javaTemplate := &Template{
		Platform:   common.PlatformGitHubActions,
		Language:   "Java",
		Framework:  "",
		Content:    getJavaGitHubActionsTemplate(),
		Filename:   ".github/workflows/ci.yml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, javaTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(javaTemplate.Platform),
			Language:   javaTemplate.Language,
			Framework:  javaTemplate.Framework,
			Content:    javaTemplate.Content,
			Filename:   javaTemplate.Filename,
			ConfigType: string(javaTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}

	// Python 项目模板
	pythonTemplate := &Template{
		Platform:   common.PlatformGitHubActions,
		Language:   "Python",
		Framework:  "",
		Content:    getPythonGitHubActionsTemplate(),
		Filename:   ".github/workflows/ci.yml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, pythonTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(pythonTemplate.Platform),
			Language:   pythonTemplate.Language,
			Framework:  pythonTemplate.Framework,
			Content:    pythonTemplate.Content,
			Filename:   pythonTemplate.Filename,
			ConfigType: string(pythonTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}

	// JavaScript 项目模板
	jsTemplate := &Template{
		Platform:   common.PlatformGitHubActions,
		Language:   "JavaScript",
		Framework:  "",
		Content:    getJavaScriptGitHubActionsTemplate(),
		Filename:   ".github/workflows/ci.yml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, jsTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(jsTemplate.Platform),
			Language:   jsTemplate.Language,
			Framework:  jsTemplate.Framework,
			Content:    jsTemplate.Content,
			Filename:   jsTemplate.Filename,
			ConfigType: string(jsTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}

	// GitHub Actions 默认模板
	defaultTemplate := &Template{
		Platform:   common.PlatformGitHubActions,
		Language:   "",
		Framework:  "",
		Content:    getDefaultGitHubActionsTemplate(),
		Filename:   ".github/workflows/ci.yml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, defaultTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(defaultTemplate.Platform),
			Language:   defaultTemplate.Language,
			Framework:  defaultTemplate.Framework,
			Content:    defaultTemplate.Content,
			Filename:   defaultTemplate.Filename,
			ConfigType: string(defaultTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}
}

// initMockTemplates 初始化 Mock 模板
func (m *TemplateManager) initMockTemplates() {
	// Mock 平台默认模板
	mockTemplate := &Template{
		Platform:   common.PlatformMock,
		Language:   "",
		Framework:  "",
		Content:    getMockTemplate(),
		Filename:   ".mock/workflows/ci.yaml",
		ConfigType: common.ConfigTypeYAML,
	}
	m.templates = append(m.templates, mockTemplate)

	// 保存到数据库
	if m.templateRepo != nil {
		dbTemplate := &repository.Template{
			Platform:   string(mockTemplate.Platform),
			Language:   mockTemplate.Language,
			Framework:  mockTemplate.Framework,
			Content:    mockTemplate.Content,
			Filename:   mockTemplate.Filename,
			ConfigType: string(mockTemplate.ConfigType),
			IsBuiltin:  true,
		}
		m.templateRepo.Create(dbTemplate)
	}
}

// getGoGitHubActionsTemplate 获取 Go 项目的 GitHub Actions 模板
func getGoGitHubActionsTemplate() string {
	return `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    # 资源配置
    resources:
      limits:
        cpu: 2
        memory: 4G
      requests:
        cpu: 1
        memory: 2G
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20
    - name: Build
      run: go build -v ./...
    - name: Test
      run: go test -v ./...
`
}

// getJavaGitHubActionsTemplate 获取 Java 项目的 GitHub Actions 模板
func getJavaGitHubActionsTemplate() string {
	return `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up JDK 11
      uses: actions/setup-java@v2
      with:
        java-version: '11'
        distribution: 'adopt'
    - name: Build with Maven
      run: mvn -B package --file pom.xml
    - name: Test
      run: mvn test
`
}

// getPythonGitHubActionsTemplate 获取 Python 项目的 GitHub Actions 模板
func getPythonGitHubActionsTemplate() string {
	return `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: 3.9
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
    - name: Test with pytest
      run: pytest
`
}

// getJavaScriptGitHubActionsTemplate 获取 JavaScript 项目的 GitHub Actions 模板
func getJavaScriptGitHubActionsTemplate() string {
	return `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Node.js
      uses: actions/setup-node@v2
      with:
        node-version: 16
    - name: Install dependencies
      run: npm install
    - name: Build
      run: npm run build --if-present
    - name: Test
      run: npm test
`
}

// getDefaultGitHubActionsTemplate 获取默认的 GitHub Actions 模板
func getDefaultGitHubActionsTemplate() string {
	return `name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Build and test
      run: |
        echo "Building and testing project..."
        # Add your build and test commands here
`
}

// getMockTemplate 获取 Mock 平台的模板
func getMockTemplate() string {
	return `name: Mock CI

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: mock-runner
    # 资源配置
    resources:
      limits:
        cpu: 2
        memory: 4G
      requests:
        cpu: 1
        memory: 2G
    steps:
    - name: Mock checkout
      run: echo "Mock checkout step"
    - name: Mock build
      run: echo "Mock build step"
    - name: Mock test
      run: echo "Mock test step"
`
}
