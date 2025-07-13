#!/bin/bash

# Kafka 功能测试脚本
# 用于验证Kafka消息队列功能是否正常工作

set -e

echo "🚀 开始测试 Kafka 消息队列功能..."

# 检查 Kafka 是否运行
echo "📋 检查 Kafka 服务状态..."
if ! docker-compose ps kafka | grep -q "Up"; then
    echo "❌ Kafka 服务未运行，请先启动 Docker Compose"
    echo "   运行: docker-compose up -d kafka zookeeper"
    exit 1
fi

echo "✅ Kafka 服务正在运行"

# 检查主题是否存在，如果不存在则创建
TOPIC="user.events"
echo "📋 检查主题: $TOPIC"

if ! kafka-topics --bootstrap-server localhost:9092 --list | grep -q "^$TOPIC$"; then
    echo "📝 创建主题: $TOPIC"
    kafka-topics --create \
        --bootstrap-server localhost:9092 \
        --topic $TOPIC \
        --partitions 3 \
        --replication-factor 1
    echo "✅ 主题创建成功"
else
    echo "✅ 主题已存在"
fi

# 启动后台消费者来监听消息
echo "📋 启动消费者监听消息..."
kafka-console-consumer \
    --bootstrap-server localhost:9092 \
    --topic $TOPIC \
    --group test-consumer \
    --timeout-ms 10000 > /tmp/kafka-messages.log 2>&1 &

CONSUMER_PID=$!
echo "✅ 消费者已启动 (PID: $CONSUMER_PID)"

# 等待消费者准备就绪
sleep 2

# 检查应用是否运行
echo "📋 检查应用服务状态..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ 应用服务未运行，请先启动应用"
    echo "   运行: make run 或 go run cmd/usercenter/main.go"
    kill $CONSUMER_PID 2>/dev/null || true
    exit 1
fi

echo "✅ 应用服务正在运行"

# 测试用户注册（会触发 Kafka 事件）
echo "📋 测试用户注册事件..."
TIMESTAMP=$(date +%s)
TEST_EMAIL="test-${TIMESTAMP}@example.com"
TEST_USERNAME="testuser${TIMESTAMP}"

REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/register \
    -H "Content-Type: application/json" \
    -H "X-Request-ID: test-request-${TIMESTAMP}" \
    -d "{
        \"username\": \"$TEST_USERNAME\",
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"password123\",
        \"first_name\": \"Test\",
        \"last_name\": \"User\"
    }")

if echo "$REGISTER_RESPONSE" | grep -q "User registered successfully"; then
    echo "✅ 用户注册成功"
    USER_ID=$(echo "$REGISTER_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   用户ID: $USER_ID"
else
    echo "❌ 用户注册失败"
    echo "   响应: $REGISTER_RESPONSE"
    kill $CONSUMER_PID 2>/dev/null || true
    exit 1
fi

# 等待消息传播
echo "⏳ 等待 Kafka 消息传播..."
sleep 3

# 检查是否收到消息
echo "📋 检查 Kafka 消息..."
if [ -f /tmp/kafka-messages.log ] && [ -s /tmp/kafka-messages.log ]; then
    echo "✅ 收到 Kafka 消息:"
    echo "----------------------------------------"
    cat /tmp/kafka-messages.log | head -5
    echo "----------------------------------------"
    
    # 验证消息内容
    if grep -q "user.registered" /tmp/kafka-messages.log && \
       grep -q "$TEST_USERNAME" /tmp/kafka-messages.log && \
       grep -q "$TEST_EMAIL" /tmp/kafka-messages.log; then
        echo "✅ 消息内容验证成功"
    else
        echo "⚠️  消息内容验证失败，但消息已发送"
    fi
else
    echo "❌ 未收到 Kafka 消息"
    echo "   请检查应用日志和 Kafka 配置"
fi

# 清理
kill $CONSUMER_PID 2>/dev/null || true
rm -f /tmp/kafka-messages.log

# 显示消费者组信息
echo "📋 消费者组状态:"
kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group usercenter 2>/dev/null || echo "   消费者组 'usercenter' 未找到"

# 显示主题信息
echo "📋 主题详情:"
kafka-topics --describe --bootstrap-server localhost:9092 --topic $TOPIC

echo ""
echo "🎉 Kafka 功能测试完成！"
echo ""
echo "📚 更多信息:"
echo "   - 查看应用日志: tail -f logs/usercenter.log"
echo "   - 查看 Kafka 日志: docker-compose logs -f kafka"
echo "   - 监控消息: kafka-console-consumer --bootstrap-server localhost:9092 --topic user.events"
echo "   - 文档: docs/kafka-integration.md" 