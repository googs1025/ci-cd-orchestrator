package repository

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"ci-cd-orchestrator/internal/db"
	"ci-cd-orchestrator/internal/models"
	"ci-cd-orchestrator/pkg/migration"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// 初始化数据库
	err := db.Init()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	testDB = db.GetDB()

	// 执行迁移
	err = migration.Run(testDB)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 运行测试
	m.Run()

	// 关闭数据库
	db.Close()
}

func TestProjectRepository(t *testing.T) {
	repo := NewProjectRepository(testDB)

	// 测试创建项目
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := repo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	if project.ID == 0 {
		t.Fatal("项目 ID 未设置")
	}

	// 测试获取项目
	getProject, err := repo.GetByID(project.ID)
	if err != nil {
		t.Fatalf("获取项目失败: %v", err)
	}

	if getProject.Name != project.Name {
		t.Errorf("项目名称不匹配: 期望 %s, 实际 %s", project.Name, getProject.Name)
	}

	// 测试获取所有项目
	projects, err := repo.GetAll()
	if err != nil {
		t.Fatalf("获取所有项目失败: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("项目列表为空")
	}

	// 测试更新项目
	project.Name = "更新后的测试项目"
	err = repo.Update(project)
	if err != nil {
		t.Fatalf("更新项目失败: %v", err)
	}

	updatedProject, err := repo.GetByID(project.ID)
	if err != nil {
		t.Fatalf("获取更新后的项目失败: %v", err)
	}

	if updatedProject.Name != project.Name {
		t.Errorf("更新后的项目名称不匹配: 期望 %s, 实际 %s", project.Name, updatedProject.Name)
	}

	// 测试删除项目
	err = repo.Delete(project.ID)
	if err != nil {
		t.Fatalf("删除项目失败: %v", err)
	}

	_, err = repo.GetByID(project.ID)
	if err == nil {
		t.Fatal("项目删除后仍能获取到")
	}
}

func TestTechStackRepository(t *testing.T) {
	// 先创建一个项目
	projectRepo := NewProjectRepository(testDB)
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := projectRepo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	testDB = db.GetDB()
	repo := NewTechStackRepository(testDB)

	// 测试创建技术栈
	techStack := &models.TechStack{
		ProjectID:     project.ID,
		Language:      "Go",
		Framework:     "",
		BuildTool:     "go build",
		TestFramework: "testing",
		Dependencies:  `{"github.com/mattn/go-sqlite3": "v1.14.34"}`,
	}

	err = repo.Create(techStack)
	if err != nil {
		t.Fatalf("创建技术栈失败: %v", err)
	}

	if techStack.ID == 0 {
		t.Fatal("技术栈 ID 未设置")
	}

	// 测试获取技术栈
	getTechStack, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取技术栈失败: %v", err)
	}

	if getTechStack.Language != techStack.Language {
		t.Errorf("技术栈语言不匹配: 期望 %s, 实际 %s", techStack.Language, getTechStack.Language)
	}

	// 测试更新技术栈
	techStack.Language = "Golang"
	err = repo.Update(techStack)
	if err != nil {
		t.Fatalf("更新技术栈失败: %v", err)
	}

	updatedTechStack, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取更新后的技术栈失败: %v", err)
	}

	if updatedTechStack.Language != techStack.Language {
		t.Errorf("更新后的技术栈语言不匹配: 期望 %s, 实际 %s", techStack.Language, updatedTechStack.Language)
	}

	// 测试删除技术栈
	err = repo.DeleteByProjectID(project.ID)
	if err != nil {
		t.Fatalf("删除技术栈失败: %v", err)
	}

	_, err = repo.GetByProjectID(project.ID)
	if err == nil {
		t.Fatal("技术栈删除后仍能获取到")
	}

	// 清理项目
	projectRepo.Delete(project.ID)
}

