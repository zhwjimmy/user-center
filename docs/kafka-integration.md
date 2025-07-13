# Kafka 消息队列集成指南

## 概述

本项目集成了Apache Kafka作为消息队列系统，用于处理用户相关的异步事件。通过Kafka，我们可以实现业务解耦、提高系统性能，并支持事件驱动架构。

## 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Service  │───▶│  Kafka Producer │───▶│     Kafka       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Event Handlers  │◀───│ Kafka Consumer  │◀───│   Topic: user   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 功能特性

### ✅ 已实现的事件类型

1. **用户注册事件** (`user.registered`)
   - 发送欢迎邮件
   - 初始化用户配置
   - 记录注册统计

2. **用户登录事件** (`user.logged_in`)
   - 记录登录日志
   - 更新最后登录时间
   - 检查异常登录

3. **密码变更事件** (`user.password_changed`)
   - 发送安全通知邮件
   - 记录安全日志

4. **用户状态变更事件** (`user.status_changed`)
   - 发送状态变更通知
   - 更新缓存

5. **用户删除事件** (`user.deleted`)
   - 清理用户数据
   - 发送删除确认

6. **用户更新事件** (`user.updated`)
   - 更新缓存
   - 同步到外部系统

### 🔧 技术特性

- **高性能**：使用IBM/sarama客户端，支持批处理和压缩
- **可靠性**：支持消息确认、重试机制和错误处理
- **可扩展性**：支持消费者组，可水平扩展
- **监控**：集成Prometheus监控和日志追踪
- **优雅关闭**：支持优雅停机，确保消息不丢失

## 配置说明

### Kafka配置 (`configs/config.yaml`)

```yaml
kafka:
  brokers: ["localhost:9092"]  # Kafka集群地址
  topics:
    user_events: "user.events"        # 用户事件主题
    user_notifications: "user.notifications"  # 通知主题
    user_analytics: "user.analytics"  # 分析主题
  group_id: "usercenter"       # 消费者组ID
```

### 生产者配置

- **确认机制**：等待所有副本确认 (`WaitForAll`)
- **重试**：最大重试3次，退避时间100ms
- **批处理**：100条消息或1MB触发发送
- **压缩**：使用Snappy压缩
- **幂等性**：启用幂等生产者

### 消费者配置

- **偏移量**：从最新位置开始消费
- **负载均衡**：轮询策略
- **自动提交**：1秒间隔自动提交偏移量
- **会话超时**：10秒
- **心跳间隔**：3秒

## 使用示例

### 1. 发布事件

```go
// 在业务服务中发布用户注册事件
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, string, error) {
    // ... 业务逻辑 ...
    
    // 发布用户注册事件
    if err := s.eventService.PublishUserRegisteredEvent(ctx, createdUser); err != nil {
        s.logger.Error("Failed to publish user registered event", zap.Error(err))
        // 不返回错误，避免影响主要业务流程
    }
    
    return createdUser, token, nil
}
```

### 2. 处理事件

```go
// 在消息处理器中处理用户注册事件
func (h *UserEventHandler) HandleUserRegistered(ctx context.Context, event *event.UserRegisteredEvent) error {
    h.logger.Info("Processing user registered event",
        zap.String("user_id", event.UserID),
        zap.String("username", event.Username),
        zap.String("email", event.Email),
    )

    // 发送欢迎邮件
    if err := h.sendWelcomeEmail(ctx, event); err != nil {
        h.logger.Error("Failed to send welcome email", zap.Error(err))
    }

    return nil
}
```

### 3. 自定义事件

```go
// 定义新的事件类型
type UserProfileUpdatedEvent struct {
    BaseEvent
    Username string `json:"username"`
    Email    string `json:"email"`
    Changes  map[string]interface{} `json:"changes"`
}

// 在生产者中添加处理逻辑
case *event.UserProfileUpdatedEvent:
    topic = p.config.GetTopicName("user_events")
    key = e.UserID
    value, err = e.ToJSON()
    // ...
```

## 部署和运行

