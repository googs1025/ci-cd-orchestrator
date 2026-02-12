# Mock CI 平台自测文档

## 1. 测试环境

- **API 基础 URL**: `http://localhost:8080/api/v1`
- **服务器状态**: 确保服务器已启动
  ```bash
  # 启动服务器
  go run cmd/server/main.go
  ```
- **预期服务器输出**:
  ```
  2026/02/12 09:53:14 数据库连接成功
  2026/02/12 09:53:14 数据库迁移成功
  2026/02/12 09:53:14 API 服务器启动在 http://localhost:8080
  ```

## 2. 测试命令

### 2.1 执行 Mock CI 管道

**命令**:
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/1/execute?platform=mock'
```

**预期响应**:
```json
{
  "status": "success",
  "data": {
    "execution_id": "<执行ID>",
    "status": "running",
    "platform": "mock"
  },
  "message": "执行管道成功"
}
```

### 2.2 获取执行详情

**命令**:
```bash
curl -X GET 'http://localhost:8080/api/v1/executions/<执行ID>'
```

**预期响应** (执行中):
```json
{
  "status": "success",
  "data": {
    "id": "<执行ID>",
    "status": "running",
    "logs": [
      {
        "stage": "init",
        "message": "Starting init stage"
      },
      {
        "stage": "init",
        "message": "Completed init stage in 2 seconds"
      },
      {
        "stage": "build",
        "message": "Starting build stage"
      }
    ]
  },
  "message": "获取执行详情成功"
}
```

**预期响应** (执行完成):
```json
{
  "status": "success",
  "data": {
    "id": "<执行ID>",
    "status": "success",
    "duration": 10,
    "metrics": {
      "total_duration": 10,
      "cpu_usage": 73.82,
      "memory_usage": 82.62,
      "test_coverage": 92.57,
      "build_size": 58720256,
      "deployment_time": 3
    },
    "platform_data": {
      "resource_config": {
        "cpu": {
          "actual": 0.74,
          "limit": 2,
          "request": 1,
          "usage_percent": 73.82
        },
        "memory": {
          "actual": 1.65,
          "limit": 4,
          "request": 2,
          "usage_percent": 82.62
        }
      }
    },
    "logs": [
      {
        "stage": "init",
        "message": "Starting init stage"
      },
      {
        "stage": "init",
        "message": "Completed init stage in 2 seconds"
      },
      {
        "stage": "build",
        "message": "Starting build stage"
      },
      {
        "stage": "build",
        "message": "Completed build stage in 3 seconds"
      },
      {
        "stage": "test",
        "message": "Starting test stage"
      },
      {
        "stage": "test",
        "message": "Completed test stage in 3 seconds"
      },
      {
        "stage": "deploy",
        "message": "Starting deploy stage"
      },
      {
        "stage": "deploy",
        "message": "Completed deploy stage in 1 seconds"
      },
      {
        "stage": "complete",
        "message": "Starting complete stage"
      },
      {
        "stage": "complete",
        "message": "Completed complete stage in 1 seconds"
      },
      {
        "stage": "complete",
        "message": "Execution completed successfully in 10 seconds"
      }
    ]
  },
  "message": "获取执行详情成功"
}
```

### 2.3 获取执行指标

**命令**:
```bash
curl -X GET 'http://localhost:8080/api/v1/executions/<执行ID>/metrics'
```

**预期响应**:
```json
{
  "status": "success",
  "data": {
    "total_duration": 10,
    "stage_durations": {
      "init": 2,
      "build": 3,
      "test": 3,
      "deploy": 1,
      "complete": 1
    },
    "success_rate": 1,
    "cpu_usage": 73.82,
    "memory_usage": 82.62,
    "test_coverage": 92.57,
    "build_size": 58720256,
    "deployment_time": 3
  },
  "message": "获取执行指标成功"
}
```

### 2.4 获取执行日志

**命令**:
```bash
curl -X GET 'http://localhost:8080/api/v1/executions/<执行ID>/logs'
```

**预期响应**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:21.877768+08:00",
      "level": "info",
      "stage": "init",
      "message": "Starting init stage"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:23.878976+08:00",
      "level": "info",
      "stage": "init",
      "message": "Completed init stage in 2 seconds"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:23.878989+08:00",
      "level": "info",
      "stage": "build",
      "message": "Starting build stage"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:26.87976+08:00",
      "level": "info",
      "stage": "build",
      "message": "Completed build stage in 3 seconds"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:26.879763+08:00",
      "level": "info",
      "stage": "test",
      "message": "Starting test stage"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:29.880755+08:00",
      "level": "info",
      "stage": "test",
      "message": "Completed test stage in 3 seconds"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:29.880757+08:00",
      "level": "info",
      "stage": "deploy",
      "message": "Starting deploy stage"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:30.881246+08:00",
      "level": "info",
      "stage": "deploy",
      "message": "Completed deploy stage in 1 seconds"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:30.881247+08:00",
      "level": "info",
      "stage": "complete",
      "message": "Starting complete stage"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:31.881948+08:00",
      "level": "info",
      "stage": "complete",
      "message": "Completed complete stage in 1 seconds"
    },
    {
      "id": "<日志ID>",
      "execution_id": "<执行ID>",
      "timestamp": "2026-02-12T09:53:31.88195+08:00",
      "level": "info",
      "stage": "complete",
      "message": "Execution completed successfully in 10 seconds"
    }
  ],
  "message": "获取执行日志成功"
}
```

