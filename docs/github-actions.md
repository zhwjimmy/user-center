# GitHub Actions 配置说明

本项目配置了完整的 GitHub Actions CI/CD 流水线，包括代码质量检查、测试、构建、发布和部署。

## 工作流概览

### 1. CI 工作流 (`ci.yml`)

**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 创建 Pull Request 到 `main` 或 `develop` 分支

**功能：**
- 代码质量检查（linting、安全扫描）
- 单元测试和集成测试
- 代码覆盖率报告
- 多平台构建
- Docker 镜像构建

**作业：**
- `test`: 运行测试套件，包括代码质量检查
- `build`: 构建多平台二进制文件
- `docker`: 构建 Docker 镜像

### 2. Release 工作流 (`release.yml`)

**触发条件：**
- 推送版本标签（如 `v1.0.0`）

**功能：**
- 构建并发布 Docker 镜像到 GitHub Container Registry
- 创建 GitHub Release
- 上传构建产物

### 3. Deploy 工作流 (`deploy.yml`)

**触发条件：**
- CI 工作流成功完成后（仅在 `main` 分支）

**功能：**
- 自动部署到测试环境
- 自动部署到生产环境
- 部署通知

### 4. Security 工作流 (`security.yml`)

**触发条件：**
- 每周一上午 9 点 UTC（定时任务）
- 手动触发
- 修改依赖文件时

**功能：**
- 代码安全扫描（gosec）
- 依赖漏洞检查（govulncheck）
- Docker 镜像安全扫描（Trivy）
- 文件系统安全扫描

## 环境配置

### 必需的环境变量

在 GitHub 仓库设置中配置以下 Secrets：

#### 数据库配置
```
POSTGRES_HOST=your-postgres-host
POSTGRES_PORT=5432
POSTGRES_USER=your-postgres-user
POSTGRES_PASSWORD=your-postgres-password
POSTGRES_DB=your-postgres-db
```

#### Redis 配置
```
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
```

#### 部署配置
```
DEPLOY_SSH_KEY=your-deploy-ssh-private-key
DEPLOY_HOST=your-deploy-host
DEPLOY_USER=your-deploy-user
```

### 环境设置

在 GitHub 仓库设置中创建以下环境：

1. **staging**: 测试环境
2. **production**: 生产环境

每个环境可以配置特定的 Secrets 和变量。

## 使用指南

### 1. 开发流程

1. 创建功能分支
2. 开发并提交代码
3. 创建 Pull Request
4. CI 工作流自动运行
5. 代码审查通过后合并

### 2. 发布流程

1. 确保所有测试通过
2. 创建版本标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. Release 工作流自动触发
4. 检查 GitHub Release 和 Container Registry

### 3. 部署流程

1. 代码合并到 `main` 分支
2. CI 工作流成功完成
3. Deploy 工作流自动触发
4. 依次部署到测试和生产环境

## 监控和通知

### 1. 工作流状态

- 在 GitHub 仓库的 "Actions" 标签页查看所有工作流状态
- 设置工作流失败时的通知

### 2. 安全扫描结果

- 在 GitHub 仓库的 "Security" 标签页查看安全扫描结果
- 配置安全警报通知

### 3. 代码覆盖率

- 集成 Codecov 查看代码覆盖率报告
- 在 Pull Request 中显示覆盖率变化

## 自定义配置

### 1. 修改触发条件

编辑对应工作流文件的 `on` 部分：

```yaml
on:
  push:
    branches: [ main, develop, feature/* ]
  pull_request:
    branches: [ main ]
```

### 2. 添加新的作业

在工作流文件中添加新的 `job`：

```yaml
new-job:
  name: New Job
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4
    # 添加更多步骤...
```

### 3. 自定义构建步骤

修改 Makefile 中的命令，或直接在工作流中使用 Go 命令：

```yaml
- name: Custom build
  run: |
    go build -ldflags="-s -w" -o myapp ./cmd/myapp
```

## 故障排除

### 1. 常见问题

**工作流失败：**
- 检查日志输出
- 验证环境变量配置
- 确认依赖服务可用性

**测试失败：**
- 检查测试环境配置
- 验证数据库连接
- 查看测试日志

**构建失败：**
- 检查 Go 版本兼容性
- 验证依赖包版本
- 查看构建日志

### 2. 调试技巧

1. 使用 `workflow_dispatch` 手动触发工作流
2. 在本地运行相同的命令
3. 检查 GitHub Actions 日志
4. 使用 `echo` 输出调试信息

## 最佳实践

1. **分支策略：**
   - 使用 `main` 作为主分支
   - 使用 `develop` 作为开发分支
   - 功能开发使用功能分支

2. **提交信息：**
   - 使用清晰的提交信息
   - 遵循约定式提交规范

3. **测试覆盖：**
   - 保持高测试覆盖率
   - 包含单元测试和集成测试

4. **安全扫描：**
   - 定期运行安全扫描
   - 及时修复安全漏洞

5. **依赖管理：**
   - 定期更新依赖
   - 使用 Dependabot 自动化更新

## 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go 官方文档](https://golang.org/doc/)
- [Docker 官方文档](https://docs.docker.com/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) 