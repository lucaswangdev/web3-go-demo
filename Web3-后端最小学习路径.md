# Web3 后端最小学习路径（基于Go）

> 目标：用最少的时间掌握核心技能，快速进入实战

---

## 📌 阶段一：Solidity 速成（2-3周）

### 必须掌握
- [ ] 数据类型（uint, address, mapping, array）
- [ ] 函数（view, pure, payable）
- [ ] 事件（Event）与日志
- [ ] 修饰器（modifier）
- [ ] ERC20 标准（transfer, approve, transferFrom）
- [ ] ERC721 标准（NFT 基础）

### 学习资源
- **CryptoZombies**（交互式教程）：https://cryptozombies.io/
- **OpenZeppelin 合约库**：https://docs.openzeppelin.com/contracts/

### 实战练习
```solidity
// 写一个简单的 ERC20 代币合约
// 写一个简单的 NFT 合约
// 理解 Transfer 事件的参数
```

### ❌ 可以跳过
- 复杂的继承和多态
- 内联汇编（assembly）
- 高级 Gas 优化技巧
- 合约升级模式（可后续学习）

---

## 📌 阶段二：go-ethereum 核心（1-2周）

### 必须掌握
- [ ] 连接以太坊节点（Infura/Alchemy）
- [ ] 读取区块和交易信息
- [ ] 查询账户余额
- [ ] 发送交易（签名 + 发送）
- [ ] 调用合约方法（读写分离）
- [ ] 监听合约事件（WebSocket）
- [ ] 解析 ABI 和事件日志

### 核心代码片段

#### 1. 连接节点
```go
package main

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "log"
)

func main() {
    // 连接 Infura 节点
    client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 获取最新区块号
    blockNumber, err := client.BlockNumber(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    log.Println("最新区块:", blockNumber)
}
```

#### 2. 查询余额
```go
account := common.HexToAddress("0x...")
balance, err := client.BalanceAt(context.Background(), account, nil)
```

#### 3. 监听合约事件（核心！）
```go
// 使用 WebSocket 连接
wsClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_API_KEY")

// 解析 ABI
contractABI, err := abi.JSON(strings.NewReader(ERC20ABI))

// 订阅 Transfer 事件
query := ethereum.FilterQuery{
    Addresses: []common.Address{contractAddress},
}

logs := make(chan types.Log)
sub, err := wsClient.SubscribeFilterLogs(context.Background(), query, logs)

for {
    select {
    case log := <-logs:
        // 解析事件
        event := struct {
            From  common.Address
            To    common.Address
            Value *big.Int
        }{}
        err := contractABI.UnpackIntoInterface(&event, "Transfer", log.Data)
        
        // 处理业务逻辑
        handleTransfer(event.From, event.To, event.Value)
    }
}
```

#### 4. 调用合约方法
```go
// 读操作（view/pure）
instance, err := NewERC20(contractAddress, client)
balance, err := instance.BalanceOf(&bind.CallOpts{}, address)

// 写操作（需要签名）
auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
tx, err := instance.Transfer(auth, toAddress, amount)
```

### 学习资源
- **go-ethereum 官方文档**：https://geth.ethereum.org/docs/developers/dapp-developer/native
- **使用 abigen 生成合约绑定**
- **ethclient 包核心 API**

### ❌ 可以跳过
- 矿池开发
- 节点同步机制
- 共识算法实现
- P2P 网络层

---

## 📌 阶段三：实战项目（2-4周）

### 项目一：ERC20 代币余额追踪服务 ⭐⭐⭐⭐⭐

#### 功能需求
1. 监听指定 ERC20 合约的 Transfer 事件
2. 实时更新每个地址的代币余额
3. 提供 REST API 查询余额
4. 支持历史交易查询

#### 技术栈
- **Go** + **go-ethereum**
- **PostgreSQL**（存储余额和交易记录）
- **Gin**（Web 框架）
- **WebSocket**（实时监听）

#### 数据库设计
```sql
-- 地址余额表
CREATE TABLE balances (
    address VARCHAR(42) PRIMARY KEY,
    balance NUMERIC(78, 0),
    updated_at TIMESTAMP
);

-- 交易记录表
CREATE TABLE transactions (
    tx_hash VARCHAR(66) PRIMARY KEY,
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    amount NUMERIC(78, 0),
    block_number BIGINT,
    timestamp TIMESTAMP
);
```

#### 核心架构
```
┌─────────────┐
│   以太坊网络  │
└──────┬──────┘
       │ WebSocket 监听 Transfer 事件
       ▼
┌─────────────────────┐
│  事件监听服务 (Go)    │
│  - 解析事件         │
│  - 更新余额         │
│  - 落库存储         │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   PostgreSQL        │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   REST API (Gin)    │
│  - GET /balance/:addr│
│  - GET /txs/:addr   │
└─────────────────────┘
```

