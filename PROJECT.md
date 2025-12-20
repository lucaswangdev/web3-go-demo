# ERC20 代币余额追踪服务

一个完整的 Web3 后端实战项目，实时监听 ERC20 合约的 Transfer 事件，追踪地址余额并提供 REST API 查询。

## 📁 项目结构

```
web3-go-demo/
├── cmd/
│   └── tracker/
│       └── main.go              # 主程序入口
├── config/
│   └── config.go                # 配置管理
├── db/
│   ├── models.go                # 数据模型
│   └── postgres.go              # 数据库操作
├── listener/
│   └── event_listener.go        # 事件监听服务
├── api/
│   └── server.go                # REST API 服务
├── examples/                    # 示例代码
│   ├── 01_basic_connection.go
│   ├── 02_listen_events.go
│   └── 03_call_contract.go
├── .env.example                 # 环境变量示例
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🚀 功能特性

### 核心功能
- ✅ 实时监听 ERC20 合约的 Transfer 事件
- ✅ 自动更新地址余额
- ✅ 存储所有交易记录
- ✅ 提供 REST API 查询余额和交易
- ✅ 支持 WebSocket 自动重连
- ✅ 数据库事务保证数据一致性

### API 端点

#### 1. 健康检查
```bash
GET /health
```

#### 2. 查询地址余额
```bash
GET /api/v1/balance/:address
```

响应示例：
```json
{
  "address": "0x1234...",
  "balance": "1000000000000000000",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 3. 查询地址交易记录
```bash
GET /api/v1/transactions/:address?limit=10&offset=0
```

响应示例：
```json
{
  "address": "0x1234...",
  "count": 10,
  "limit": 10,
  "offset": 0,
  "transactions": [
    {
      "tx_hash": "0xabc...",
      "from_address": "0x1234...",
      "to_address": "0x5678...",
      "amount": "1000000",
      "block_number": 12345678,
      "timestamp": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 4. 查询交易详情
```bash
GET /api/v1/transaction/:txhash
```

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **区块链**: go-ethereum (geth)
- **数据库**: PostgreSQL + GORM
- **Web 框架**: Gin
- **配置**: 环境变量

## 📦 安装依赖

### 1. 安装 Go 依赖

```bash
go mod tidy
```

需要安装以下包：
```bash
go get github.com/ethereum/go-ethereum
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

### 2. 安装 PostgreSQL

**macOS**:
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian**:
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**创建数据库**:
```bash
psql -U postgres
CREATE DATABASE erc20_tracker;
\q
```

## ⚙️ 配置

### 1. 复制环境变量文件

```bash
cp .env.example .env
```

### 2. 编辑 `.env` 文件

```env
# 重要：替换为你的 Infura/Alchemy API Key
ETHEREUM_RPC=https://mainnet.infura.io/v3/YOUR_API_KEY
ETHEREUM_WS_RPC=wss://mainnet.infura.io/ws/v3/YOUR_API_KEY

# ERC20 合约地址（默认为 USDT）
CONTRACT_ADDRESS=0xdAC17F958D2ee523a2206206994597C13D831ec7

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=erc20_tracker

# API 端口
API_PORT=8080
```

### 3. 获取 API Key

#### Infura
1. 访问 https://infura.io/
2. 注册账号并创建项目
3. 复制 Project ID

#### Alchemy（备选）
1. 访问 https://www.alchemy.com/
2. 注册账号并创建应用
3. 复制 API Key

## 🎯 运行项目

### 1. 启动服务

```bash
# 方式 1: 直接运行
go run cmd/tracker/main.go

# 方式 2: 构建后运行
go build -o tracker cmd/tracker/main.go
./tracker
```

### 2. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 查询余额（替换为真实地址）
curl http://localhost:8080/api/v1/balance/0x123...

# 查询交易记录
curl http://localhost:8080/api/v1/transactions/0x123...?limit=5
```

### 3. 查看日志

服务启动后会输出类似信息：
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

当有 Transfer 事件时：
```
📤 Transfer 事件 | From: 0x123... | To: 0x456... | Amount: 1000000 | Block: 12345678
✅ 交易已保存: 0xabc...
✅ 余额已更新
```

## 🧪 测试

### 使用测试网
强烈建议先在测试网（Sepolia）测试：

1. 修改 `.env`:
```env
ETHEREUM_RPC=https://sepolia.infura.io/v3/YOUR_API_KEY
ETHEREUM_WS_RPC=wss://sepolia.infura.io/ws/v3/YOUR_API_KEY
CONTRACT_ADDRESS=0x... # 测试网上的 ERC20 合约地址
```

2. 部署或使用测试网上的 ERC20 合约
3. 获取测试币：https://sepoliafaucet.com/

## 📊 数据库表结构

### balances 表
```sql
CREATE TABLE balances (
    address VARCHAR(42) PRIMARY KEY,
    balance VARCHAR(78) NOT NULL,
    updated_at TIMESTAMP
);
```

### transactions 表
```sql
CREATE TABLE transactions (
    tx_hash VARCHAR(66) PRIMARY KEY,
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    amount VARCHAR(78) NOT NULL,
    block_number BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX idx_from_address ON transactions(from_address);
CREATE INDEX idx_to_address ON transactions(to_address);
CREATE INDEX idx_block_number ON transactions(block_number);
```

## 🔧 常见问题

### 1. 连接数据库失败
```
错误: failed to connect database
解决: 检查 PostgreSQL 是否运行，检查 .env 中的数据库配置
```

### 2. WebSocket 连接失败
```
错误: failed to connect to WebSocket
解决: 检查 Infura/Alchemy API Key 是否正确，检查网络连接
```

### 3. 没有收到事件
```
原因: USDT 等热门合约在主网上交易频繁，应该很快看到事件
检查: 确认合约地址正确，确认 WebSocket 连接成功
```

### 4. 余额不准确
```
原因: 从中间开始监听，之前的交易没有同步
解决: 需要同步历史数据（见进阶功能）
```

## 🚀 进阶功能

### 1. 同步历史数据
```go
// 查询历史事件
query := ethereum.FilterQuery{
    FromBlock: big.NewInt(startBlock),
    ToBlock:   big.NewInt(endBlock),
    Addresses: []common.Address{contractAddress},
}

logs, err := client.FilterLogs(context.Background(), query)
```

### 2. 添加缓存（Redis）
```go
// 查询余额时先查缓存
balance := redis.Get("balance:" + address)
if balance == "" {
    balance = database.GetBalance(address)
    redis.Set("balance:"+address, balance, 60*time.Second)
}
```

### 3. 支持多个合约
```go
// 配置中添加多个合约地址
contracts := []string{
    "0xdAC17F958D2ee523a2206206994597C13D831ec7", // USDT
    "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
}

// 为每个合约创建监听器
for _, addr := range contracts {
    go startListener(addr)
}
```

### 4. 消息队列（Kafka）
```go
// 将事件发送到 Kafka
producer.Send("transfer-events", event)

// 消费者处理事件
consumer.Subscribe("transfer-events", handleEvent)
```

## 📚 学习资源

- [go-ethereum 文档](https://geth.ethereum.org/docs/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [ERC20 标准](https://eips.ethereum.org/EIPS/eip-20)

## 🎓 下一步学习

完成这个项目后，你可以：

1. ✅ 理解 Web3 后端架构
2. ✅ 掌握事件监听和解析
3. ✅ 熟悉区块链数据存储

**接下来学习**：
- Solidity 合约编写（CryptoZombies）
- 合约交互（发送交易、签名）
- NFT (ERC721) 相关开发
- Gas 优化和 Nonce 管理
- 生产级部署（Docker, K8s）

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**🎯 记住：这是一个学习项目，不要在生产环境中直接使用。生产环境需要添加更多的错误处理、监控、安全措施等。**
