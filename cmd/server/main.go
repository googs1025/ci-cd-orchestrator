package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"ci-cd-orchestrator/cmd/server/api"
	"ci-cd-orchestrator/internal/cicd/template"
	"ci-cd-orchestrator/internal/db"
	"ci-cd-orchestrator/internal/repository"
	"ci-cd-orchestrator/pkg/migration"
)

func main() {
	// 初始化数据库
	err := db.Init()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 执行数据库迁移
	dbConn := db.GetDB()
	err = migration.Run(dbConn)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化内置模板
	initBuiltinTemplates(dbConn)

	// 创建路由
	router := api.SetupRouter(dbConn)

	// 启动服务器
	port := 8080 // 使用端口 8080
	serverAddr := fmt.Sprintf(":%d", port)

	log.Printf("API 服务器启动在 http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// initBuiltinTemplates 初始化内置模板
func initBuiltinTemplates(dbConn *sql.DB) {
	// 创建模板仓库
	templateRepo := repository.NewTemplateRepository(dbConn)

	// 创建模板管理器，这会自动初始化内置模板
	template.NewTemplateManager(templateRepo)

	log.Println("内置模板初始化完成")
}