### 2.5 停止执行

**命令**:
```bash
curl -X POST 'http://localhost:8080/api/v1/executions/<执行ID>/stop'
```

**预期响应**:
```json
{
  "status": "success",
  "data": {
    "execution_id": "<执行ID>",
    "status": "cancelled"
  },
  "message": "停止执行成功"
}
```

### 2.6 获取执行历史

**命令**:
```bash
curl -X GET 'http://localhost:8080/api/v1/projects/1/executions'
```

**预期响应**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "<执行ID>",
      "project_id": "1",
      "platform": "mock",
      "status": "success",
      "start_time": "2026-02-12T09:53:21.877675+08:00",
      "end_time": "2026-02-12T09:53:31.881948+08:00",
      "duration": 10
    }
  ],
  "message": "获取执行历史成功"
}
```

### 2.7 分析管道优化建议

**命令**:
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/1/analyze-optimization'
```

**预期响应**:
```json
{
  "status": "success",
  "data": {
    "suggestions": [
      {
        "type": "performance",
        "description": "执行时长过长",
        "suggestion": "考虑优化构建和测试流程，减少执行时间"
      }
    ],
    "metrics": {
      "average_duration": 10.0,
      "success_rate": 1.0,
      "failure_rate": 0.0,
      "recent_durations": [10, 10, 10, 10, 10]
    }
  },
  "message": "分析优化建议成功"
}
```

### 2.8 获取优化建议列表

**命令**:
```bash
curl -X GET 'http://localhost:8080/api/v1/projects/1/optimization-suggestions'
```

**预期响应**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "<建议ID>",
      "project_id": 1,
      "type": "performance",
      "description": "执行时长过长",
      "suggestion": "考虑优化构建和测试流程，减少执行时间",
      "applied": false,
      "created_at": "2026-02-12T09:53:31.88195+08:00"
    }
  ],
  "message": "获取优化建议列表成功"
}
```

### 2.9 应用优化建议

**命令**:
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/1/apply-optimization' \
  -H 'Content-Type: application/json' \
  -d '{"optimization_id": <建议ID>}'
```

**预期响应**:
```json
{
  "status": "success",
  "data": null,
  "message": "应用优化建议成功"
}
```

## 3. 测试场景

### 3.1 基本执行测试

**步骤**:
1. 启动服务器
2. 执行 Mock CI 管道
3. 获取执行详情，验证状态为 "running"
4. 等待约 10 秒
5. 再次获取执行详情，验证状态为 "success"
6. 检查指标数据和资源使用情况

**预期结果**:
- 执行成功完成
- 收集到完整的指标数据
- 资源使用情况在合理范围内

### 3.2 资源配置测试

**步骤**:
1. 执行 Mock CI 管道
2. 等待执行完成
3. 检查 `platform_data.resource_config` 字段

**预期结果**:
```json
{
  "resource_config": {
    "cpu": {
      "actual": 0.7-0.9,
      "limit": 2,
      "request": 1,
      "usage_percent": 70-90
    },
    "memory": {
      "actual": 1.5-1.9,
      "limit": 4,
      "request": 2,
      "usage_percent": 75-95
    }
  }
}
```

### 3.3 取消执行测试

**步骤**:
1. 执行 Mock CI 管道
2. 立即调用停止执行接口
3. 获取执行详情，验证状态为 "cancelled"
4. 检查日志中是否包含取消信息

**预期结果**:
- 执行状态为 "cancelled"
- 日志中包含 "Execution cancelled by user" 信息

## 4. 指标说明

### 4.1 收集的指标

