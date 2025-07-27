# UserCenter Migration Tool

这是 UserCenter 项目的数据库迁移工具，基于 Goose 实现，提供简单易用的数据库迁移管理功能。

## 功能特性

- ✅ 执行数据库迁移 (up)
- ✅ 回滚数据库迁移 (down)
- ✅ 查看迁移状态 (status)
- ✅ 创建新的迁移文件 (create)
- ✅ 支持配置文件和环境变量
- ✅ 完整的日志记录
- ✅ 优雅的错误处理

## 快速开始

### 1. 构建工具

```bash
make build-migrate
```

### 2. 查看帮助

```bash
./bin/migrate -help
```

### 3. 查看版本

```bash
./bin/migrate -version
```

## 使用方法

### 查看迁移状态

```bash
make migrate-status
# 或者
./bin/migrate -action status
```

### 执行迁移

```bash
make migrate-up
# 或者
./bin/migrate -action up
```

### 回滚迁移

```bash
make migrate-down
# 或者
./bin/migrate -action down -steps 1
```

### 创建新迁移文件

```bash
make migrate-create name=add_user_profile
# 或者
./bin/migrate -action create -name add_user_profile
```

## 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-action` | string | "up" | 迁移操作: up, down, status, create |
| `-config` | string | "configs/config.yaml" | 配置文件路径 |
| `-name` | string | "" | 迁移名称 (create 操作必需) |
| `-steps` | int | 1 | 回滚步数 (down 操作) |
| `-version` | bool | false | 显示版本信息 |

## 配置文件

工具使用与主应用相同的配置文件 (`configs/config.yaml`)，主要使用以下配置：

```yaml
database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    dbname: "usercenter"
    sslmode: "disable"
```

## 环境变量

支持通过环境变量覆盖配置：

```bash
export USERCENTER_DATABASE_POSTGRES_HOST=localhost
export USERCENTER_DATABASE_POSTGRES_PORT=5432
export USERCENTER_DATABASE_POSTGRES_USER=postgres
export USERCENTER_DATABASE_POSTGRES_PASSWORD=password
export USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter
export USERCENTER_DATABASE_POSTGRES_SSLMODE=disable
```

## 迁移文件格式

迁移文件使用 Goose 格式，包含 Up 和 Down 两个部分：

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
```

## 依赖要求

- Go 1.23.1+
- PostgreSQL 15+
- Goose v3+

## 安装依赖

```bash
# 安装 Goose
make install-goose
```

## 故障排除

### 1. 数据库连接失败

检查数据库服务是否运行：

```bash
docker-compose ps postgres
```

### 2. 配置文件不存在

确保配置文件存在：

```bash
ls -la configs/config.yaml
```

### 3. 权限问题

确保有足够的数据库权限：

```bash
psql -h localhost -U postgres -d usercenter -c "SELECT version();"
```

## 开发说明

### 项目结构

```
cmd/migrate/
├── main.go          # 主入口文件
├── wire.go          # 依赖注入配置
├── wire_gen.go      # Wire 生成的代码
└── migrations.go    # 迁移逻辑实现
```

### 构建和测试

```bash
# 构建
go build -o bin/migrate ./cmd/migrate

# 测试
go test ./cmd/migrate/...

# 生成 Wire 代码
wire ./cmd/migrate
```

## 最佳实践

1. **备份数据库**: 在生产环境执行迁移前，务必备份数据库
2. **测试迁移**: 在开发环境充分测试迁移脚本
3. **版本控制**: 将迁移文件纳入版本控制
4. **回滚计划**: 准备回滚计划以应对迁移失败
5. **监控日志**: 关注迁移执行日志，及时发现问题 