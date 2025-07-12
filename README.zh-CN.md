# UserCenter - 用户中心服务

[![Go 版本](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![测试覆盖率](https://img.shields.io/badge/Coverage-80%25-brightgreen.svg)](./coverage.html)
[![构建状态](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## 📖 目录

- [项目概述](#项目概述)
- [核心功能](#核心功能)
- [技术栈](#技术栈)
- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [故障排除](#故障排除)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [开发指南](#开发指南)
- [测试说明](#测试说明)
- [CI/CD](#cicd)
- [部署方案](#部署方案)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## 🎯 项目概述

UserCenter 是一个基于 Go 语言构建的生产就绪的用户中心服务，提供完整的用户管理功能，包括注册、登录、查询和列表等核心功能。该项目遵循标准 Go 项目布局，采用现代化的技术栈，支持高并发、高可用性和可扩展性。

### 核心功能

- 🔐 **用户认证**：基于 JWT 的用户注册和登录
- 🔍 **用户查询**：支持条件过滤的用户信息查询
- 📋 **用户列表**：支持分页和排序的用户列表
- 🏥 **健康检查**：服务状态监控端点
- 🛡️ **安全特性**：输入校验、速率限制、CORS 支持
- 🌍 **国际化**：多语言支持（中文/英文）
- 🔄 **优雅停机**：安全的服务关闭机制
- 📊 **可观测性**：完整的监控、日志和链路追踪

## 🚀 核心功能

### 认证与授权
- 基于 JWT 的无状态认证
- 使用 bcrypt 进行密码哈希（成本 12）
- 基于角色的访问控制
- Token 刷新机制
- 安全的会话管理

### 用户管理
- 支持邮箱验证的用户注册
- 用户资料管理
- 账户状态管理（活跃、非活跃、暂停）
- 软删除支持
- 批量用户操作
- UUID 用户标识符
- 密码强度验证

### API 特性
- RESTful API 设计
- 全面的输入验证
- 速率限制（通用、登录专用、注册专用）
- 请求 ID 追踪
- CORS 配置
- Swagger/OpenAPI 文档
- 国际化支持（中文/英文）
- 优雅的错误处理

### 监控与可观测性
- 所有依赖的健康检查端点
- Prometheus 指标收集
- 使用 Zap 的结构化日志
- 使用 OpenTelemetry 的分布式追踪
- 性能监控
- 实时性能分析（pprof）
- 自定义业务指标

## 🛠️ 技术栈

### 核心框架
- **Web 框架**：[Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- **依赖注入**：[Wire](https://github.com/google/wire) - 编译时依赖注入
- **API 文档**：[Swagger](https://github.com/swaggo/gin-swagger) - 自动生成 OpenAPI 3.0 文档

### 数据存储
- **主数据库**：[PostgreSQL](https://www.postgresql.org/) + [GORM](https://gorm.io/) - 用户核心数据
- **辅助数据库**：[MongoDB](https://www.mongodb.com/) - 日志和会话数据
- **缓存**：[Redis](https://redis.io/) - 高性能缓存
- **数据库迁移**：[Goose](https://github.com/pressly/goose) - 数据库版本控制
- **UUID 生成**：PostgreSQL pgcrypto 扩展

### 消息和任务处理
- **消息队列**：[Kafka](https://kafka.apache.org/) - 事件消费
- **异步任务**：[Asynq](https://github.com/hibiken/asynq) - 后台任务处理

### 监控和日志
- **日志**：[Zap](https://github.com/uber-go/zap) - 高性能结构化日志
- **监控**：[Prometheus](https://prometheus.io/) - 指标收集
- **分布式追踪**：[OpenTelemetry](https://opentelemetry.io/) - 分布式追踪

### 安全和工具
- **认证**：[JWT](https://github.com/golang-jwt/jwt) - 无状态认证
- **国际化**：[go-i18n](https://github.com/nicksnyder/go-i18n) - 多语言支持
- **配置**：YAML 配置文件
- **代码质量**：[golangci-lint](https://golangci-lint.run/) - 代码质量检查

## 📋 环境要求

### 系统要求
- Go 1.23.1 或更高版本
- PostgreSQL 13+
- MongoDB 5.0+
- Redis 6.0+
- Apache Kafka 2.8+

### 开发工具
```bash
# 安装 Go
# 参考：https://golang.org/doc/install

# 安装 Wire (依赖注入)
go install github.com/google/wire/cmd/wire@latest

# 安装 Mockgen (Mock 生成)
go install github.com/golang/mock/mockgen@latest

# 安装 Goose (数据库迁移)
go install github.com/pressly/goose/v3/cmd/goose@latest

# 安装 golangci-lint (代码质量检查)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# 安装 Swagger 生成工具
go install github.com/swaggo/swag/cmd/swag@latest

# 安装覆盖率工具
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest

# 安装性能分析工具
go install github.com/google/pprof@latest

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd user-center
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置环境
```bash
# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件
vim configs/config.yaml
```

### 4. 初始化数据库
```bash
# 运行数据库迁移
make migrate-up

# 或者手动运行
goose -dir migrations postgres "user=username password=password dbname=usercenter sslmode=disable" up
```

### 5. 生成 Wire 依赖注入代码
```bash
make wire
```

### 6. 生成 Swagger 文档
```bash
make swagger
```

### 7. 运行服务
```bash
# 开发环境
make run

# 或者直接运行
go run cmd/usercenter/main.go

# 生产环境
make build
./bin/usercenter
```

### 8. 验证服务
```bash
# 检查健康状态
curl http://localhost:8080/health

# 访问 Swagger 文档
open http://localhost:8080/swagger/index.html
```

## 🔧 故障排除

### 常见问题

#### 1. 数据库连接失败
```bash
# 检查 PostgreSQL 服务状态
sudo systemctl status postgresql

# 检查连接参数
psql -h localhost -U username -d usercenter

# 常见解决方案：
# - 确保 PostgreSQL 服务正在运行
# - 验证用户名和密码
# - 检查数据库是否存在
# - 确认端口 5432 可访问
```

#### 2. 数据库迁移失败
```bash
# 检查迁移状态
make migrate-status

# 重置数据库（谨慎使用）
make migrate-down
make migrate-up

# 常见解决方案：
# - 确保数据库表结构正确
# - 检查 pgcrypto 扩展是否已安装
# - 验证用户权限
```

#### 3. 服务启动失败
```bash
# 检查端口占用
lsof -i :8080

# 查看详细错误日志
make run 2>&1 | tee server.log

# 常见解决方案：
# - 确保所有依赖服务正在运行
# - 检查配置文件语法
# - 验证环境变量设置
```

#### 4. JWT 认证问题
```bash
# 检查 JWT 配置
grep -r "jwt" configs/

# 常见解决方案：
# - 确保 JWT 密钥已正确设置
# - 检查 Token 过期时间配置
# - 验证 Token 格式
```

#### 5. 依赖服务问题
```bash
# 检查 Redis 连接
redis-cli ping

# 检查 MongoDB 连接
mongosh --eval "db.runCommand('ping')"

# 检查 Kafka 连接
kafka-topics --list --bootstrap-server localhost:9092

# 常见解决方案：
# - 确保所有服务正在运行
# - 检查网络连接
# - 验证配置参数
```

### 日志分析
```bash
# 查看实时日志
tail -f logs/usercenter.log

# 搜索错误日志
grep -i error logs/usercenter.log

# 查看性能指标
curl http://localhost:8080/metrics
```

### 性能调优
```bash
# 启用性能分析
export USERCENTER_PROFILING=true

# 监控内存使用
go tool pprof http://localhost:8080/debug/pprof/heap

# 监控 CPU 使用
go tool pprof http://localhost:8080/debug/pprof/profile
```

## ⚙️ 配置说明

项目支持多种配置方式，优先级从高到低：

1. **环境变量**：`USERCENTER_` 前缀
2. **配置文件**：`configs/config.yaml`
3. **默认值**：代码中的默认配置

### 主要配置项

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test

database:
  postgres:
    host: "localhost"
    port: 5432
    user: "username"
    password: "password"
    dbname: "usercenter"
    sslmode: "disable"
  
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "usercenter_logs"
  
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0

kafka:
  brokers: ["localhost:9092"]
  topics:
    user_events: "user.events"

jwt:
  secret: "your-secret-key"
  expiry: "24h"

logging:
  level: "info"
  format: "json"

monitoring:
  prometheus:
    enabled: true
    port: 9090
  
  tracing:
    enabled: true
    endpoint: "http://localhost:14268/api/traces"
```

## 📚 API 文档

### Swagger 文档访问

启动服务后，可通过以下地址访问 API 文档：

- **Swagger UI**：http://localhost:8080/swagger/index.html
- **OpenAPI JSON**：http://localhost:8080/swagger/doc.json

### API 端点

#### 1. 健康检查
```bash
# 基础健康检查
GET /health

# 详细健康检查
GET /health/detailed

# 指标端点
GET /metrics
```

#### 2. 用户管理
```bash
# 用户注册
POST /api/v1/users/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password"
}

# 用户登录
POST /api/v1/users/login
{
  "email": "john@example.com",
  "password": "secure_password"
}

# 获取用户资料
GET /api/v1/users/profile
Authorization: Bearer <jwt_token>

# 更新用户资料
PUT /api/v1/users/profile
Authorization: Bearer <jwt_token>
{
  "username": "john_doe_updated",
  "email": "john.updated@example.com"
}

# 获取用户列表（支持分页和过滤）
GET /api/v1/users?page=1&limit=20&status=active&search=john
Authorization: Bearer <jwt_token>

# 获取特定用户
GET /api/v1/users/{id}
Authorization: Bearer <jwt_token>

# 删除用户
DELETE /api/v1/users/{id}
Authorization: Bearer <jwt_token>
```

## 🛠️ 开发指南

### 项目结构
```
user-center/
├── cmd/usercenter/          # 应用程序入口点
│   ├── main.go             # 主应用程序
│   └── wire.go             # Wire 依赖注入
├── internal/               # 私有应用程序代码
│   ├── config/             # 配置管理
│   ├── model/              # 领域实体（GORM 模型）
│   ├── dto/                # 数据传输对象
│   ├── service/            # 业务逻辑层
│   ├── repository/         # 数据访问层
│   ├── handler/            # HTTP 处理器（控制器）
│   ├── middleware/         # HTTP 中间件
│   ├── server/             # 服务器设置和路由
│   └── database/           # 数据库连接
├── pkg/                    # 共享包
│   ├── logger/             # 日志工具
│   └── jwt/                # JWT 工具
├── configs/                # 配置文件
├── migrations/             # 数据库迁移
├── docs/                   # 生成的文档
├── Makefile                # 构建和开发任务
├── Dockerfile              # 容器配置
└── README.md               # 此文件
```

### 可用的 Make 命令
```bash
# 开发
make run                    # 开发模式运行
make build                  # 构建二进制文件
make clean                  # 清理构建产物
make wire                   # 生成 Wire 依赖注入
make swagger                # 生成 Swagger 文档

# 测试
make test                   # 运行所有测试
make test-coverage          # 运行测试并生成覆盖率报告
make test-coverage-xml      # 运行测试并生成 XML 覆盖率报告
make test-short             # 仅运行短测试
make test-race              # 运行竞态检测测试
make mockgen                # 生成测试用的 Mock

# 数据库
make migrate-up             # 运行数据库迁移
make migrate-down           # 回滚数据库迁移
make migrate-status         # 检查迁移状态

# 代码质量
make lint                   # 运行 golangci-lint
make fmt                    # 格式化代码
make vet                    # 运行 go vet

# Docker
make docker-build           # 构建 Docker 镜像
make docker-run             # 运行 Docker 容器
make docker-clean           # 清理 Docker 产物

# 工具
make help                   # 显示所有可用命令
make profiling              # 启用性能分析
make logs                   # 查看实时日志
```

## 🧪 测试说明

### 运行测试
```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行测试并生成 XML 覆盖率报告（用于 CI）
make test-coverage-xml

# 仅运行单元测试（跳过集成测试）
make test-short

# 运行竞态检测测试
make test-race

# 生成测试用的 Mock
make mockgen

# 运行特定测试
go test -run TestUserService_CreateUser ./...
```

### 测试覆盖率
项目目标测试覆盖率达到 80% 以上。覆盖率报告生成在：
- `coverage.out` - 原始覆盖率数据
- `coverage.html` - HTML 覆盖率报告
- `coverage.xml` - XML 覆盖率报告（用于 CI 集成）

### 测试结构
- **单元测试**：测试单个函数和方法
- **集成测试**：测试数据库操作和 API 端点
- **Mock 测试**：使用 gomock 进行依赖模拟
- **Mock 生成**：使用 `mockgen` 自动生成 Mock
- **性能测试**：基准测试和压力测试
- **安全测试**：JWT 和认证测试

## 🔄 CI/CD

本项目使用 GitHub Actions 进行持续集成和部署。CI/CD 流水线包括代码质量检查、测试、构建和自动化部署。

### 工作流

#### 1. CI 工作流 (`ci.yml`)
- **触发条件**：推送到 `main`/`develop` 分支，Pull Requests
- **功能特性**：
  - 单元和集成测试（含覆盖率）
  - Mock 生成和依赖注入代码生成
  - XML 覆盖率报告（用于 CI 集成）
  - 针对快速执行进行优化（并行测试、缓存）

#### 2. 发布工作流 (`release.yml`)
- **触发条件**：版本标签推送（如 `v1.0.0`）
- **功能特性**：
  - 构建并发布 Docker 镜像到 GitHub Container Registry
  - 创建 GitHub Releases 并包含资源文件
  - 多架构支持（linux/amd64, linux/arm64）

#### 3. 部署工作流 (`deploy.yml`)
- **触发条件**：`main` 分支上 CI 成功完成后
- **功能特性**：
  - 自动部署到测试环境
  - 自动部署到生产环境
  - 部署通知

#### 4. 安全扫描工作流 (`security.yml`)
- **触发条件**：每周定时、手动触发、依赖变更
- **功能特性**：
  - 代码安全扫描（gosec）
  - 依赖漏洞检查（govulncheck）
  - Docker 镜像安全扫描（Trivy）
  - 文件系统安全扫描

### 设置

1. 在仓库设置中**启用 GitHub Actions**
2. **配置密钥**用于数据库连接和部署
3. **设置环境**用于测试和生产
4. **配置 Dependabot**用于自动依赖更新
5. **确保 `go.sum` 已提交**（不要忽略）以确保可重现构建

### 使用方法

```bash
# 创建新发布
git tag v1.0.0
git push origin v1.0.0

# 检查工作流状态
# 访问：https://github.com/username/user-center/actions

# 查看安全扫描结果
# 访问：https://github.com/username/user-center/security
```

详细配置和故障排除请参阅 [GitHub Actions 文档](docs/github-actions.md)。

## 🚀 部署方案

### Docker 部署
```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 或使用 docker-compose
docker-compose up -d
```

### 生产环境部署
```bash
# 构建生产版本
make build

# 设置环境变量
export USERCENTER_ENV=production
export USERCENTER_DB_HOST=your-db-host
export USERCENTER_DB_PASSWORD=your-db-password

# 运行服务
./bin/usercenter
```

### Kubernetes 部署
```bash
# 应用 Kubernetes 清单
kubectl apply -f k8s/

# 检查部署状态
kubectl get pods -l app=usercenter
```

## 🤝 贡献指南

我们欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'Add amazing feature'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

### 开发规范
- 遵循 Go 编码标准
- 编写全面的测试
- 更新文档
- 使用规范的提交信息
- 确保所有测试通过后再提交

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

## 🔗 相关链接

- [英文文档](README.md)
- [项目主页](https://github.com/username/user-center)
- [问题反馈](https://github.com/username/user-center/issues)
- [讨论区](https://github.com/username/user-center/discussions)
- [Docker Hub](https://hub.docker.com/r/username/user-center)
- [GitHub Container Registry](https://github.com/username/user-center/packages) 