| 指标名称 | 类型 | 范围 | 描述 |
|---------|------|------|------|
| **总执行时间** | 整数（秒） | 10 秒 | 整个 CI 流程的总执行时长 |
| **CPU 使用率** | 浮点数（%） | 70-90% | 基于 request 的使用率 |
| **内存使用率** | 浮点数（%） | 75-95% | 基于 request 的使用率 |
| **测试覆盖率** | 浮点数（%） | 75-95% | 模拟的测试代码覆盖率 |
| **构建大小** | 整数（字节） | 10-60 MB | 模拟的构建产物大小 |
| **部署时间** | 整数（秒） | 1-5 秒 | 模拟的部署过程时长 |

### 4.2 资源配置与实际使用比较

| 资源类型 | 配置上限 | 配置请求 | 实际使用 | 使用率 |
|---------|---------|---------|---------|-------|
| CPU | 2 核 | 1 核 | 0.7-0.9 核 | 70-90% |
| 内存 | 4 GB | 2 GB | 1.5-1.9 GB | 75-95% |

## 5. 故障排除

### 5.1 常见问题

| 问题 | 可能原因 | 解决方案 |
|------|---------|---------|
| 服务器启动失败 | 端口被占用 | 检查并关闭占用 8080 端口的进程 |
| 执行失败 | 服务器未运行 | 确保服务器已启动 |
| 指标数据缺失 | 执行未完成 | 等待执行完成后再查询指标 |
| 资源配置信息缺失 | CI 配置解析失败 | 检查 CI 配置格式是否正确 |

### 5.2 调试命令

**检查服务器状态**:
```bash
curl -I http://localhost:8080/api/v1/projects
```

**检查执行状态**:
```bash
curl -X GET 'http://localhost:8080/api/v1/executions/<执行ID>' | jq '.data.status'
```

**查看完整响应**:
```bash
curl -X GET 'http://localhost:8080/api/v1/executions/<执行ID>' | jq
```

## 6. 测试脚本

### 6.1 完整测试脚本

```bash
#!/bin/bash

# 启动服务器（如果未启动）
# go run cmd/server/main.go &
sleep 2

echo "=== 测试 1: 执行 Mock CI 管道 ==="
execution_response=$(curl -s -X POST 'http://localhost:8080/api/v1/projects/1/execute?platform=mock')
echo "执行响应:"
echo "$execution_response"

execution_id=$(echo "$execution_response" | grep -o '"execution_id":"[^"]*"' | cut -d'"' -f4)
echo "\n执行ID: $execution_id"

echo "\n=== 测试 2: 获取执行详情（执行中） ==="
sleep 3
detail_response=$(curl -s -X GET "http://localhost:8080/api/v1/executions/$execution_id")
echo "执行中状态:"
echo "$detail_response" | grep -o '"status":"[^"]*"'

echo "\n=== 测试 3: 等待执行完成 ==="
sleep 8

echo "\n=== 测试 4: 获取执行详情（执行完成） ==="
detail_response=$(curl -s -X GET "http://localhost:8080/api/v1/executions/$execution_id")
echo "执行完成状态:"
echo "$detail_response" | grep -o '"status":"[^"]*"'

echo "\n=== 测试 5: 检查指标数据 ==="
echo "CPU 使用率:"
echo "$detail_response" | grep -o '"cpu_usage":[0-9.]*'
echo "内存使用率:"
echo "$detail_response" | grep -o '"memory_usage":[0-9.]*'

echo "\n=== 测试 6: 检查资源配置 ==="
echo "资源配置:"
echo "$detail_response" | grep -A 15 '"resource_config"'

echo "\n=== 测试 7: 获取执行历史 ==="
history_response=$(curl -s -X GET 'http://localhost:8080/api/v1/projects/1/executions')
echo "执行历史:"
echo "$history_response"

echo "\n=== 测试完成 ==="
```

**使用方法**:
```bash
# 保存为 test_mock_ci.sh
chmod +x test_mock_ci.sh
./test_mock_ci.sh
```

## 7. 总结

Mock CI 平台提供了完整的 CI 执行模拟功能，包括：

1. **执行流程模拟**：模拟完整的 CI 执行流程，包括 init、build、test、deploy、complete 等阶段
2. **指标收集**：收集 CPU、内存、执行时间等关键指标
3. **资源配置管理**：支持设置资源使用上限和请求量
4. **配置与实际比较**：提供资源配置与实际使用的详细比较
5. **完整的 API 接口**：支持执行、查询、取消等操作

通过本测试文档，您可以全面测试 Mock CI 平台的各项功能，验证其是否符合预期。