### 1. 启动Kafka

```bash
# 使用Docker Compose启动Kafka
docker-compose up -d kafka zookeeper

# 检查Kafka状态
docker-compose ps kafka
```

### 2. 创建主题

```bash
# 创建用户事件主题
kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic user.events \
  --partitions 3 \
  --replication-factor 1

# 列出所有主题
kafka-topics --list --bootstrap-server localhost:9092
```

### 3. 启动应用

```bash
# 构建应用
make build

# 启动应用
./bin/usercenter
```

### 4. 测试消息

```bash
# 注册新用户（会触发用户注册事件）
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 查看Kafka消息
kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic user.events \
  --from-beginning
```

## 监控和调试

### 1. 查看消费者组状态

```bash
kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group usercenter
```

### 2. 查看主题详情

```bash
kafka-topics --describe --bootstrap-server localhost:9092 --topic user.events
```

### 3. 监控指标

应用暴露以下Prometheus指标：

- `kafka_producer_messages_total` - 生产者发送消息总数
- `kafka_consumer_messages_total` - 消费者处理消息总数
- `kafka_producer_errors_total` - 生产者错误总数
- `kafka_consumer_errors_total` - 消费者错误总数

### 4. 日志查看

```bash
# 查看应用日志
tail -f logs/usercenter.log | grep kafka

# 查看Docker容器日志
docker-compose logs -f kafka
```

## 故障排除

### 常见问题

1. **连接失败**
   ```
   Error: failed to create kafka producer: kafka: client has run out of available brokers
   ```
   - 检查Kafka是否正在运行
   - 验证broker地址配置
   - 检查网络连接

2. **消息发送失败**
   ```
   Error: kafka: Failed to produce message to topic
   ```
   - 检查主题是否存在
   - 验证权限配置
   - 查看Kafka日志

3. **消费者无法消费**
   ```
   Error: kafka: error while consuming messages
   ```
   - 检查消费者组配置
   - 验证主题权限
   - 检查偏移量设置

### 解决方案

1. **重置消费者组偏移量**
   ```bash
   kafka-consumer-groups --bootstrap-server localhost:9092 \
     --group usercenter --reset-offsets --to-earliest \
     --topic user.events --execute
   ```

2. **清理主题数据**
   ```bash
   kafka-topics --delete --bootstrap-server localhost:9092 --topic user.events
   ```

3. **检查Kafka健康状态**
   ```bash
   kafka-broker-api-versions --bootstrap-server localhost:9092
   ```

## 最佳实践

### 1. 事件设计

- 使用版本化的事件结构
- 包含足够的上下文信息
- 保持事件的幂等性
- 使用有意义的事件类型名称

### 2. 错误处理

- 实现重试机制
- 记录详细的错误日志
- 使用死信队列处理失败消息
- 监控错误率和延迟

### 3. 性能优化

- 合理设置批处理大小
- 使用压缩减少网络传输
- 优化分区数量
- 监控消费者延迟

### 4. 安全考虑

- 使用SSL/TLS加密传输
- 配置SASL认证
- 限制主题访问权限
- 定期轮换密钥

## 扩展功能

### 1. 添加新的事件类型

1. 在 `internal/kafka/event/` 中定义新事件
2. 在 `producer.go` 中添加处理逻辑
3. 在 `consumer.go` 中添加消费逻辑
4. 在 `handler.go` 中实现业务处理

### 2. 集成外部服务

- 邮件服务集成
- 短信通知服务
- 数据分析平台
- 第三方API调用

### 3. 高级特性

- 事件溯源 (Event Sourcing)
- CQRS 模式
- 分布式事务 (Saga)
- 事件重放功能

## 参考资料

- [Apache Kafka 官方文档](https://kafka.apache.org/documentation/)
- [IBM Sarama 客户端文档](https://github.com/IBM/sarama)
- [事件驱动架构最佳实践](https://microservices.io/patterns/data/event-driven-architecture.html)
- [Kafka 性能调优指南](https://kafka.apache.org/documentation/#config) 