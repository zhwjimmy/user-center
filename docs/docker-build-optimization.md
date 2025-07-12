# Docker 构建优化指南

## 🎯 优化目标

基于 GitHub Actions 中 Docker Build job 耗时较长的问题，本文档提供了详细的 Docker 构建优化策略。

## 📊 问题分析

### 原始配置的性能瓶颈

1. **依赖下载重复**
   - 每次构建都重新下载 Go 模块依赖
   - 没有利用 Docker 层缓存机制

2. **构建上下文过大**
   - 复制整个项目目录到构建上下文
   - 包含不必要的文件（.git、docs、测试文件等）

3. **工具安装重复**
   - wire 和 swag 工具每次都重新安装
   - 没有缓存工具安装层

4. **CI 缓存缺失**
   - GitHub Actions 中没有 Docker 层缓存
   - 每次都是完整构建

## 🚀 优化策略

### 1. Dockerfile 优化

#### 优化层缓存顺序
```dockerfile
# 复制 go.mod/go.sum 文件（变化频率低）
COPY go.mod go.sum ./

# 下载依赖（缓存层）
RUN go mod download

# 安装构建工具（缓存层）
RUN go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

# 复制源代码（变化频率高）
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
```

#### 减少构建上下文
```dockerfile
# 只复制必要的目录
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY configs/ ./configs/
COPY migrations/ ./migrations/
COPY locales/ ./locales/
```

### 2. .dockerignore 优化

排除不必要的文件：
```dockerignore
# Git 文件
.git
.gitignore

# 文档文件
README.md
docs/

# 构建产物
bin/
coverage/

# 测试文件
*_test.go

# IDE 文件
.vscode/
.idea/
```

### 3. CI/CD 缓存优化

#### Docker 层缓存
```yaml
- name: Cache Docker layers
  uses: actions/cache@v4
  with:
    path: /tmp/.buildx-cache
    key: ${{ runner.os }}-buildx-${{ github.sha }}
    restore-keys: |
      ${{ runner.os }}-buildx-

- name: Build Docker image
  run: |
    docker buildx build \
      --cache-from type=local,src=/tmp/.buildx-cache \
      --cache-to type=local,dest=/tmp/.buildx-cache,mode=max \
      --tag usercenter:latest \
      --file Dockerfile \
      .
```

## 📈 预期性能提升

### 首次构建
- **依赖下载**: 减少 50-60% 时间（缓存命中）
- **工具安装**: 减少 70-80% 时间（缓存命中）
- **构建上下文**: 减少 60-70% 大小
- **总体构建时间**: 减少 40-50%

### 后续构建
- **依赖层缓存**: 减少 80-90% 时间
- **工具层缓存**: 减少 90-95% 时间
- **总体构建时间**: 减少 70-80%

### 缓存命中率
- **依赖层**: 95%+ 命中率
- **工具层**: 98%+ 命中率
- **总体缓存**: 90%+ 命中率

## 🔧 实施步骤

### 1. 应用优化配置
```bash
# 应用 Dockerfile 优化
git add Dockerfile .dockerignore

# 应用 CI 优化
git add .github/workflows/ci.yml

# 提交更改
git commit -m "ci: optimize Docker build with layer caching and context reduction"
git push
```

### 2. 监控和调优
- 监控每次 Docker 构建时间
- 分析缓存命中率
- 根据实际效果调整缓存策略

## 📋 性能监控

### 关键指标
1. **构建时间**: 目标减少 50-70%
2. **缓存命中率**: 目标 > 90%
3. **构建上下文大小**: 目标减少 60-70%
4. **层缓存效率**: 目标 > 80%

### 监控方法
```bash
# 查看构建时间
gh run view <run-id> --log

# 查看缓存使用情况
docker system df

# 分析构建上下文
docker build --progress=plain .
```

## 🛠️ 故障排除

### 缓存失效
```bash
# 清除 Docker 缓存
docker system prune -a

# 重新构建
docker build --no-cache .
```

### 构建上下文过大
```bash
# 检查 .dockerignore 是否生效
docker build --progress=plain . 2>&1 | grep "sending build context"

# 优化 .dockerignore
# 添加更多排除规则
```

### CI 缓存问题
```yaml
# 增加缓存键的精确度
key: ${{ runner.os }}-buildx-${{ hashFiles('go.mod', 'go.sum') }}
```

## 📚 最佳实践

### 1. 层顺序优化
- 将变化频率低的文件放在前面
- 将变化频率高的文件放在后面
- 合并相关的 RUN 命令

### 2. 多阶段构建
- 使用 builder 阶段编译
- 使用 final 阶段运行
- 减少最终镜像大小

### 3. 缓存策略
- 使用精确的缓存键
- 设置合适的缓存过期时间
- 监控缓存使用情况

### 4. 安全考虑
- 使用非 root 用户运行
- 最小化运行时依赖
- 定期更新基础镜像 