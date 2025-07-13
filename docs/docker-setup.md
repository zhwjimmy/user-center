# Docker 环境设置指南

## 概述

本项目使用 Docker Compose 来管理所有依赖服务，包括数据库、缓存、消息队列和监控工具。**应用服务本身通过本地方式运行，便于开发和调试。**

## 架构说明

### 服务分层
- **依赖服务**：通过 Docker Compose 管理（数据库、缓存、消息队列等）
- **应用服务**：本地运行，便于热重载、调试和开发

### 优势
- 🚀 **快速开发**：本地运行应用服务，支持热重载
- 🔧 **易于调试**：可以直接调试本地代码
- 📊 **资源优化**：依赖服务容器化，应用服务本地化
- 🔄 **灵活配置**：开发环境配置更灵活

## 服务列表

### 核心服务
- **PostgreSQL 15** - 主数据库 (端口: 5432)
- **MongoDB 6.0** - 日志数据库 (端口: 27017)
- **Redis 7** - 缓存和会话存储 (端口: 6379)

### 消息队列
- **Apache Kafka 7.3.0** - 事件流处理 (端口: 9092)
- **Apache Zookeeper 7.3.0** - Kafka 协调服务 (端口: 2181)

### 监控和可观测性
- **Jaeger 1.50** - 分布式追踪 (端口: 16686 - UI, 14268 - HTTP collector)
- **Prometheus 2.48.1** - 指标收集 (端口: 9090)

## 快速开始

### 1. 启动依赖服务

```bash
# 启动所有依赖服务
docker-compose up -d

# 检查服务状态
docker-compose ps
```

### 2. 配置本地环境

创建 `.env` 文件（可选，用于自动加载环境变量）：

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

# 加载环境变量
source .env
```

### 3. 启动应用服务

```bash
# 方式一：使用环境变量启动
USERCENTER_DATABASE_POSTGRES_HOST=localhost \
USERCENTER_DATABASE_POSTGRES_PORT=5432 \
USERCENTER_DATABASE_POSTGRES_USER=postgres \
USERCENTER_DATABASE_POSTGRES_PASSWORD=password \
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter \
USERCENTER_DATABASE_POSTGRES_SSLMODE=disable \
./bin/usercenter

# 方式二：如果已创建 .env 文件
source .env && ./bin/usercenter

# 方式三：开发模式（自动重新编译）
make run-dev
```

### 4. 验证服务

```bash
# 检查应用服务健康状态
curl http://localhost:8080/health

# 检查依赖服务状态
docker-compose ps
```

## 服务访问地址

### 数据库服务
- **PostgreSQL**: `localhost:5432`
  - 用户名: `postgres`
  - 密码: `password`
  - 数据库: `usercenter`

- **MongoDB**: `localhost:27017`
  - 数据库: `usercenter_logs`

- **Redis**: `localhost:6379`
  - 无密码

### 消息队列
- **Kafka**: `localhost:9092`
- **Zookeeper**: `localhost:2181`

### 监控服务
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090

### 应用服务
- **API 服务**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Swagger 文档**: http://localhost:8080/swagger/index.html

## 开发工作流

### 1. 日常开发流程

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 启动应用服务（开发模式）
make run-dev

# 3. 进行开发...

# 4. 停止服务
docker-compose down
```

### 2. 代码修改后重启

```bash
# 应用服务会自动重新编译和启动
# 或者手动重启
pkill usercenter
make run-dev
```

### 3. 依赖服务管理

```bash
# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 重启特定服务
docker-compose restart [service-name]

# 停止所有服务
docker-compose down
```

## 配置说明

### 网络配置
所有服务都在 `mynetwork` 网络中，使用子网 `172.20.0.0/16`。

### 数据持久化
所有数据都通过 Docker volumes 持久化：
- `postgres_data` - PostgreSQL 数据
- `mongodb_data` - MongoDB 数据
- `redis_data` - Redis 数据
- `kafka_data` - Kafka 数据
- `zookeeper_data` - Zookeeper 数据
- `jaeger_data` - Jaeger 数据
- `prometheus_data` - Prometheus 数据

### 健康检查
所有服务都配置了健康检查，确保服务启动顺序正确。

## 开发环境配置

### 1. 环境变量配置

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

### 2. 使用 direnv（推荐）

安装 `direnv` 实现自动环境变量加载：

```bash
# macOS
brew install direnv

# 在项目根目录创建 .envrc 文件
echo "source .env" > .envrc
direnv allow
```

### 3. 配置文件

确保 `configs/config.yaml` 中的配置与 Docker 服务匹配：

```yaml
database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    dbname: "usercenter"
  
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

monitoring:
  tracing:
    endpoint: "http://localhost:14268/api/traces"
```

## 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   lsof -i :5432
   lsof -i :27017
   lsof -i :6379
   lsof -i :9092
   lsof -i :8080
   ```

2. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs -f [service-name]
   
   # 重启特定服务
   docker-compose restart [service-name]
   ```

3. **应用服务连接失败**
   ```bash
   # 检查环境变量
   env | grep USERCENTER
   
   # 检查数据库连接
   psql -h localhost -U postgres -d usercenter
   ```

4. **Kafka 集群 ID 冲突**
   ```bash
   # 清理 Kafka 数据
   docker-compose down
   docker volume rm user-center_kafka_data user-center_zookeeper_data
   docker-compose up -d
   ```

### 重置环境

```bash
# 停止所有服务并删除数据
docker-compose down -v

# 重新启动
docker-compose up -d
```

## 生产环境注意事项

1. **安全性**
   - 修改默认密码
   - 使用环境变量管理敏感信息
   - 配置防火墙规则

2. **性能**
   - 调整资源限制
   - 配置适当的连接池大小
   - 启用数据压缩

3. **监控**
   - 配置告警规则
   - 设置日志轮转
   - 监控磁盘空间

## 有用的命令

```bash
# 进入容器
docker-compose exec postgres psql -U postgres -d usercenter
docker-compose exec redis redis-cli
docker-compose exec mongodb mongosh

# 备份数据
docker-compose exec postgres pg_dump -U postgres usercenter > backup.sql

# 查看资源使用
docker stats

# 清理未使用的资源
docker system prune -a

# 查看应用服务进程
ps aux | grep usercenter

# 查看应用服务端口
lsof -i :8080
```

## 开发工具集成

### IDE 配置

#### VS Code
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

#### GoLand
在运行配置中添加环境变量：
- `USERCENTER_DATABASE_POSTGRES_HOST=localhost`
- `USERCENTER_DATABASE_POSTGRES_PORT=5432`
- `USERCENTER_DATABASE_POSTGRES_USER=postgres`
- `USERCENTER_DATABASE_POSTGRES_PASSWORD=password`
- `USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter`
- `USERCENTER_DATABASE_POSTGRES_SSLMODE=disable` 