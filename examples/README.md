# Go-Ethereum 示例代码

这个目录包含了 7 个完整的 go-ethereum 使用示例，从基础到进阶，涵盖 Web3 后端开发的核心技能。

## 📚 示例列表

### 基础示例

#### [01_basic_connection.go](01_basic_connection.go)
**连接以太坊节点并查询基本信息**

学习内容：
- 连接到 Infura/Alchemy 节点
- 查询最新区块号
- 查询账户余额
- 获取区块详情
- 查询交易信息

适合：初学者，第一次接触 go-ethereum

```bash
go run examples/01_basic_connection.go
```

---

#### [02_listen_events.go](02_listen_events.go)
**监听智能合约事件**

学习内容：
- 使用 WebSocket 连接
- 订阅合约事件
- 解析 Transfer 事件
- 实时监控链上活动

适合：理解事件驱动的 Web3 开发

```bash
go run examples/02_listen_events.go
```

---

#### [03_call_contract.go](03_call_contract.go)
**调用智能合约方法**

学习内容：
- 使用 abigen 生成合约绑定
- 调用只读方法（balanceOf）
- 查询合约状态
- ERC20 标准接口

适合：需要与合约交互的场景

```bash
go run examples/03_call_contract.go
```

---

### 进阶示例

#### [04_send_transaction.go](04_send_transaction.go)
**发送交易到区块链**

学习内容：
- 加载私钥
- 签名交易
- 发送 ETH 转账
- 等待交易确认
- **安全提示和最佳实践**

适合：需要发送交易的应用

⚠️ **警告：仅在测试网使用！**

```bash
go run examples/04_send_transaction.go
```

---

#### [05_parse_logs.go](05_parse_logs.go)
**解析历史事件日志**

学习内容：
- 查询指定区块范围的事件
- 解析 indexed 和非 indexed 参数
- 过滤特定地址的事件
- 批量处理历史数据
- 统计分析链上活动

适合：数据分析、历史数据同步

```bash
go run examples/05_parse_logs.go
```

---

#### [06_estimate_gas.go](06_estimate_gas.go)
**Gas 费估算和优化**

学习内容：
- 估算交易 Gas Limit
- 获取 Gas Price 建议
- EIP-1559 Gas 费机制
- 计算交易成本
- Gas 优化策略

适合：关注成本优化的应用

```bash
go run examples/06_estimate_gas.go
```

---

#### [07_nonce_management.go](07_nonce_management.go)
**Nonce 管理和并发控制**

学习内容：
- 理解 Nonce 的作用
- 查询和管理 Nonce
- 并发安全的 Nonce 管理器
- 处理 Nonce 冲突
- 交易加速和取消

适合：高频交易、批量发送交易

```bash
go run examples/07_nonce_management.go
```

---

## 🎯 学习路径

### 第一周：基础（阶段二）
1. ✅ [`01_basic_connection.go`](01_basic_connection.go) - 连接节点
2. ✅ [`02_listen_events.go`](02_listen_events.go) - 监听事件
3. ✅ [`03_call_contract.go`](03_call_contract.go) - 调用合约

### 第二周：进阶
4. ⭐ [`04_send_transaction.go`](04_send_transaction.go) - 发送交易
5. ⭐ [`05_parse_logs.go`](05_parse_logs.go) - 解析日志
6. ⭐ [`06_estimate_gas.go`](06_estimate_gas.go) - Gas 估算

### 第三周：生产级
7. ⭐ [`07_nonce_management.go`](07_nonce_management.go) - Nonce 管理
8. 🚀 开始实战项目（参考 [`../cmd/tracker/main.go`](../cmd/tracker/main.go)）

---

## 📋 前提条件

