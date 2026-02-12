# CI/CD 管道管理系统 - 前端

这是一个基于Vue 3的轻量版前端应用，用于管理CI/CD管道系统。

## 技术栈

- Vue 3
- Vite
- Axios
- 原生CSS

## 功能特性

- **项目管理**：创建、查看、删除项目
- **技术栈分析**：分析项目技术栈
- **管道配置**：生成CI管道配置
- **执行管理**：执行管道、查看执行历史
- **指标分析**：查看执行指标
- **优化建议**：分析执行历史并提供优化建议

## 项目结构

```
frontend/
├── src/
│   ├── main.js          # Vue应用入口
│   ├── App.vue          # 主组件
│   └── style.css        # 全局样式
├── index.html           # HTML入口文件
├── package.json         # 项目依赖
├── vite.config.js       # Vite配置
└── .gitignore           # Git忽略文件
```

## 安装和运行

### 前提条件

- Node.js 18+
- npm 9+

### 安装步骤

1. **安装依赖**

```bash
npm install
```

2. **构建项目**

```bash
npm run build
```

3. **运行开发服务器**

```bash
npm run dev
```

前端应用将运行在 http://localhost:3000

### 后端服务

前端需要与后端服务器配合使用：

- 后端服务器运行在 http://localhost:8080
- API请求会通过Vite代理自动转发到后端

## API集成

前端应用集成了以下API端点：

- **项目管理**：/api/projects
- **技术栈分析**：/api/projects/{id}/analyze, /api/projects/{id}/tech-stack
- **管道配置**：/api/projects/{id}/generate-pipeline
- **执行管理**：/api/projects/{id}/execute, /api/projects/{id}/executions
- **执行详情**：/api/executions/{id}, /api/executions/{id}/metrics, /api/executions/{id}/logs
- **优化建议**：/api/projects/{id}/analyze-optimization

## 使用说明

1. **创建项目**：在项目管理页面填写项目名称和路径
2. **分析技术栈**：点击项目的"分析技术栈"按钮
3. **生成管道**：点击项目的"生成管道"按钮
4. **执行管道**：点击项目的"执行管道"按钮
5. **查看执行**：点击项目的"查看执行"按钮
6. **分析优化**：点击项目的"分析优化"按钮

## 界面设计

- 响应式布局，适配不同屏幕尺寸
- 现代化的UI设计，使用蓝色主题
- 清晰的导航结构
- 直观的操作界面
- 实时的消息提示

## 注意事项

- 确保后端服务器运行在 http://localhost:8080
- 前端开发服务器运行在 http://localhost:3000
- 所有API请求会自动代理到后端服务器
