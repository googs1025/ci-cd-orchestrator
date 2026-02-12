package api

import (
	"database/sql"
	"net/http"

	"ci-cd-orchestrator/cmd/server/handlers"
	"ci-cd-orchestrator/cmd/server/middleware"
	"ci-cd-orchestrator/internal/execution"
	"ci-cd-orchestrator/internal/repository"
)

// SetupRouter 设置路由
func SetupRouter(dbConn *sql.DB) http.Handler {
	// 创建路由器
	mux := http.NewServeMux()

	// 初始化执行管理器和引擎
	executionManager := execution.NewManager()
	mockEngine := execution.NewMockEngine()
	githubEngine := execution.NewGitHubActionsEngine()
	executionManager.RegisterEngine("mock", mockEngine)
	executionManager.RegisterEngine("github_actions", githubEngine)

	// 初始化仓库
	templateRepo := repository.NewTemplateRepository(dbConn)

	// 创建处理器实例
	projectHandler := handlers.NewProjectHandler()
	techStackHandler := handlers.NewTechStackHandler()
	pipelineHandler := handlers.NewPipelineHandler(templateRepo)
	templateHandler := handlers.NewTemplateHandler(templateRepo)
	executionHandler := handlers.NewExecutionHandler(executionManager)
	metricHandler := handlers.NewMetricHandler()
	optimizationHandler := handlers.NewOptimizationHandler(executionManager)

	// API 版本前缀
	apiPrefix := "/api/v1"

	// 项目管理路由
	mux.HandleFunc(apiPrefix+"/projects", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET":  projectHandler.ListProjects,
		"POST": projectHandler.CreateProject,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET":    projectHandler.GetProject,
		"PUT":    projectHandler.UpdateProject,
		"DELETE": projectHandler.DeleteProject,
	}))

	// 技术栈分析路由
	mux.HandleFunc(apiPrefix+"/projects/{id}/analyze", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": techStackHandler.AnalyzeProject,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}/tech-stack", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": techStackHandler.GetTechStack,
	}))

	// 管道配置路由
	mux.HandleFunc(apiPrefix+"/projects/{id}/generate-pipeline", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": pipelineHandler.GeneratePipeline,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}/pipeline", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": pipelineHandler.GetPipeline,
		"PUT": pipelineHandler.UpdatePipeline,
	}))

	// 执行历史路由
	mux.HandleFunc(apiPrefix+"/projects/{id}/execute", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": executionHandler.ExecutePipeline,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}/executions", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": executionHandler.ListExecutions,
	}))
	mux.HandleFunc(apiPrefix+"/executions/{id}", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": executionHandler.GetExecution,
	}))
	mux.HandleFunc(apiPrefix+"/executions/{id}/stop", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": executionHandler.StopExecution,
	}))
	mux.HandleFunc(apiPrefix+"/executions/{id}/metrics", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": executionHandler.GetExecutionMetrics,
	}))
	mux.HandleFunc(apiPrefix+"/executions/{id}/logs", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": executionHandler.GetExecutionLogs,
	}))

	// 指标路由
	mux.HandleFunc(apiPrefix+"/projects/{id}/metrics", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": metricHandler.ListMetrics,
	}))

	// 优化建议路由
	mux.HandleFunc(apiPrefix+"/projects/{id}/analyze-optimization", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": optimizationHandler.AnalyzeOptimization,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}/optimization-suggestions", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET": optimizationHandler.ListOptimizationSuggestions,
	}))
	mux.HandleFunc(apiPrefix+"/projects/{id}/apply-optimization", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": optimizationHandler.ApplyOptimization,
	}))

	// 模板管理路由
	mux.HandleFunc(apiPrefix+"/templates", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET":  templateHandler.ListTemplates,
		"POST": templateHandler.CreateTemplate,
	}))
	mux.HandleFunc(apiPrefix+"/templates/{id}", middleware.MethodHandler(map[string]http.HandlerFunc{
		"GET":    templateHandler.GetTemplate,
		"PUT":    templateHandler.UpdateTemplate,
		"DELETE": templateHandler.DeleteTemplate,
	}))
	mux.HandleFunc(apiPrefix+"/templates/{id}/reset", middleware.MethodHandler(map[string]http.HandlerFunc{
		"POST": templateHandler.ResetTemplate,
	}))

	return mux
}
