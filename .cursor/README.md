# Cursor Rules for UserCenter Project

这个目录包含了按照最新 Cursor MDC 规范组织的项目规则文件。

## 📁 规则文件结构

```
.cursor/rules/
├── core-rules/                 # 核心项目规则
│   ├── project-context-always.mdc      # 项目上下文 (Always)
│   ├── communication-rules.mdc         # 沟通规范 (Always)
│   └── version-control.mdc             # 版本控制 (Auto Attached)
├── go-rules/                   # Go 语言规则  
│   ├── go-coding-standards-auto.mdc     # Go 编码标准 (Auto Attached)
│   ├── dependency-injection.mdc         # 依赖注入 (Auto Attached)
│   ├── kafka-best-practices.mdc         # Kafka 最佳实践 (Auto Attached)
│   └── ci-cd-best-practices.mdc         # CI/CD 最佳实践 (Auto Attached)
├── api-rules/                  # API 设计规则
│   └── rest-api-design-agent.mdc        # REST API 设计 (Agent Requested)
├── database-rules/             # 数据库相关规则
│   └── database-patterns-auto.mdc       # 数据库模式 (Auto Attached)
├── testing-rules/              # 测试相关规则
│   ├── testing-standards-auto.mdc       # 测试标准 (Auto Attached)
│   └── go-testing-best-practices.mdc    # Go 测试最佳实践 (Auto Attached)
├── middleware-rules/           # 中间件相关规则
│   └── middleware-patterns-auto.mdc     # 中间件模式 (Auto Attached)
└── security-rules/             # 安全相关规则
    └── security-review-manual.mdc       # 安全审查 (Manual)
```

## 🎯 规则类型说明

### Always Rules (总是应用)
- **文件**: `*-always.mdc`
- **触发**: 每次对话和 Cmd+K 请求
- **用途**: 提供项目基础上下文和核心标准
- **示例**: 项目架构信息、核心开发原则

### Auto Attached Rules (自动附加)
- **文件**: `*-auto.mdc` 
- **触发**: 当编辑匹配 `globs` 模式的文件时
- **用途**: 针对特定文件类型的编码标准
- **示例**: Go 文件编码规范、测试文件标准

### Agent Requested Rules (代理请求)
- **文件**: `*-agent.mdc`
- **触发**: AI 代理根据上下文判断是否需要
- **用途**: 特定场景下的指导原则
- **示例**: API 设计规范、架构决策指南

### Manual Rules (手动引用)
- **文件**: `*-manual.mdc`
- **触发**: 使用 `@规则名` 手动引用
- **用途**: 专门的检查清单或审查指南
- **示例**: 安全审查清单、性能优化指南

## 🚀 使用方法

### 1. 自动生效的规则
大多数规则会根据你正在编辑的文件自动生效：
- 编辑 `.go` 文件 → Go 编码标准自动应用
- 编辑 `*_test.go` 文件 → 测试标准自动应用
- 编辑数据库相关文件 → 数据库模式自动应用

### 2. 手动引用规则
在对话中使用 `@` 符号引用特定规则：
```
@security-review-manual 请帮我审查这个认证模块的安全性
```

### 3. 项目上下文
项目基础信息(Always 规则)会在每次对话中自动提供，确保 AI 了解：
- 技术栈 (Go + Gin + PostgreSQL + Redis)
- 项目架构 (Clean Architecture)
- 核心功能和标准

## 🔧 规则管理

### 查看当前规则
在 Cursor 中：
1. 打开命令面板 (Cmd/Ctrl + Shift + P)
2. 搜索 "Cursor Rules"
3. 查看所有活跃的规则

### 添加新规则
1. 在适当的子目录下创建新的 `.mdc` 文件
2. 按照 MDC 格式编写规则内容
3. 设置正确的 frontmatter (description, globs, alwaysApply)

### 修改现有规则
直接编辑对应的 `.mdc` 文件，Cursor 会自动重新加载。

## 💡 最佳实践

### 规则编写
- 保持规则简洁明了 (目标 25 行，最多 50 行)
- 提供具体的示例和反例
- 使用描述性的规则名称
- 包含适当的表情符号和格式化

### 规则组织
- 按功能领域分组 (Go、API、数据库等)
- 避免规则重复和冲突
- 定期审查和清理过时的规则

### 性能考虑
- 避免过多的 Always 规则
- 使用具体的 glob 模式而不是通配符
- 保持规则内容精简

## 🔄 规则演进

随着项目发展，规则应该：
1. **有机演进** - 根据实际需求添加新规则
2. **定期清理** - 移除过时或冗余的规则  
3. **持续优化** - 根据 AI 表现调整规则内容
4. **团队同步** - 确保团队成员了解规则变更

## 📚 参考文档

- [Cursor Rules 官方文档](https://docs.cursor.com/context/rules)
- [MDC 格式规范](https://docs.cursor.com/context/rules#example-mdc-rule)
- [最佳实践指南](https://docs.cursor.com/context/rules#best-practices) 