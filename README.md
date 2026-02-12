# CI/CD 管道管理系统

## 项目概述

CI/CD 管道管理系统是一个用于自动化构建、测试和部署软件项目的平台。该系统能够分析项目技术栈，生成适合的 CI/CD 配置，并支持多平台执行。

## 核心功能

- **项目管理**：创建、查看、删除项目
- **技术栈分析**：自动识别项目的技术栈（语言、框架、构建工具等）
- **CI/CD 配置生成**：根据技术栈生成适合的 CI/CD 配置文件
- **模板管理**：管理 CI/CD 配置模板，支持内置模板和自定义模板
- **执行管理**：执行 CI/CD 管道，查看执行历史和指标
- **优化建议**：分析执行数据，生成优化建议

## 支持的平台

- **GitHub Actions**：真实的 CI/CD 平台（现阶段只实现生成 ci.yaml 功能，其他功能还未实现）
- **Mock**：模拟的 CI/CD 平台，用于测试和开发

## 文档目录

详细文档请查看 `docs/` 目录：

- **docs/system/** - 系统文档，包括对外解释文档和设计文档
- **docs/test/** - 测试文档，包括 Mock CI 平台自测文档和接口测试文档
- **docs/implementation/** - 实现文档，包括文档更新记录

## 快速开始

### 后端部署

1. **环境要求**：
   - Go 1.20+
   - SQLite 3+

2. **启动步骤**：
   ```bash
   # 克隆代码库
   git clone https://github.com/googs1025/ci-cd-orchestrator
   cd ci-cd-orchestrator

   # 安装依赖
   go mod tidy

   # 启动服务器
   go run cmd/server/main.go
   ```

   服务器默认启动在 `http://localhost:8080`

### 前端部署

1. **环境要求**：
   - Node.js 16+
   - npm 7+

2. **启动步骤**：
   ```bash
   # 进入前端目录
   cd frontend

   # 安装依赖
   npm install

   # 启动开发服务器
   npm run dev
   ```

   前端开发服务器默认启动在 `http://localhost:3000`