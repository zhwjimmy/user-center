# Docker 环境设置指南

## 概述

本项目使用 Docker Compose 来管理所有依赖服务，包括数据库、缓存、消息队列和监控工具。

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

### 1. 启动所有服务

```bash
# 启动所有服务
make docker-compose-up

# 或者直接使用 docker-compose
docker-compose up -d
```

### 2. 检查服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f postgres
```

### 3. 停止服务

```bash
# 停止所有服务
make docker-compose-down

# 或者直接使用 docker-compose
docker-compose down
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

### 1. 更新配置文件

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

### 2. 运行应用

```bash
# 启动依赖服务
make docker-compose-up

# 运行应用
make run-dev
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
   ```

2. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs -f [service-name]
   
   # 重启特定服务
   docker-compose restart [service-name]
   ```

3. **数据丢失**
   ```bash
   # 删除所有数据（谨慎使用）
   docker-compose down -v
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
``` 