### 1. 安装 Go
```bash
go version  # 需要 Go 1.21+
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 获取 API Key
- **Infura**: https://infura.io/
- **Alchemy**: https://www.alchemy.com/

### 4. 修改配置
在每个示例文件中，替换：
```go
rpcURL := "https://mainnet.infura.io/v3/YOUR_API_KEY"
```

---

## 🔑 关键概念对照表

| 概念 | 说明 | 相关示例 |
|------|------|---------|
| **Ethereum Client** | 以太坊客户端连接 | 01 |
| **Block** | 区块 | 01 |
| **Transaction** | 交易 | 01, 04 |
| **Event/Log** | 事件/日志 | 02, 05 |
| **Smart Contract** | 智能合约 | 03, 04 |
| **ABI** | 应用程序二进制接口 | 02, 03, 05 |
| **Gas** | 计算费用 | 04, 06 |
| **Nonce** | 交易序号 | 04, 07 |
| **Private Key** | 私钥 | 04 |
| **WebSocket** | 双向通信协议 | 02 |

---

## 🛠️ 使用技巧

### 运行单个示例
```bash
go run examples/01_basic_connection.go
```

### 查看示例代码
每个示例都包含详细的：
- 功能说明
- 使用场景
- 注释说明
- 知识点总结
- 最佳实践

### 修改和实验
推荐：
1. 先运行示例，看到输出
2. 阅读代码和注释
3. 修改参数，观察变化
4. 尝试自己的需求

---

## ⚠️ 重要提示

### 关于测试网
- 使用 **Sepolia** 或 **Goerli** 测试网
- 从水龙头获取免费测试币
  - https://sepoliafaucet.com/
  - https://goerlifaucet.com/

### 关于主网
- ⚠️ 主网操作需要真实 ETH
- 💰 每笔交易都会消耗 Gas 费
- 🔒 永远不要暴露私钥

### 关于私钥
- ❌ 不要在代码中硬编码私钥
- ❌ 不要提交包含私钥的代码
- ✅ 使用环境变量
- ✅ 生产环境使用 KMS

---

## 📊 示例对比

| 示例 | 难度 | 耗时 | 网络 | 需要私钥 | 费用 |
|------|------|------|------|---------|------|
| 01 | ⭐ | 5分钟 | 主网/测试网 | ❌ | 免费 |
| 02 | ⭐⭐ | 10分钟 | 主网/测试网 | ❌ | 免费 |
| 03 | ⭐⭐ | 15分钟 | 主网/测试网 | ❌ | 免费 |
| 04 | ⭐⭐⭐ | 20分钟 | 测试网 | ✅ | Gas费 |
| 05 | ⭐⭐⭐ | 15分钟 | 主网/测试网 | ❌ | 免费 |
| 06 | ⭐⭐ | 10分钟 | 主网/测试网 | ❌ | 免费 |
| 07 | ⭐⭐⭐⭐ | 20分钟 | 测试网 | ❌ | 免费 |

---

## 🔗 参考资源

### 官方文档
- [go-ethereum 文档](https://geth.ethereum.org/docs/)
- [以太坊开发文档](https://ethereum.org/developers)

### 在线工具
- [Etherscan](https://etherscan.io/) - 区块链浏览器
- [Remix IDE](https://remix.ethereum.org/) - 在线 Solidity IDE
- [Gas Tracker](https://etherscan.io/gastracker) - Gas 价格

### 推荐阅读
- [Go Ethereum Book](https://goethereumbook.org/)
- [Solidity by Example](https://solidity-by-example.org/)

---

## 💡 下一步

完成这些示例后，你可以：

### 1. 深化理解
- 修改示例代码
- 尝试其他合约
- 部署自己的合约

### 2. 实战项目
- 开始 [ERC20 余额追踪服务](../PROJECT.md)
- 构建自己的 Web3 应用

### 3. 学习 Solidity
- [CryptoZombies](https://cryptozombies.io/)
- [智能合约入门指南](../智能合约入门指南.md)

---

## 🤝 贡献

发现问题或有改进建议？欢迎：
- 提交 Issue
- 创建 Pull Request
- 分享你的使用经验

---

## 📄 许可证

MIT License

---

**🎓 记住：实践是最好的学习方式！运行代码，修改参数，观察结果。**
