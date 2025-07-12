# CI/CD 性能优化指南

## 🎯 优化目标

基于 GitHub Actions 工作流中 "Initialize containers" 阶段耗时较长的问题，本文档提供了详细的性能优化策略。

## 📊 问题分析

### 原始配置的性能瓶颈

1. **容器初始化慢**
   - PostgreSQL 使用完整镜像 (`postgres:15`)
   - 健康检查间隔过长 (10s)
   - 重试次数过多 (5次)

2. **工具安装效率低**
   - 每次重新安装 Go 工具
   - 缓存策略不够精确
   - 串行安装工具

3. **依赖管理优化空间**
   - 缺少 Go 模块缓存
   - 依赖下载串行化

## 🚀 优化策略

### 1. 容器优化

#### 使用轻量级镜像
```yaml
# 优化前
image: postgres:15

# 优化后
image: postgres:15-alpine
```

#### 减少健康检查时间
```yaml
# 优化前
--health-interval 10s
--health-timeout 5s
--health-retries 5

# 优化后
--health-interval 3s
--health-timeout 2s
--health-retries 2
```

### 2. 缓存优化

#### 合并缓存策略
```yaml
- name: Cache Go tools and dependencies
  uses: actions/cache@v4
  with:
    path: |
      $HOME/go/bin
      $HOME/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('go.mod', 'go.sum', 'Makefile') }}
    restore-keys: |
      ${{ runner.os }}-go-${{ hashFiles('go.mod', 'go.sum') }}-
      ${{ runner.os }}-go-
```

### 3. 并行化优化

#### 并行安装工具
```bash
# 并行安装工具
(
  echo "Installing tools..."
  go install github.com/golang/mock/mockgen@v1.6.0 &
  go install github.com/axw/gocov/gocov@latest &
  go install github.com/AlekSi/gocov-xml@latest &
  go install github.com/google/wire/cmd/wire@v0.6.0 &
  wait
)
```

#### 并行下载依赖
```bash
# 并行下载依赖
(
  echo "Downloading dependencies..."
  go mod download &
  go mod tidy &
  wait
)
```

#### 并行代码生成
```bash
# 并行生成代码
make mock &
make wire &
wait
```

### 4. 测试优化

#### 增加并行度
```bash
# 优化前
go test -v -coverprofile=coverage/coverage.out -covermode=atomic -p=4 ./...

# 优化后
go test -v -coverprofile=coverage/coverage.out -covermode=atomic -p=8 ./...
```

## 📈 预期性能提升

### 容器初始化时间
- **PostgreSQL**: 从 ~30s 减少到 ~15s
- **Redis**: 从 ~10s 减少到 ~5s
- **总体**: 减少 40-50% 的容器启动时间

### 工具安装时间
- **首次运行**: 保持不变 (~30s)
- **后续运行**: 从 ~30s 减少到 ~5s (缓存命中)
- **总体**: 减少 80% 的工具安装时间

### 依赖下载时间
- **首次运行**: 保持不变 (~20s)
- **后续运行**: 从 ~20s 减少到 ~2s (缓存命中)
- **总体**: 减少 90% 的依赖下载时间

### 测试执行时间
- **并行度提升**: 从 4 增加到 8
- **预期提升**: 减少 30-40% 的测试执行时间

## 🔧 实施步骤

### 1. 立即优化 (ci.yml)
```bash
# 应用基础优化
git add .github/workflows/ci.yml
git commit -m "ci: optimize container initialization and caching"
git push
```

### 2. 高级优化 (ci-fast.yml)
```bash
# 应用高级优化
git add .github/workflows/ci-fast.yml
git commit -m "ci: add fast CI configuration with parallel execution"
git push
```

### 3. 监控和调优
- 监控每次 CI 运行时间
- 分析缓存命中率
- 根据实际效果调整参数

## 📋 性能监控

### 关键指标
1. **容器启动时间**: 目标 < 20s
2. **工具安装时间**: 目标 < 10s (缓存命中)
3. **依赖下载时间**: 目标 < 5s (缓存命中)
4. **测试执行时间**: 目标减少 30-40%
5. **总体 CI 时间**: 目标减少 50-60%

### 监控方法
```bash
# 查看 CI 运行时间
gh run list --limit 10

# 查看具体步骤时间
gh run view <run-id> --log
```

## 🛠️ 故障排除

### Codecov 配置问题
```yaml
# 错误配置 (v3 版本)
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage/coverage.xml  # ❌ 已废弃

# 正确配置 (v3 版本)
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage/coverage.xml  # ✅ 新参数名
```

### 缓存问题
```bash
# 清除缓存
gh run rerun <run-id> --clear-cache
```

### 容器健康检查失败
```yaml
# 增加重试次数
--health-retries 3
```

### 并行执行冲突
```bash
# 串行执行作为备选
make mock
make wire
```

## 📚 参考资源

- [GitHub Actions 缓存最佳实践](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows)
- [Docker 健康检查优化](https://docs.docker.com/engine/reference/builder/#healthcheck)
- [Go 测试并行化](https://golang.org/pkg/testing/#hdr-Parallel_Subtests)

## 🔄 持续优化

### 定期审查
- 每月审查 CI 性能指标
- 分析缓存命中率
- 评估新工具和依赖的影响

### 渐进式优化
- 先应用基础优化
- 监控效果后再应用高级优化
- 根据项目规模调整参数

---

*最后更新: 2024年12月* 