func TestPipelineRepository(t *testing.T) {
	// 先创建一个项目
	projectRepo := NewProjectRepository(testDB)
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := projectRepo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	testDB = db.GetDB()
	repo := NewPipelineRepository(testDB)

	// 测试创建管道配置
	pipeline := &models.Pipeline{
		ProjectID: project.ID,
		Platform:  "GitHub Actions",
		Config:    `name: Test CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v2\n      - run: go test ./...`,
	}

	err = repo.Create(pipeline)
	if err != nil {
		t.Fatalf("创建管道配置失败: %v", err)
	}

	if pipeline.ID == 0 {
		t.Fatal("管道配置 ID 未设置")
	}

	// 测试根据项目 ID 获取管道配置
	pipelines, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取管道配置失败: %v", err)
	}

	if len(pipelines) == 0 {
		t.Fatal("管道配置列表为空")
	}

	// 测试根据 ID 获取管道配置
	getPipeline, err := repo.GetByID(pipeline.ID)
	if err != nil {
		t.Fatalf("获取管道配置失败: %v", err)
	}

	if getPipeline.Platform != pipeline.Platform {
		t.Errorf("管道配置平台不匹配: 期望 %s, 实际 %s", pipeline.Platform, getPipeline.Platform)
	}

	// 测试更新管道配置
	pipeline.Config = `name: Updated Test CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v2\n      - run: go test ./...\n      - run: go build ./...`
	err = repo.Update(pipeline)
	if err != nil {
		t.Fatalf("更新管道配置失败: %v", err)
	}

	updatedPipeline, err := repo.GetByID(pipeline.ID)
	if err != nil {
		t.Fatalf("获取更新后的管道配置失败: %v", err)
	}

	if updatedPipeline.Config != pipeline.Config {
		t.Error("更新后的管道配置不匹配")
	}

	// 测试删除管道配置
	err = repo.Delete(pipeline.ID)
	if err != nil {
		t.Fatalf("删除管道配置失败: %v", err)
	}

	_, err = repo.GetByID(pipeline.ID)
	if err == nil {
		t.Fatal("管道配置删除后仍能获取到")
	}

	// 清理项目
	projectRepo.Delete(project.ID)
}

func TestExecutionRepository(t *testing.T) {
	// 先创建一个项目和管道配置
	projectRepo := NewProjectRepository(testDB)
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := projectRepo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	pipelineRepo := NewPipelineRepository(testDB)
	pipeline := &models.Pipeline{
		ProjectID: project.ID,
		Platform:  "GitHub Actions",
		Config:    `name: Test CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v2\n      - run: go test ./...`,
	}

	err = pipelineRepo.Create(pipeline)
	if err != nil {
		t.Fatalf("创建管道配置失败: %v", err)
	}

	testDB = db.GetDB()
	repo := NewExecutionRepository(testDB)

	// 测试创建执行历史
	now := time.Now()
	execution := &models.Execution{
		ProjectID:  project.ID,
		PipelineID: pipeline.ID,
		Status:     "success",
		StartTime:  now.Add(-time.Minute),
		EndTime:    now,
		Duration:   60,
		Logs:       "Test passed",
	}

	err = repo.Create(execution)
	if err != nil {
		t.Fatalf("创建执行历史失败: %v", err)
	}

	if execution.ID == 0 {
		t.Fatal("执行历史 ID 未设置")
	}

	// 测试根据项目 ID 获取执行历史
	executions, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取执行历史失败: %v", err)
	}

	if len(executions) == 0 {
		t.Fatal("执行历史列表为空")
	}

	// 测试根据 ID 获取执行历史
	getExecution, err := repo.GetByID(execution.ID)
	if err != nil {
		t.Fatalf("获取执行历史失败: %v", err)
	}

	if getExecution.Status != execution.Status {
		t.Errorf("执行历史状态不匹配: 期望 %s, 实际 %s", execution.Status, getExecution.Status)
	}

	// 测试更新执行状态
	newStatus := "failed"
	newEndTime := time.Now()
	newDuration := 120
	newLogs := "Test failed"

	err = repo.UpdateStatus(execution.ID, newStatus, newEndTime, newDuration, newLogs)
	if err != nil {
		t.Fatalf("更新执行状态失败: %v", err)
	}

	updatedExecution, err := repo.GetByID(execution.ID)
	if err != nil {
		t.Fatalf("获取更新后的执行历史失败: %v", err)
	}

	if updatedExecution.Status != newStatus {
		t.Errorf("更新后的执行状态不匹配: 期望 %s, 实际 %s", newStatus, updatedExecution.Status)
	}

	// 清理项目
	projectRepo.Delete(project.ID)
}

