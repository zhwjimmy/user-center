# 本地开发指南

## 概述

本指南介绍如何在本地环境中开发和调试 UserCenter 服务。我们采用混合架构：依赖服务通过 Docker Compose 管理，应用服务本地运行。

## 架构设计

### 服务分层
```
┌─────────────────┐    ┌─────────────────┐
│   应用服务      │    │   依赖服务      │
│   (本地运行)    │◄──►│   (Docker)      │
│                 │    │                 │
│ • UserCenter    │    │ • PostgreSQL    │
│ • 热重载        │    │ • MongoDB       │
│ • 调试支持      │    │ • Redis         │
│ • 快速开发      │    │ • Kafka         │
│                 │    │ • Jaeger        │
└─────────────────┘    │ • Prometheus    │
                       └─────────────────┘
```

### 优势
- 🚀 **快速开发**：本地运行应用服务，支持热重载
- 🔧 **易于调试**：可以直接调试本地代码
- 📊 **资源优化**：依赖服务容器化，应用服务本地化
- 🔄 **灵活配置**：开发环境配置更灵活

## 快速开始

### 1. 启动开发环境

```bash
# 一键启动开发环境
make dev-start

# 或者分步执行
make setup-env      # 创建 .env 文件
make docker-compose-up  # 启动依赖服务
```

### 2. 启动应用服务

```bash
# 开发模式（推荐）
make run-dev

# 或者直接运行
source .env && ./bin/usercenter
```

### 3. 验证服务

```bash
# 检查应用服务健康状态
curl http://localhost:8080/health

# 检查依赖服务状态
make docker-compose-ps
```

## 开发工作流

### 日常开发流程

```bash
# 1. 启动开发环境
make dev-start

# 2. 启动应用服务
make run-dev

# 3. 进行开发...

# 4. 停止开发环境
make dev-stop
```

### 代码修改

应用服务支持热重载，代码修改后会自动重新编译和启动。

### 依赖服务管理

```bash
# 查看服务状态
make docker-compose-ps

# 查看服务日志
make docker-compose-logs

# 重启特定服务
docker-compose restart [service-name]
```

## 环境配置

### 环境变量

推荐使用 `.env` 文件管理环境变量：

```bash
# .env 文件内容
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PORT=5432
USERCENTER_DATABASE_POSTGRES_USER=postgres
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter
USERCENTER_DATABASE_POSTGRES_SSLMODE=disable

# 可选：Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 可选：Kafka 配置
KAFKA_BROKERS=localhost:9092
```

### 使用 direnv（推荐）

安装 `direnv` 实现自动环境变量加载：

```bash
# macOS
brew install direnv

# 在项目根目录创建 .envrc 文件
echo "source .env" > .envrc
direnv allow
```

## 服务访问

### 应用服务
- **API 服务**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **指标端点**: http://localhost:8080/metrics

### 依赖服务
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **PostgreSQL**: localhost:5432
- **MongoDB**: localhost:27017
- **Redis**: localhost:6379
- **Kafka**: localhost:9092

## 调试技巧

### 1. 应用服务调试

```bash
# 使用 delve 调试器
dlv debug ./cmd/usercenter

# 或者使用 IDE 调试
# VS Code: 按 F5
# GoLand: 右键 -> Debug
```

### 2. 数据库调试

```bash
# 连接 PostgreSQL
docker-compose exec postgres psql -U postgres -d usercenter

# 连接 MongoDB
docker-compose exec mongodb mongosh

# 连接 Redis
docker-compose exec redis redis-cli
```

### 3. 日志查看

```bash
# 查看应用服务日志
tail -f logs/usercenter.log

# 查看依赖服务日志
make docker-compose-logs

# 查看特定服务日志
docker-compose logs -f postgres
```

## 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行短测试
make test-short

# 运行测试并生成覆盖率报告
make test-coverage
```

### 集成测试

```bash
# 运行集成测试
make test-integration

# 运行 Kafka 测试
make test-kafka
```

### API 测试

```bash
# 测试用户注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 测试用户登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   lsof -i :8080
   lsof -i :5432
   ```

2. **数据库连接失败**
   ```bash
   # 检查环境变量
   env | grep USERCENTER
   
   # 检查数据库服务
   make docker-compose-ps
   ```

3. **Kafka 集群 ID 冲突**
   ```bash
   # 清理 Kafka 数据
   docker-compose down
   docker volume rm user-center_kafka_data user-center_zookeeper_data
   docker-compose up -d
   ```

### 重置环境

```bash
# 完全重置开发环境
make dev-stop
docker-compose down -v
make dev-start
```

## IDE 配置

### VS Code

在 `.vscode/settings.json` 中添加：

```json
{
  "go.toolsEnvVars": {
    "USERCENTER_DATABASE_POSTGRES_HOST": "localhost",
    "USERCENTER_DATABASE_POSTGRES_PORT": "5432",
    "USERCENTER_DATABASE_POSTGRES_USER": "postgres",
    "USERCENTER_DATABASE_POSTGRES_PASSWORD": "password",
    "USERCENTER_DATABASE_POSTGRES_DBNAME": "usercenter",
    "USERCENTER_DATABASE_POSTGRES_SSLMODE": "disable"
  }
}
```

### GoLand

在运行配置中添加环境变量：
- `USERCENTER_DATABASE_POSTGRES_HOST=localhost`
- `USERCENTER_DATABASE_POSTGRES_PORT=5432`
- `USERCENTER_DATABASE_POSTGRES_USER=postgres`
- `USERCENTER_DATABASE_POSTGRES_PASSWORD=password`
- `USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter`
- `USERCENTER_DATABASE_POSTGRES_SSLMODE=disable`

## 性能优化

### 1. 开发环境优化

```bash
# 使用 Go 模块缓存
export GOMODCACHE=$HOME/.cache/go-mod

# 启用 Go 构建缓存
export GOCACHE=$HOME/.cache/go-build
```

### 2. 数据库优化

```bash
# PostgreSQL 连接池优化
# 在 configs/config.yaml 中调整
database:
  postgres:
    max_open_conns: 25
    max_idle_conns: 10
    max_lifetime: "5m"
```

### 3. 监控和调试

```bash
# 启用性能分析
export USERCENTER_PROFILING=true

# 查看内存使用
go tool pprof http://localhost:8080/debug/pprof/heap

# 查看 CPU 使用
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 最佳实践

### 1. 代码组织
- 遵循 Go 项目标准布局
- 使用依赖注入（Wire）
- 编写单元测试和集成测试

### 2. 配置管理
- 使用环境变量管理敏感信息
- 使用 `.env` 文件简化开发
- 区分开发、测试、生产环境

### 3. 日志和监控
- 使用结构化日志
- 添加适当的指标
- 配置分布式追踪

### 4. 错误处理
- 使用适当的错误类型
- 添加错误上下文
- 实现优雅降级

## 有用的命令

```bash
# 开发环境管理
make dev-start          # 启动开发环境
make dev-stop           # 停止开发环境
make run-dev            # 运行应用服务

# 服务管理
make docker-compose-ps  # 查看服务状态
make docker-compose-logs # 查看服务日志
make setup-env          # 设置环境变量

# 测试
make test               # 运行测试
make test-coverage      # 生成覆盖率报告
make test-kafka         # 运行 Kafka 测试

# 代码质量
make lint               # 代码检查
make fmt                # 代码格式化
make vet                # 代码验证
``` 