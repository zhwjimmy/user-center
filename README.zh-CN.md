# UserCenter - 用户中心服务

[![Go 版本](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![测试覆盖率](https://img.shields.io/badge/Coverage-80%25-brightgreen.svg)](./coverage.html)
[![构建状态](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## 🎯 项目概述

UserCenter 是一个基于 Go 语言构建的生产就绪的用户中心服务，采用清洁架构设计和事件驱动模式。提供完整的用户管理功能，包括注册、登录、查询和列表等核心功能，支持高并发、高可用性和可扩展性。

### ✨ 核心功能

- 🔐 **JWT 认证**：基于 JWT 的用户注册和登录，安全密码哈希
- 👥 **用户管理**：UUID 用户标识，软删除支持
- 🚀 **事件驱动架构**：使用 Apache Kafka 进行异步处理
- 🏗️ **清洁架构**：清晰的关注点分离
- 🛡️ **安全特性**：速率限制、CORS、输入验证
- 📊 **可观测性**：健康检查、指标、结构化日志
- 🌍 **国际化**：多语言支持（中文/英文）
- 🔄 **优雅停机**：依赖管理和安全关闭

## 🚀 快速开始

### 环境要求
- Go 1.23.1+
- PostgreSQL 13+
- Redis 6.0+
- Apache Kafka 2.8+

### 1. 克隆和设置
```bash
git clone <repository-url>
cd user-center
go mod download
```

### 2. 启动依赖服务
```bash
docker-compose up -d
```

### 3. 配置环境
```bash
# 创建 .env 文件
cat > .env << EOF
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PORT=5432
USERCENTER_DATABASE_POSTGRES_USER=postgres
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter
USERCENTER_DATABASE_POSTGRES_SSLMODE=disable
EOF
```

### 4. 运行服务
```bash
# 生成依赖注入代码
go generate ./cmd/usercenter

# 开发模式运行
make run-dev
```

### 5. 验证安装
```bash
# 健康检查
curl http://localhost:8080/health

# API 文档
open http://localhost:8080/swagger/index.html
```

## 📚 文档

📖 **完整文档**：[docs/README.md](docs/README.md)

### 快速链接
- 🏗️ [架构设计](docs/architecture.md)
- 🚀 [快速开始指南](docs/getting-started.md)
- 📖 [API 参考](docs/api-reference.md)
- 🛠️ [开发指南](docs/development.md)
- 🚀 [部署指南](docs/deployment.md)
- 🔧 [配置指南](docs/configuration.md)
- 🐛 [故障排除](docs/troubleshooting.md)
- 📊 [Kafka 集成](docs/kafka-integration.md)

## 🏗️ 架构设计

UserCenter 遵循清洁架构原则，在基础设施层和业务层之间实现清晰的分离：

```
┌─────────────────────────────────────────────────────────────┐
│                    业务层                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   服务层    │  │  处理器层   │  │   事件层    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  基础设施层                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   数据库    │  │    缓存     │  │   消息队列  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### 核心设计原则
- **依赖倒置**：业务层定义接口，基础设施层实现
- **事件驱动**：使用 Kafka 进行异步事件处理
- **关注点分离**：层间清晰边界
- **可测试性**：依赖注入便于测试

## 🛠️ 开发

### 项目结构
```
user-center/
├── cmd/usercenter/          # 应用程序入口点
├── internal/               # 私有应用程序代码
│   ├── infrastructure/     # 外部依赖（数据库、缓存、消息队列）
│   ├── events/            # 事件驱动架构
│   ├── service/           # 业务逻辑
│   ├── handler/           # HTTP 处理器
│   └── middleware/        # HTTP 中间件
├── docs/                  # 文档
├── configs/               # 配置文件
└── migrations/            # 数据库迁移
```

### 可用命令
```bash
# 开发
make run-dev              # 热重载运行
make build                # 构建二进制文件
make wire                 # 生成依赖注入

# 测试
make test                 # 运行所有测试
make test-coverage        # 运行覆盖率测试

# 数据库
make migrate-up           # 运行迁移
make migrate-down         # 回滚迁移

# 文档
make swagger              # 生成 API 文档
```

## 🤝 贡献

我们欢迎贡献！详情请参阅 [贡献指南](docs/contributing.md)。

### 快速贡献步骤
1. Fork 项目
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -am 'Add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

## 🔗 链接

- [英文文档](README.md)
- [项目主页](https://github.com/zhwjimmy/user-center)
- [问题反馈](https://github.com/zhwjimmy/user-center/issues)
- [讨论区](https://github.com/zhwjimmy/user-center/discussions) 