func TestMetricRepository(t *testing.T) {
	// 先创建一个项目、管道配置和执行历史
	projectRepo := NewProjectRepository(testDB)
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := projectRepo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	pipelineRepo := NewPipelineRepository(testDB)
	pipeline := &models.Pipeline{
		ProjectID: project.ID,
		Platform:  "GitHub Actions",
		Config:    `name: Test CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v2\n      - run: go test ./...`,
	}

	err = pipelineRepo.Create(pipeline)
	if err != nil {
		t.Fatalf("创建管道配置失败: %v", err)
	}

	executionRepo := NewExecutionRepository(testDB)
	now := time.Now()
	execution := &models.Execution{
		ProjectID:  project.ID,
		PipelineID: pipeline.ID,
		Status:     "success",
		StartTime:  now.Add(-time.Minute),
		EndTime:    now,
		Duration:   60,
		Logs:       "Test passed",
	}

	err = executionRepo.Create(execution)
	if err != nil {
		t.Fatalf("创建执行历史失败: %v", err)
	}

	testDB = db.GetDB()
	repo := NewMetricRepository(testDB)

	// 测试创建指标
	metric := &models.Metric{
		ExecutionID: execution.ID,
		Name:        "build_time",
		Value:       15.5,
	}

	err = repo.Create(metric)
	if err != nil {
		t.Fatalf("创建指标失败: %v", err)
	}

	if metric.ID == 0 {
		t.Fatal("指标 ID 未设置")
	}

	// 测试根据执行 ID 获取指标
	metrics, err := repo.GetByExecutionID(execution.ID)
	if err != nil {
		t.Fatalf("获取指标失败: %v", err)
	}

	if len(metrics) == 0 {
		t.Fatal("指标列表为空")
	}

	// 测试根据项目 ID 获取指标
	projectMetrics, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取项目指标失败: %v", err)
	}

	if len(projectMetrics) == 0 {
		t.Fatal("项目指标列表为空")
	}

	// 清理项目
	projectRepo.Delete(project.ID)
}

func TestOptimizationRepository(t *testing.T) {
	// 先创建一个项目
	projectRepo := NewProjectRepository(testDB)
	project := &models.Project{
		Name:          "测试项目",
		Description:   "这是一个测试项目",
		RepositoryURL: "https://github.com/test/test",
	}

	err := projectRepo.Create(project)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	testDB = db.GetDB()
	repo := NewOptimizationRepository(testDB)

	// 测试创建优化建议
	optimization := &models.Optimization{
		ProjectID:   project.ID,
		Type:        "build",
		Description: "构建时间过长",
		Suggestion:  "使用缓存加速构建",
		Applied:     false,
	}

	err = repo.Create(optimization)
	if err != nil {
		t.Fatalf("创建优化建议失败: %v", err)
	}

	if optimization.ID == 0 {
		t.Fatal("优化建议 ID 未设置")
	}

	// 测试根据项目 ID 获取优化建议
	optimizations, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取优化建议失败: %v", err)
	}

	if len(optimizations) == 0 {
		t.Fatal("优化建议列表为空")
	}

	// 测试获取未应用的优化建议
	unappliedOptimizations, err := repo.GetUnappliedByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取未应用的优化建议失败: %v", err)
	}

	if len(unappliedOptimizations) == 0 {
		t.Fatal("未应用的优化建议列表为空")
	}

	// 测试标记为已应用
	err = repo.MarkAsApplied(optimization.ID)
	if err != nil {
		t.Fatalf("标记优化建议为已应用失败: %v", err)
	}

	updatedOptimization, err := repo.GetByProjectID(project.ID)
	if err != nil {
		t.Fatalf("获取更新后的优化建议失败: %v", err)
	}

	if !updatedOptimization[0].Applied {
		t.Fatal("优化建议未标记为已应用")
	}

	// 清理项目
	projectRepo.Delete(project.ID)
}
