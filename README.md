# UserCenter - 用户中心服务

[![Go Version](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Coverage-80%25-brightgreen.svg)](./coverage.html)

## 项目简介

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

## 技术栈

### 核心框架
- **Web 框架**：[Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- **依赖注入**：[Wire](https://github.com/google/wire) - 编译时依赖注入
- **API 文档**：[Swagger](https://github.com/swaggo/gin-swagger) - 自动生成 OpenAPI 3.0 文档

### 数据存储
- **主数据库**：[PostgreSQL](https://www.postgresql.org/) + [GORM](https://gorm.io/) - 用户核心数据
- **辅助数据库**：[MongoDB](https://www.mongodb.com/) - 日志和会话数据
- **缓存**：[Redis](https://redis.io/) - 高性能缓存
- **数据库迁移**：[Goose](https://github.com/pressly/goose) - 数据库版本控制

### 消息和任务
- **消息队列**：[Kafka](https://kafka.apache.org/) - 事件消费
- **异步任务**：[Asynq](https://github.com/hibiken/asynq) - 后台任务处理

### 监控和日志
- **日志**：[Zap](https://github.com/uber-go/zap) - 高性能结构化日志
- **监控**：[Prometheus](https://prometheus.io/) - 指标收集
- **链路追踪**：[OpenTelemetry](https://opentelemetry.io/) - 分布式追踪

### 其他
- **认证**：[JWT](https://github.com/golang-jwt/jwt) - 无状态认证
- **国际化**：[go-i18n](https://github.com/nicksnyder/go-i18n) - 多语言支持
- **配置**：YAML 配置文件
- **代码质量**：[golangci-lint](https://golangci-lint.run/) - 代码规范检查

## 依赖和前提条件

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

# 安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 安装 Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# 安装 golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# 安装 Swagger 生成工具
go install github.com/swaggo/swag/cmd/swag@latest
```

## 安装和运行

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

## 配置说明

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

## API 文档和使用示例

### Swagger 文档访问

启动服务后，可通过以下地址访问 API 文档：

- **Swagger UI**：http://localhost:8080/swagger/index.html
- **OpenAPI JSON**：http://localhost:8080/swagger/doc.json

### API 端点

#### 1. 健康检查
```bash
curl -X GET http://localhost:8080/health
```

#### 2. 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh-CN" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword123"
  }'
```

#### 3. 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123"
  }'
```

#### 4. 查询用户（需要认证）
```bash
curl -X GET "http://localhost:8080/api/v1/users/1" \
  -H "Authorization: Bearer <jwt-token>"
```

#### 5. 用户列表（支持分页和排序）
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&size=10&sort=created_at&order=desc" \
  -H "Authorization: Bearer <jwt-token>"
```

## 国际化支持

### 支持的语言
- **中文**：`zh-CN`
- **英文**：`en-US`

### 使用方式

通过 HTTP 头部 `Accept-Language` 指定语言：

```bash
# 中文响应
curl -H "Accept-Language: zh-CN" http://localhost:8080/api/v1/users

# 英文响应
curl -H "Accept-Language: en-US" http://localhost:8080/api/v1/users
```

### 添加新语言

1. 在 `locales/` 目录下添加新的语言文件
2. 更新 `internal/i18n/i18n.go` 中的语言配置
3. 重新编译项目

## Request ID 追踪

### 自动生成
每个 HTTP 请求都会自动生成唯一的 Request ID，用于：
- 日志记录
- 链路追踪
- 错误调试

### 手动指定
可通过 HTTP 头部 `X-Request-ID` 手动指定：

```bash
curl -H "X-Request-ID: custom-request-123" http://localhost:8080/api/v1/users
```

### 日志中的 Request ID
所有日志都包含 Request ID 字段：

```json
{
  "level": "info",
  "time": "2023-12-01T10:00:00Z",
  "request_id": "req-123456789",
  "message": "User login successful",
  "user_id": 1
}
```

## 优雅停机

服务支持优雅停机，确保在收到终止信号时：

1. 停止接受新请求
2. 等待当前请求完成（最多 30 秒）
3. 关闭数据库连接
4. 清理资源

### 支持的信号
- `SIGINT` (Ctrl+C)
- `SIGTERM` (Docker stop)

### 配置停机超时
```yaml
server:
  shutdown_timeout: "30s"
```

## 开发指南

### 代码生成
```bash
# 生成 Wire 依赖注入代码
make wire

# 生成 Swagger 文档
make swagger

# 生成 Mock 文件
make mock
```

### 测试
```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行基准测试
make benchmark
```

### 代码质量检查
```bash
# 运行 linter
make lint

# 格式化代码
make fmt

# 检查代码安全性
make security
```

### 数据库操作
```bash
# 创建新的迁移文件
make migrate-create name=add_user_table

# 应用迁移
make migrate-up

# 回滚迁移
make migrate-down

# 查看迁移状态
make migrate-status
```

## 部署指南

### Docker 部署
```bash
# 构建镜像
make docker-build

# 运行容器
docker run -p 8080:8080 user-center:latest
```

### Docker Compose 部署
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f usercenter
```

### Kubernetes 部署
```bash
# 应用配置
kubectl apply -f deployments/k8s/

# 查看服务状态
kubectl get pods -l app=usercenter
```

## 监控和可观测性

### Prometheus 指标
- HTTP 请求指标：`http_requests_total`
- 响应时间：`http_request_duration_seconds`
- 数据库连接：`database_connections`
- 自定义业务指标

访问地址：http://localhost:9090/metrics

### 链路追踪
使用 OpenTelemetry 进行分布式追踪，支持：
- Jaeger
- Zipkin
- 其他兼容后端

### 日志聚合
推荐使用 ELK Stack：
- Elasticsearch：日志存储
- Logstash：日志处理
- Kibana：日志可视化

## 故障排除

### 常见问题

1. **服务启动失败**
   - 检查配置文件格式
   - 确认数据库连接
   - 查看端口占用情况

2. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接参数
   - 确认网络连通性

3. **JWT 认证失败**
   - 检查 JWT 密钥配置
   - 确认 Token 格式
   - 验证 Token 过期时间

### 日志调试
```bash
# 查看实时日志
tail -f logs/usercenter.log

# 根据 Request ID 查询日志
grep "req-123456789" logs/usercenter.log
```

## 贡献指南

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 提交更改：`git commit -am 'Add new feature'`
4. 推送分支：`git push origin feature/new-feature`
5. 创建 Pull Request

### 代码规范
- 遵循 Go 官方代码规范
- 通过 golangci-lint 检查
- 保持测试覆盖率 80% 以上
- 添加必要的注释和文档

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目维护者：[Your Name](mailto:your.email@example.com)
- 问题反馈：[GitHub Issues](https://github.com/your-org/user-center/issues)
- 文档：[项目文档](https://docs.your-domain.com/user-center) 