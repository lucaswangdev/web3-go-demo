# 🚀 快速启动指南

## 前提条件

- Go 1.21+ 已安装
- PostgreSQL 已安装并运行
- Infura 或 Alchemy API Key

## 步骤 1: 设置数据库

```bash
# 启动 PostgreSQL (macOS)
brew services start postgresql@15

# 或 (Ubuntu)
sudo systemctl start postgresql

# 创建数据库
psql -U postgres
CREATE DATABASE erc20_tracker;
\q
```

## 步骤 2: 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入你的 API Key
vim .env
```

**重要**：必须修改以下配置：
```env
ETHEREUM_RPC=https://mainnet.infura.io/v3/YOUR_API_KEY  # 替换为你的 API Key
ETHEREUM_WS_RPC=wss://mainnet.infura.io/ws/v3/YOUR_API_KEY
DB_PASSWORD=your_password  # 替换为你的数据库密码
```

## 步骤 3: 安装依赖

```bash
go mod tidy
```

## 步骤 4: 运行项目

```bash
go run cmd/tracker/main.go
```

你应该看到类似输出：
```
🚀 启动 ERC20 代币余额追踪服务...
配置加载完成: Contract=0xdAC17F958D2ee523a2206206994597C13D831ec7
数据库连接成功，表迁移完成
开始监听合约事件: 0xdAC17F958D2ee523a2206206994597C13D831ec7
成功订阅事件，等待 Transfer 事件...
✅ 服务启动成功
📡 监听合约: 0xdAC17F958D2ee523a2206206994597C13D831ec7
🌐 API 服务: http://localhost:8080
```

## 步骤 5: 测试 API

打开新的终端窗口：

```bash
# 健康检查
curl http://localhost:8080/health

# 查询 Binance 热钱包的 USDT 余额（示例）
curl http://localhost:8080/api/v1/balance/0x28C6c06298d514Db089934071355E5743bf21d60

# 查询交易记录
curl http://localhost:8080/api/v1/transactions/0x28C6c06298d514Db089934071355E5743bf21d60?limit=5
```

## 常见问题

### Q1: 数据库连接失败？
```bash
# 检查 PostgreSQL 是否运行
pg_isready

# 重启 PostgreSQL
brew services restart postgresql@15
```

### Q2: 没有收到事件？
- 确认 API Key 正确
- 确认 WebSocket 连接成功
- USDT 主网交易频繁，应该很快看到事件

### Q3: 想使用测试网？
修改 .env：
```env
ETHEREUM_WS_RPC=wss://sepolia.infura.io/ws/v3/YOUR_API_KEY
CONTRACT_ADDRESS=0x...  # 测试网合约地址
```

## 下一步

1. ✅ 服务正常运行后，观察控制台输出的 Transfer 事件
2. ✅ 使用 curl 或 Postman 测试 API
3. ✅ 查看数据库中的数据：
   ```bash
   psql -U postgres -d erc20_tracker
   SELECT * FROM balances LIMIT 10;
   SELECT * FROM transactions LIMIT 10;
   ```

4. ✅ 阅读 [`PROJECT.md`](PROJECT.md:1) 了解完整功能
5. ✅ 尝试添加新功能（支持多合约、缓存等）

## 性能优化建议

### 开发环境
- 使用测试网（Sepolia）
- 监听小流量合约

### 生产环境
- 添加 Redis 缓存
- 使用消息队列（Kafka）
- 多个 WebSocket 连接（高可用）
- 添加监控（Prometheus）

---

🎉 **恭喜！你已经成功搭建了第一个 Web3 后端服务！**