#### 关键代码结构
```
project/
├── main.go                 # 程序入口
├── listener/
│   └── event_listener.go   # 事件监听
├── db/
│   ├── postgres.go         # 数据库连接
│   └── models.go           # 数据模型
├── api/
│   └── server.go           # REST API
├── contracts/
│   ├── erc20.go            # abigen 生成的合约绑定
│   └── erc20.abi           # 合约 ABI
└── config/
    └── config.go           # 配置管理
```

#### 实现步骤
1. [ ] 使用 abigen 生成 ERC20 合约 Go 绑定
2. [ ] 实现 WebSocket 连接和事件订阅
3. [ ] 实现事件解析和数据库更新逻辑
4. [ ] 实现 REST API（查询余额和交易）
5. [ ] 添加错误处理和重连机制
6. [ ] 添加日志和监控

---

### 项目二：NFT 铸造和查询服务（可选）

#### 功能需求
1. 部署 ERC721 合约
2. 提供 API 铸造 NFT
3. 查询 NFT 所有者
4. 监听 Transfer 事件更新所有权

---

## 📌 阶段四：生产级优化（进阶）

### 必须掌握
- [ ] 高可用架构（多节点冗余）
- [ ] 消息队列（Kafka/RabbitMQ）
- [ ] 事件去重和幂等性
- [ ] Gas 费估算和优化
- [ ] 私钥安全存储（KMS）
- [ ] 交易重试和 Nonce 管理
- [ ] 监控告警（Prometheus + Grafana）

### 架构升级
```
┌─────────────┐     ┌─────────────┐
│  Infura 节点 │     │  Alchemy 节点│
└──────┬──────┘     └──────┬──────┘
       │多节点冗余          │
       └────────┬───────────┘
                ▼
      ┌──────────────────┐
      │  事件监听服务 (Go) │
      └────────┬──────────┘
               │
               ▼
      ┌──────────────────┐
      │   Kafka 消息队列  │
      └────────┬──────────┘
               │
        ┌──────┴──────┐
        ▼             ▼
   ┌────────┐   ┌────────┐
   │ Consumer1│   │ Consumer2│
   └────┬───┘   └────┬───┘
        │            │
        └──────┬─────┘
               ▼
        ┌────────────┐
        │ PostgreSQL │
        └────────────┘
```

---

## 🎯 学习建议

### 1. **边学边做**
- 不要只看文档，必须写代码
- 每个知识点都写一个小 demo
- 从简单合约开始（ERC20）

### 2. **使用测试网**
- Goerli / Sepolia 测试网
- 从水龙头获取免费测试币
- 先在测试网验证，再上主网

### 3. **阅读优秀项目**
- **Uniswap SDK**：https://github.com/Uniswap/sdk-core
- **OpenZeppelin**：https://github.com/OpenZeppelin/openzeppelin-contracts
- **go-ethereum examples**：https://goethereumbook.org/

### 4. **关注安全**
- 私钥绝对不能明文存储
- 使用环境变量 + KMS
- 永远不要在前端暴露私钥

---

## 🔗 重要资源

### 开发工具
- **Remix IDE**：在线 Solidity IDE
- **Hardhat**：智能合约开发框架
- **abigen**：生成 Go 合约绑定
- **Etherscan**：区块链浏览器

### 测试网水龙头
- Goerli：https://goerlifaucet.com/
- Sepolia：https://sepoliafaucet.com/

### 节点服务商
- **Infura**：https://infura.io/
- **Alchemy**：https://www.alchemy.com/
- **QuickNode**：https://www.quicknode.com/

---

## 📊 时间线总结

| 阶段 | 内容 | 时间 | 产出 |
|------|------|------|------|
| 阶段一 | Solidity 基础 | 2-3周 | 能读懂合约代码 |
| 阶段二 | go-ethereum | 1-2周 | 能调用合约和监听事件 |
| 阶段三 | 实战项目 | 2-4周 | 完整的 Web3 后端服务 |
| 阶段四 | 生产优化 | 持续 | 生产级系统架构 |

**总计：6-10周可以入门并完成实战项目**

---

## ✅ 检查清单

### Solidity 基础
- [ ] 能读懂 ERC20 合约
- [ ] 能读懂 ERC721 合约
- [ ] 理解 Transfer 事件

### Go-Ethereum
- [ ] 能连接节点查询数据
- [ ] 能发送交易
- [ ] 能监听合约事件
- [ ] 能调用合约方法

### 实战能力
- [ ] 完成 ERC20 余额追踪项目
- [ ] 理解生产架构模式
- [ ] 掌握错误处理和重试

---

## 🚀 快速开始

### 今天就开始（Day 1）
```bash
# 1. 安装 go-ethereum
go get github.com/ethereum/go-ethereum

# 2. 注册 Infura 获取 API Key
# https://infura.io/

# 3. 写第一个程序：连接节点并查询最新区块
```

### 本周目标（Week 1）
- 完成 CryptoZombies 前 3 课
- 能读懂 ERC20 合约
- 用 Go 连接节点并查询余额

---

**记住：Web3 开发的核心是理解链上数据和事件，Go 只是工具。先理解业务逻辑，再写代码！**
